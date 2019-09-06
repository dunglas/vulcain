package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var jsonRe = regexp.MustCompile(`(?i)\bjson\b`)

// Gateway is the main struct
type Gateway struct {
	options      *Options
	reverseProxy *httputil.ReverseProxy
	server       *http.Server
}

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
		newJSON = g.traverseJSON("preload", query["preload"], newJSON, newJSON, func(relation string) {
			log.Infof("Push Me! %s", relation)
		})
	}

	// Construct the new JSON document by traversing the existing one
	body, err := json.Marshal(newJSON)
	if err != nil {
		return
	}
	log.Printf("URL: %s, Body %s", req.RequestURI, body)

	brw.body = body

	/*if req.RequestURI == "/books.jsonld" {
		if pusher, ok := rw.(http.Pusher); ok {
			// Push is supported.
			if err := pusher.Push("/books/1.jsonld", &http.PushOptions{Header: req.Header}); err != nil {
				log.Printf("Failed to push: %v", err)
			} else {
				log.Println("Pushed!")
			}
		}
	}*/
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
	rp.ModifyResponse = func(r *http.Response) error {
		r.Header.Add("X-Push", "On")
		return nil
	}

	return &Gateway{
		options,
		rp,
		nil,
	}
}
