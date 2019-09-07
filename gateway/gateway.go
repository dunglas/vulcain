package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

var jsonRe = regexp.MustCompile(`(?i)\bjson\b`)

// Gateway is the main struct
type Gateway struct {
	options      *Options
	reverseProxy *httputil.ReverseProxy
	server       *http.Server
	pushers      *pushers
}

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var p *pusher
	if mainPusher, ok := rw.(http.Pusher); ok {
		// We'll be able to get rid of this hack when https://github.com/golang/go/issues/20566 will be resolved
		explicitRequestID := req.Header.Get("Vulcain-Explicit-Request-ID")
		if explicitRequestID == "" {
			// Explicit client-initiated request
			p = &pusher{internalPusher: mainPusher}
			defer p.Wait()

			explicitRequestID = uuid.Must(uuid.NewV4()).String()
			req.Header.Add("Vulcain-Explicit-Request-ID", explicitRequestID)

			g.pushers.add(explicitRequestID, p)
			defer g.pushers.remove(explicitRequestID)
		} else {
			// Push request
			p, _ = g.pushers.get(explicitRequestID)
			defer p.Done()
		}
	}

	query := req.URL.Query()
	if len(query["preload"]) == 0 && len(query["fields"]) == 0 {
		// No reserved query parameters, don't buffer
		g.reverseProxy.ServeHTTP(rw, req)
		return
	}

	// I'm not fond of this... But I'll live with it until https://github.com/golang/go/issues/19307 is resolved!
	if accept := req.Header.Get("Accept"); accept != "" && !strings.Contains(accept, "*/*") && !jsonRe.MatchString(accept) {
		// Accept header doesn't include a JSON MIME type, don't buffer
		g.reverseProxy.ServeHTTP(rw, req)
		return
	}

	// Assume the response will be a JSON document, buffer the request
	brw := newBufferedResponseWriter(rw)
	defer brw.send()
	g.reverseProxy.ServeHTTP(brw, req)

	if !jsonRe.MatchString(brw.Header().Get("Content-Type")) {
		// Ignore non-JSON documents
		return
	}

	var currentJSON interface{}
	if err := json.Unmarshal(brw.bodyContent(), &currentJSON); err != nil {
		// Invalid JSON
		return
	}

	var newJSON interface{}
	if len(query["fields"]) == 0 {
		newJSON = currentJSON
	} else {
		newJSON = g.traverseJSON("fields", query["fields"], currentJSON, nil, nil)
	}

	if len(query["preload"]) != 0 {
		pushOptions := &http.PushOptions{Header: req.Header}
		newJSON = g.traverseJSON("preload", query["preload"], newJSON, newJSON, func(u *url.URL) {
			uStr := u.String()
			// TODO: allow to disable Server Push from the config
			if !u.IsAbs() && p != nil {
				// HTTP/2, and relative relation, push!
				if err := p.Push(uStr, pushOptions); err == nil {
					log.WithFields(log.Fields{"relation": uStr}).Info("Relation pushed")
					return
				}
			}

			log.Info("Add header")

			// Use preload Link headers as fallback (https://www.w3.org/TR/preload/)
			// TODO: allow to set the nopush attribute using the configuration (https://www.w3.org/TR/preload/#server-push-http-2)
			// TODO: send 103 early hints responses (https://tools.ietf.org/html/rfc8297)
			brw.Header().Add("Link", "<"+uStr+">; rel=preload; as=fetch")
		})
	}

	// Construct the new JSON document by traversing the existing one
	body, err := json.Marshal(newJSON)
	if err != nil {
		return
	}
	log.Printf("URL: %s, Body %s", req.RequestURI, body)

	brw.body = body
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
	rp := httputil.NewSingleHostReverseProxy(options.Upstream)

	return &Gateway{
		options,
		rp,
		nil,
		&pushers{pusherMap: make(map[string]*pusher)},
	}
}
