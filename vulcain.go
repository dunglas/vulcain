// Package vulcain helps implementing the Vulcain protocol (https://vulcain.rocks).
// It provides helper functions to parse HTTP requests containing "preload" and "fields" directives,
// to extract and push using HTTP/2 Server Push the relations of a JSON document matched by the "preload" directive,
// and to modify the JSON document according to both directives.
//
// The package can be used in conjunction with Go's httputil.ReverseProxy as well as in any HTTP handler.
package vulcain

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/dunglas/httpsfv"
	"github.com/getkin/kin-openapi/openapi3filter"
	"go.uber.org/zap"
)

var (
	jsonRe   = regexp.MustCompile(`(?i)\bjson\b`)
	preferRe = regexp.MustCompile(`\s*selector="?json-pointer"?`)
)

// Options allows to set the configuration
type Options struct {
	// OpenAPIFile contains the path to an OpenAPI definition (in YAML or JSON) documenting the relation between resources
	// This is useful only for non-hypermedia APIs
	OpenAPIFile string

	// MaxPushes is the maximum number of resources to push
	// Set MaxPushes to -1 to disable the limit
	MaxPushes int

	// Logger is the logger to use
	Logger *zap.Logger
}

// Vulcain parses and allow to modify the HTTP response according to the "preload" and "fields" directives extracted from the request
// Use New to create a Vulcain instance
type Vulcain struct {
	pushers *pushers
	openAPI *openAPI
	logger  *zap.Logger
}

// New creates a Vulcain instance
func New(options Options) *Vulcain {
	logger := options.Logger
	if options.Logger == nil {
		logger = zap.NewNop()
	}

	var o *openAPI
	if options.OpenAPIFile != "" {
		o = newOpenAPI(options.OpenAPIFile, logger)
	}

	return &Vulcain{
		&pushers{maxPushes: options.MaxPushes, pusherMap: make(map[string]*waitPusher), logger: logger},
		o,
		logger,
	}
}

// extractFromRequest extracts the "fields" and "preload" directives from the appropriate HTTP headers and query parameters
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

// getOpenAPIRoute gets the openapi3filter.Route instance corresponding to the given URL
func (v *Vulcain) getOpenAPIRoute(url *url.URL, route *openapi3filter.Route, routeTested bool) *openapi3filter.Route {
	if routeTested || v.openAPI == nil {
		return route
	}

	return v.openAPI.getRoute(url)
}

// CanApply checks is Vulcain's directives can be applied for this request and response
// CanApply must always be called
func (v *Vulcain) CanApply(rw http.ResponseWriter, req *http.Request, responseStatus int, responseHeaders http.Header) bool {
	pusher := v.pushers.getPusherForRequest(rw, req)

	// Not a success, or not JSON: don't modify the response
	if responseStatus < 200 ||
		responseStatus > 300 ||
		!jsonRe.MatchString(responseHeaders.Get("Content-Type")) {
		v.pushers.finish(req, pusher)

		return false
	}

	query := req.URL.Query()

	// No Vulcain hints: don't modify the response
	if req.Header.Get("Preload") == "" &&
		req.Header.Get("Fields") == "" &&
		query.Get("preload") == "" &&
		query.Get("fields") == "" {
		v.pushers.finish(req, pusher)

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

	v.pushers.finish(req, pusher)

	return false
}

// Apply pushes the requested relations and return rewrote response if necessary
// The content of the HTTP response is not automatically modified,
// it's the responsibility of the user to use the returned updated content.
// On the other hand headers are automatically updated.
// CanApply must always be called before Apply, or internal pushers will leak
func (v *Vulcain) Apply(req *http.Request, rw http.ResponseWriter, responseBody io.Reader, responseHeaders http.Header) ([]byte, error) {
	pusher := v.pushers.getPusherForRequest(rw, req)
	f, p, fieldsHeader, fieldsQuery, preloadHeader, preloadQuery := extractFromRequest(req)

	currentBody, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return nil, err
	}

	tree := &node{}
	tree.importPointers(preload, p)
	tree.importPointers(fields, f)

	var (
		oaRoute                         *openapi3filter.Route
		oaRouteTested, addPreloadToVary bool
	)
	newBody := v.traverseJSON(currentBody, tree, len(f) > 0, func(n *node, val string) string {
		var (
			u        *url.URL
			useOA    bool
			newValue string
		)

		oaRoute, oaRouteTested = v.getOpenAPIRoute(req.URL, oaRoute, oaRouteTested), true
		if u, useOA, err = v.parseRelation(n.String(), val, oaRoute); err != nil {
			return ""
		}

		// Never rewrite values when using OpenAPI, use header instead of query parameters
		if (preloadQuery || fieldsQuery) && !useOA {
			urlRewriter(u, n)
			newValue = u.String()
		}

		if n.preload {
			addPreloadToVary = !v.push(u, pusher, req, responseHeaders, n, preloadHeader, fieldsHeader)
		}

		return newValue
	})

	responseHeaders.Set("Content-Length", strconv.Itoa(len(newBody)))
	if fieldsHeader {
		responseHeaders.Add("Vary", "Fields")
	}
	if addPreloadToVary {
		responseHeaders.Add("Vary", "Preload")
	}

	v.pushers.finish(req, pusher)

	return newBody, nil
}

// addPreloadHeader sets preload Link rel=preload headers as fallback when Server Push isn't available (https://www.w3.org/TR/preload/)
func (v *Vulcain) addPreloadHeader(h http.Header, link string) {
	h.Add("Link", "<"+link+">; rel=preload; as=fetch")
	v.logger.Debug("link preload header added", zap.String("relation", link))
}

// push pushes a relation or add a Link rel=preload header as a fallback
// TODO: allow to set the nopush attribute using the configuration (https://www.w3.org/TR/preload/#server-push-http-2)
// TODO: send 103 early hints responses (https://tools.ietf.org/html/rfc8297)
func (v *Vulcain) push(u *url.URL, pusher *waitPusher, req *http.Request, newHeaders http.Header, n *node, preloadHeader, fieldsHeader bool) bool {
	url := u.String()
	if pusher == nil || u.IsAbs() {
		v.addPreloadHeader(newHeaders, url)
		return false
	}

	pushOptions := &http.PushOptions{Header: req.Header.Clone()}
	pushOptions.Header.Set(internalRequestHeader, pusher.id)
	pushOptions.Header.Del("Preload")
	pushOptions.Header.Del("Fields")
	pushOptions.Header.Del("Te") // Trailing headers aren't supported by Firefox for pushes, and we don't use them

	if preloadHeader {
		if preload := n.httpList(preload, ""); len(preload) > 0 {
			if v, err := httpsfv.Marshal(preload); err == nil {
				pushOptions.Header.Set("Preload", v)
			}
		}
	}
	if fieldsHeader {
		if f := n.httpList(fields, ""); len(f) > 0 {
			if v, err := httpsfv.Marshal(f); err == nil {
				pushOptions.Header.Set("Fields", v)
			}
		}
	}

	// HTTP/2, and relative relation, push!
	if err := pusher.Push(url, pushOptions); err != nil {
		// Don't add the preload header for something already pushed
		if errors.Is(err, errRelationAlreadyPushed) {
			return true
		}

		v.addPreloadHeader(newHeaders, url)
		v.logger.Debug("failed to push", zap.Stringer("node", n), zap.String("relation", url), zap.Error(err))

		return false
	}

	v.logger.Debug("relation pushed", zap.String("relation", url))
	return true
}

// parseRelation returns the URL of a relation, using OpenAPI to build it if necessary
func (v *Vulcain) parseRelation(selector, rel string, oaRoute *openapi3filter.Route) (*url.URL, bool, error) {
	var useOA bool
	if oaRoute != nil {
		if oaRel := v.openAPI.getRelation(oaRoute, selector, rel); oaRel != "" {
			rel = oaRel
			useOA = true
		}
	}

	u, err := url.Parse(rel)
	if err == nil {
		return u, useOA, nil
	}

	v.logger.Debug("the relation is an invalid URL", zap.String("node", selector), zap.String("relation", rel), zap.Error(err))

	return nil, useOA, err
}
