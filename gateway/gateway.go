package gateway

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

var jsonRe = regexp.MustCompile(`(?i)\bjson\b`)

// Gateway is the main struct
type Gateway struct {
	Options *Options
	server  *http.Server
	pushers *pushers
}

func addToVary(r *http.Response, header string) {
	v := r.Header.Get("Vary")
	if v == "" {
		r.Header.Set("Vary", header)
		return
	}

	r.Header.Set("Vary", v+", "+header)
}

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	pusher, explicitRequest, explicitRequestID := g.getPusher(rw, req)

	rp := httputil.NewSingleHostReverseProxy(g.Options.Upstream)
	rp.ModifyResponse = func(r *http.Response) error {
		query := req.URL.Query()
		var useFieldsHeader, useFieldsQuery, usePreloadHeader, usePreloadQuery bool

		if len(req.Header["Fields"]) > 0 {
			useFieldsHeader = true
		} else if len(query["fields"]) > 0 {
			useFieldsQuery = true
		}

		if len(req.Header["Preload"]) > 0 {
			usePreloadHeader = true
		} else if len(query["preload"]) > 0 {
			usePreloadQuery = true
		}

		if (!useFieldsHeader && !useFieldsQuery && !usePreloadHeader && !usePreloadQuery) || !jsonRe.MatchString(r.Header.Get("Content-Type")) {
			// No Vulcain hints, or not JSON: don't modify the response
			if pusher != nil {
				if explicitRequest {
					g.pushers.remove(explicitRequestID)
				} else {
					pusher.Done()
				}
			}

			return nil
		}

		currentBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		tree := &node{}
		if usePreloadHeader {
			tree.importPointers(Preload, req.Header["Preload"])
		}
		if usePreloadQuery {
			tree.importPointers(Preload, query["preload"])
		}
		if useFieldsHeader {
			tree.importPointers(Fields, req.Header["Fields"])
		}
		if useFieldsQuery {
			tree.importPointers(Fields, query["fields"])
		}

		newBody := traverseJSON(currentBody, tree, useFieldsHeader || useFieldsQuery, func(u *url.URL, n *node) {
			if usePreloadQuery || useFieldsQuery {
				urlRewriter(u, n)
			}

			if !usePreloadHeader && !usePreloadQuery {
				return
			}

			uStr := u.String()
			// TODO: allow to disable Server Push from the config
			if !u.IsAbs() && pusher != nil {
				pushOptions := &http.PushOptions{Header: req.Header}
				pushOptions.Header.Del("Preload")
				pushOptions.Header.Del("Fields")
				pushOptions.Header.Del("Te") // Trailing headers aren't supported by Firefox for pushes, and we don't use them

				if usePreloadHeader {
					for _, pp := range n.strings(Preload, "") {
						if pp != "/" {
							pushOptions.Header.Add("Preload", pp)
						}
					}
				}
				if useFieldsHeader {
					for _, fp := range n.strings(Fields, "") {
						if fp != "/" {
							pushOptions.Header.Add("Fields", fp)
						}
					}
				}

				// HTTP/2, and relative relation, push!
				err := pusher.Push(uStr, pushOptions)
				if err == nil {
					log.WithFields(log.Fields{"relation": uStr}).Debug("Relation pushed")
					return
				}
				log.WithFields(log.Fields{"relation": uStr, "reason": err.Error()}).Debug("Failed to push")
				if _, ok := err.(*relationAlreadyPushedError); ok {
					// Don't add the preload header for something already pushed
					return
				}
			}

			// Use preload Link headers as fallback (https://www.w3.org/TR/preload/)
			// TODO: allow to set the nopush attribute using the configuration (https://www.w3.org/TR/preload/#server-push-http-2)
			// TODO: send 103 early hints responses (https://tools.ietf.org/html/rfc8297)
			r.Header.Add("Link", "<"+uStr+">; rel=preload; as=fetch")
			log.WithFields(log.Fields{"relation": uStr}).Debug("Link preload header added")
		})

		if useFieldsHeader {
			addToVary(r, "Fields")
		}
		if usePreloadHeader {
			addToVary(r, "Preload")
		}

		if pusher != nil {
			if explicitRequest {
				pusher.Wait()
				g.pushers.remove(explicitRequestID)
			} else {
				// Relations pushed
				pusher.Done()
			}
		}

		newBodyBuffer := bytes.NewBuffer(newBody)
		r.Body = ioutil.NopCloser(newBodyBuffer)
		r.Header["Content-Length"] = []string{fmt.Sprint(newBodyBuffer.Len())}

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		// Adapted from the default ErrorHandler
		log.Errorf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)

		if pusher != nil && !explicitRequest {
			pusher.Done()
		}
	}
	rp.ServeHTTP(rw, req)
}

func (g *Gateway) getPusher(rw http.ResponseWriter, req *http.Request) (p *waitPusher, explicitRequest bool, explicitRequestID string) {
	internalPusher, ok := rw.(http.Pusher)
	if !ok {
		// Not an HTTP/2 connection
		return nil, false, ""
	}

	// Need https://github.com/golang/go/issues/20566 to get rid of this hack
	explicitRequestID = req.Header.Get("Vulcain-Explicit-Request")
	if explicitRequestID != "" {
		p, ok = g.pushers.get(explicitRequestID)
		if !ok {
			// Should not happen, is an attacker forging an evil request?
			log.WithFields(log.Fields{"uri": req.RequestURI, "explicitRequestID": explicitRequestID}).Debug("Pusher not found")
			explicitRequestID = ""
		}
	}

	if explicitRequestID == "" {
		// Explicit request
		explicitRequestID = uuid.Must(uuid.NewV4()).String()
		p = newWaitPusher(internalPusher, g.Options.MaxPushes)
		req.Header.Set("Vulcain-Explicit-Request", explicitRequestID)
		g.pushers.add(explicitRequestID, p)

		return p, true, explicitRequestID
	}

	return p, false, explicitRequestID
}

// NewGatewayFromEnv creates a gateway using the configuration set in env vars
func NewGatewayFromEnv() (*Gateway, error) {
	options, err := NewOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	return NewGateway(options), nil
}

// NewGateway creates a Vulcain gateway instance
func NewGateway(options *Options) *Gateway {
	return &Gateway{
		options,
		nil,
		&pushers{pusherMap: make(map[string]*waitPusher)},
	}
}
