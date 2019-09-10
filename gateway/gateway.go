package gateway

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sync"

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
	rp := httputil.NewSingleHostReverseProxy(g.Options.Upstream)
	defer rp.ServeHTTP(rw, req)

	var wg sync.WaitGroup
	parentPusher, childrenPusher, requestID := g.retrievePushers(rw, req)
	if childrenPusher != nil {
		// Block until reverse proxy's goroutine finishes
		wg.Add(1)
		// Cleanup
		defer g.pushers.remove(requestID)
		defer wg.Wait()
		// Wait for the subrequests to finish
		defer childrenPusher.Wait()

		if parentPusher != nil {
			// Mark this subrequest as finished
			defer parentPusher.Done()
		}
	}

	rp.ModifyResponse = func(r *http.Response) error {
		if childrenPusher != nil {
			defer wg.Done()
		}

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

		if !useFieldsHeader && !useFieldsQuery && !usePreloadHeader && !usePreloadQuery {
			// No reserved query parameters, nothing to do
			return nil
		}

		if !jsonRe.MatchString(r.Header.Get("Content-Type")) {
			// Not JSON, nothing to do
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
			if !u.IsAbs() && childrenPusher != nil {
				pushOptions := &http.PushOptions{Header: req.Header}
				pushOptions.Header.Del("Preload")
				pushOptions.Header.Del("Fields")
				pushOptions.Header.Del("Te") // Trailing headers aren't supported by Firefox for pushes, and we don't use them

				for _, pp := range n.strings(Preload, "") {
					if pp != "/" {
						pushOptions.Header.Add("Preload", pp)
					}
				}
				for _, fp := range n.strings(Preload, "") {
					if fp != "/" {
						pushOptions.Header.Add("Fields", fp)
					}
				}

				// HTTP/2, and relative relation, push!
				if err := childrenPusher.Push(uStr, pushOptions); err == nil {
					log.WithFields(log.Fields{"relation": uStr}).Debug("Relation pushed")
					return
				}
				log.WithFields(log.Fields{"relation": uStr, "reason": err}).Info("Failed to push")
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

		newBodyBuffer := bytes.NewBuffer(newBody)
		r.Body = ioutil.NopCloser(newBodyBuffer)
		r.Header["Content-Length"] = []string{fmt.Sprint(newBodyBuffer.Len())}

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		if childrenPusher != nil {
			defer wg.Done()
		}

		// Adapted from the default ErrorHandler
		log.Errorf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
	}
}

func (g *Gateway) retrievePushers(rw http.ResponseWriter, req *http.Request) (parentPusher *pusher, childrenPusher *pusher, requestID string) {
	internalPusher, ok := rw.(http.Pusher)
	if !ok {
		// Not an HTTP/2 connection
		return nil, nil, ""
	}

	// Need https://github.com/golang/go/issues/20566 to get rid of this hack
	parentID := req.Header.Get("Vulcain-Parent")
	if parentID != "" {
		parentPusher, ok = g.pushers.get(parentID)
		if ok {
			internalPusher = parentPusher.internalPusher
		} else {
			// Should not happen, is an attacker forging an evil request?
			log.WithFields(log.Fields{"uri": req.RequestURI, "parentID": parentID}).Debug("Pusher not found")
			parentID = ""
		}
	}

	// Store a new waitPusher to be used by children
	requestID = uuid.Must(uuid.NewV4()).String()
	req.Header.Set("Vulcain-Parent", requestID)

	childrenPusher = &pusher{internalPusher: internalPusher}
	g.pushers.add(requestID, childrenPusher)

	return parentPusher, childrenPusher, requestID
}

// NewGatewayFromEnv creates a gateway using the configuration set in env vars
func NewGatewayFromEnv() (*Gateway, error) {
	options, err := NewOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	return NewGateway(options), nil
}

// NewGateway creates a gateway
func NewGateway(options *Options) *Gateway {
	return &Gateway{
		options,
		nil,
		&pushers{pusherMap: make(map[string]*pusher)},
	}
}
