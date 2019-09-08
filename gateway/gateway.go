package gateway

import (
	"bytes"
	"encoding/json"
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

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rp := httputil.NewSingleHostReverseProxy(g.Options.Upstream)
	defer rp.ServeHTTP(rw, req)

	p, explicitRequestID, explicitRequest := g.retrieveMainPusher(rw, req)

	rp.ModifyResponse = func(r *http.Response) error {
		if p != nil {
			if explicitRequest {
				defer p.Wait()
				defer g.pushers.remove(explicitRequestID)
			} else {
				defer p.Done()
			}
		}

		query := req.URL.Query()
		if len(query["preload"]) == 0 && len(query["fields"]) == 0 {
			// No reserved query parameters, nothing to do
			return nil
		}

		if !jsonRe.MatchString(r.Header.Get("Content-Type")) {
			// Not JSON, nothing to do
			return nil
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		var currentJSON interface{}
		if err := json.Unmarshal(body, &currentJSON); err != nil {
			// Invalid JSON
			return nil
		}

		var newJSON interface{}
		if len(query["fields"]) == 0 {
			newJSON = currentJSON
		} else {
			newJSON = g.traverseJSON("fields", query["fields"], currentJSON, nil, nil)
		}

		if len(query["preload"]) != 0 {
			pushOptions := &http.PushOptions{Header: req.Header}
			pushOptions.Header.Del("Te") // Trailing headers aren't supported by Firefox for pushes, and we don't use them
			newJSON = g.traverseJSON("preload", query["preload"], newJSON, newJSON, func(u *url.URL) {
				uStr := u.String()
				// TODO: allow to disable Server Push from the config
				if !u.IsAbs() && p != nil {
					// HTTP/2, and relative relation, push!
					if err := p.Push(uStr, pushOptions); err == nil {
						log.WithFields(log.Fields{"relation": uStr}).Debug("Relation pushed")
						return
					}
					log.WithFields(log.Fields{"relation": uStr, "reason": err}).Info("Failed to push")
				}

				log.WithFields(log.Fields{"relation": uStr}).Debug("Link preload header added")

				// Use preload Link headers as fallback (https://www.w3.org/TR/preload/)
				// TODO: allow to set the nopush attribute using the configuration (https://www.w3.org/TR/preload/#server-push-http-2)
				// TODO: send 103 early hints responses (https://tools.ietf.org/html/rfc8297)
				r.Header.Add("Link", "<"+uStr+">; rel=preload; as=fetch")
			})
		}

		// Construct the new JSON document by traversing the existing one
		newBodyContent, err := json.Marshal(newJSON)
		if err != nil {
			return err
		}

		newBody := bytes.NewBuffer(newBodyContent)
		r.Body = ioutil.NopCloser(newBody)
		r.Header["Content-Length"] = []string{fmt.Sprint(newBody.Len())}

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		if !explicitRequest {
			// Don't block the explicit request if there is an error in a push request
			p.Done()
		}

		// Adapted from the default ErrorHandler
		log.Errorf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
	}
}

func (g *Gateway) retrieveMainPusher(rw http.ResponseWriter, req *http.Request) (*pusher, string, bool) {
	mainPusher, ok := rw.(http.Pusher)
	if !ok {
		return nil, "", true
	}

	// Need https://github.com/golang/go/issues/20566 to get rid of this hack
	explicitRequestID := req.Header.Get("Vulcain-Explicit-Request-ID")
	if explicitRequestID == "" {
		// Explicit client-initiated request
		p := &pusher{internalPusher: mainPusher}

		explicitRequestID := uuid.Must(uuid.NewV4()).String()
		req.Header.Add("Vulcain-Explicit-Request-ID", explicitRequestID)

		g.pushers.add(explicitRequestID, p)

		return p, explicitRequestID, true
	}

	// Push request
	p, _ := g.pushers.get(explicitRequestID)
	if p == nil {
		log.WithFields(log.Fields{"uri": req.RequestURI, "explicitRequestID": explicitRequestID}).Debug("Pusher not found")

		return nil, "", true
	}

	return p, explicitRequestID, false
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
