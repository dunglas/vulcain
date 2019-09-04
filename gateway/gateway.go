package gateway

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

// Gateway is the main struct
type Gateway struct {
	options      *Options
	reverseProxy *httputil.ReverseProxy
	server       *http.Server
}

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var buffer bytes.Buffer
	brw := &bufferedResponseWriter{rw, &buffer}
	g.reverseProxy.ServeHTTP(brw, req)
	respBody := buffer.Bytes()
	log.Printf("URL: %s, Body %s", req.RequestURI, respBody)

	if req.RequestURI == "/books.jsonld" {
		if pusher, ok := rw.(http.Pusher); ok {
			// Push is supported.
			if err := pusher.Push("/books/1.jsonld", &http.PushOptions{Header: req.Header}); err != nil {
				log.Printf("Failed to push: %v", err)
			} else {
				log.Println("Pushed!")
			}
		}
	}

	io.Copy(rw, bytes.NewReader(respBody))
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
