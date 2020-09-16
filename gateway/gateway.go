package gateway

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/dunglas/httpsfv"
	"github.com/getkin/kin-openapi/openapi3filter"
	log "github.com/sirupsen/logrus"
)

var (
	jsonRe   = regexp.MustCompile(`(?i)\bjson\b`)
	preferRe = regexp.MustCompile(`\s*selector="?json-pointer"?`)
)

// Gateway is the main struct
type Gateway struct {
	options *Options
	server  *http.Server
	pushers *pushers
	openAPI *openAPI
}

func extractFromRequest(req *http.Request) (fields, preload httpsfv.List, fieldsHeader, fieldsQuery, preloadHeader, preloadQuery bool) {
	query := req.URL.Query()
	var err error
	if len(req.Header["Fields"]) > 0 {
		if fields, err = httpsfv.UnmarshalList(req.Header["Fields"]); err == nil {
			fieldsHeader = true
		}
	}

	if !fieldsHeader && len(query["fields"]) > 0 {
		if fields, err = httpsfv.UnmarshalList(query["fields"]); err == nil {
			fieldsQuery = true
		}
	}

	if len(req.Header["Preload"]) > 0 {
		if preload, err = httpsfv.UnmarshalList(req.Header["Preload"]); err == nil {
			preloadHeader = true
		}
	}

	if !preloadHeader && len(query["preload"]) > 0 {
		if preload, err = httpsfv.UnmarshalList(query["preload"]); err == nil {
			preloadQuery = true
		}
	}

	return fields, preload, fieldsHeader, fieldsQuery, preloadHeader, preloadQuery
}

func (g *Gateway) getOpenAPIRoute(url *url.URL, route *openapi3filter.Route, routeTested bool) *openapi3filter.Route {
	if routeTested || g.openAPI == nil {
		return route
	}

	return g.openAPI.getRoute(url)
}

func canParse(responseHeaders http.Header, req *http.Request, fields, preload httpsfv.List) bool {
	if (len(fields) == 0 && len(preload) == 0) || !jsonRe.MatchString(responseHeaders.Get("Content-Type")) {
		// No Vulcain hints, or not JSON: don't modify the response
		return false
	}

	prefers, ok := req.Header["Prefer"]
	if !ok {
		return true
	}

	for _, p := range prefers {
		if preferRe.MatchString(p) {
			return true
		}
	}

	return false
}

func (g *Gateway) Apply(req *http.Request, rw http.ResponseWriter, responseBody io.Reader, responseHeaders http.Header) ([]byte, http.Header, error) {
	pusher := g.pushers.getPusherForRequest(rw, req)
	fields, preload, fieldsHeader, fieldsQuery, preloadHeader, preloadQuery := extractFromRequest(req)
	if !canParse(responseHeaders, req, fields, preload) {
		g.pushers.cleanupAfterRequest(req, pusher, false)
		return nil, nil, nil
	}

	currentBody, err := ioutil.ReadAll(responseBody)
	if err != nil {
		g.pushers.cleanupAfterRequest(req, pusher, false)
		return nil, nil, err
	}

	tree := &node{}
	tree.importPointers(Preload, preload)
	tree.importPointers(Fields, fields)

	var (
		oaRoute                         *openapi3filter.Route
		oaRouteTested, addPreloadToVary bool
	)
	newBody := traverseJSON(currentBody, tree, len(fields) > 0, func(n *node, v string) string {
		var (
			u        *url.URL
			useOA    bool
			newValue string
		)

		oaRoute, oaRouteTested = g.getOpenAPIRoute(req.URL, oaRoute, oaRouteTested), true
		if u, useOA, err = g.parseRelation(n.String(), v, oaRoute); err != nil {
			return ""
		}

		// Never rewrite values when using OpenAPI, use header instead of query parameters
		if (preloadQuery || fieldsQuery) && !useOA {
			urlRewriter(u, n)
			newValue = u.String()
		}

		if len(preload) > 0 {
			addPreloadToVary = !g.push(u, pusher, req, &responseHeaders, n, preloadHeader, fieldsHeader)
		}

		return newValue
	})

	if fieldsHeader {
		responseHeaders.Add("Vary", "Fields")
	}
	if addPreloadToVary {
		responseHeaders.Add("Vary", "Preload")
	}

	g.pushers.cleanupAfterRequest(req, pusher, true)

	return newBody, responseHeaders, nil
}

func (g *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rp := httputil.NewSingleHostReverseProxy(g.options.Upstream)
	rp.ModifyResponse = func(resp *http.Response) error {
		newBody, newHeaders, err := g.Apply(req, rw, resp.Body, resp.Header)
		if newBody == nil {
			return err
		}

		newBodyBuffer := bytes.NewBuffer(newBody)
		resp.Body = ioutil.NopCloser(newBodyBuffer)
		resp.Header = newHeaders
		resp.Header["Content-Length"] = []string{fmt.Sprint(newBodyBuffer.Len())}

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		// Adapted from the default ErrorHandler
		log.Errorf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
	}

	proto := "https"
	if req.TLS == nil {
		proto = "http"
	}

	req.Header.Set("X-Forwarded-Proto", proto)
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Header.Del("X-Forwarded-For")
	rp.ServeHTTP(rw, req)
}

// addPreloadHeader sets preload Link headers as fallback when Server Push isn't available (https://www.w3.org/TR/preload/)
func addPreloadHeader(h *http.Header, link string) {
	h.Add("Link", "<"+link+">; rel=preload; as=fetch")
	log.WithFields(log.Fields{"relation": link}).Debug("Link preload header added")
}

// TODO: allow to set the nopush attribute using the configuration (https://www.w3.org/TR/preload/#server-push-http-2)
// TODO: send 103 early hints responses (https://tools.ietf.org/html/rfc8297)
func (g *Gateway) push(u *url.URL, pusher *waitPusher, req *http.Request, newHeaders *http.Header, n *node, preloadHeader, fieldsHeader bool) bool {
	url := u.String()
	if pusher == nil || u.IsAbs() {
		addPreloadHeader(newHeaders, url)
		return false
	}

	pushOptions := &http.PushOptions{Header: req.Header.Clone()}
	pushOptions.Header.Set(internalRequestHeader, pusher.id)
	pushOptions.Header.Del("Preload")
	pushOptions.Header.Del("Fields")
	pushOptions.Header.Del("Te") // Trailing headers aren't supported by Firefox for pushes, and we don't use them

	if preloadHeader {
		if preload := n.httpList(Preload, ""); len(preload) > 0 {
			if v, err := httpsfv.Marshal(preload); err == nil {
				pushOptions.Header.Set("Preload", v)
			}
		}
	}
	if fieldsHeader {
		if fields := n.httpList(Fields, ""); len(fields) > 0 {
			if v, err := httpsfv.Marshal(fields); err == nil {
				pushOptions.Header.Set("Fields", v)
			}
		}
	}

	// HTTP/2, and relative relation, push!
	if err := pusher.Push(url, pushOptions); err != nil {
		// Don't add the preload header for something already pushed
		if _, ok := err.(*relationAlreadyPushedError); ok {
			return true
		}

		addPreloadHeader(newHeaders, url)
		log.WithFields(log.Fields{
			"node":     n.String(),
			"relation": url,
			"reason":   err,
		}).Debug("Failed to push")

		return false
	}

	log.WithFields(log.Fields{"relation": url}).Debug("Relation pushed")
	return true
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
	var o *openAPI
	if options.OpenAPIFile != "" {
		o = newOpenAPI(options.OpenAPIFile)
	}

	return &Gateway{
		options,
		nil,
		&pushers{maxPushes: options.MaxPushes, pusherMap: make(map[string]*waitPusher)},
		o,
	}
}

func (g *Gateway) parseRelation(selector, rel string, oaRoute *openapi3filter.Route) (*url.URL, bool, error) {
	var useOA bool
	if oaRoute != nil {
		if oaRel := g.openAPI.getRelation(oaRoute, selector, rel); oaRel != "" {
			rel = oaRel
			useOA = true
		}
	}

	u, err := url.Parse(rel)
	if err == nil {
		return u, useOA, nil
	}

	log.WithFields(
		log.Fields{
			"node":     selector,
			"relation": rel,
			"reason":   err,
		}).Debug("The relation is an invalid URL")

	return nil, useOA, err
}
