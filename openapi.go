package vulcain

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"go.uber.org/zap"
)

// openAPI is used to find the URL of a relation using an OpenAPI description
type openAPI struct {
	swagger *openapi3.T
	router  routers.Router
	logger  *zap.Logger
}

// newOpenAPI creates a ne openAPI instance
func newOpenAPI(file string, logger *zap.Logger) *openAPI {
	swagger, err := openapi3.NewLoader().LoadFromFile(file)
	if err != nil {
		panic(err)
	}

	router, err := legacy.NewRouter(swagger)
	if err != nil {
		panic(err)
	}

	return &openAPI{
		swagger,
		router,
		logger,
	}
}

// getRoute gets the routers.Route instance related to the given URL
func (o *openAPI) getRoute(url *url.URL) *routers.Route {
	route, _, err := o.router.FindRoute(&http.Request{Method: "GET", URL: url})
	if err != nil {
		o.logger.Debug("route not found in the OpenAPI specification", zap.Stringer("url", url), zap.Error(err))
	}

	return route
}

// getRelation generated the link for the given parameters
// TODO: support operationRef in addition to operationId
func (o *openAPI) getRelation(r *routers.Route, selector, value string) string {
	for code, responseRef := range r.Operation.Responses.Map() {
		if (!strings.HasPrefix(code, "2")) || responseRef.Value == nil {
			continue
		}

		if rel := o.generateLinkForResponse(responseRef.Value, selector, value); rel != "" {
			return rel
		}
	}

	// Fallback on the default response
	if d := r.Operation.Responses.Default(); d != nil && d.Value != nil {
		if rel := o.generateLinkForResponse(d.Value, selector, value); rel != "" {
			return rel
		}
	}

	o.logger.Error("openAPI Link not found (using operationRef isn't supported yet)")

	return ""
}

// generateLinkForResponse uses the openapi3.Response extracted from the OpenAPI description to generate a URL
func (o *openAPI) generateLinkForResponse(response *openapi3.Response, selector, value string) string {
	for _, linkRef := range response.Links {
		if linkRef == nil || linkRef.Value == nil {
			continue
		}

		var parameter string
		for p, s := range linkRef.Value.Parameters {
			if s == "$response.body#"+selector {
				parameter = p
				break
			}
		}

		if parameter != "" && linkRef.Value.OperationID != "" {
			return o.generateLink(linkRef.Value.OperationID, parameter, value)
		}
	}

	return ""
}

// generateLink uses the template IRI extracted from the OpenAPI description to generate a URL
func (o *openAPI) generateLink(operationID, parameter, value string) string {
	for path, i := range o.swagger.Paths.Map() {
		if op := i.GetOperation("GET"); op != nil && op.OperationID == operationID {
			return strings.ReplaceAll(path, "{"+parameter+"}", value)
		}
	}

	o.logger.Debug("operation not found in the OpenAPI specification", zap.String("operationID", operationID))

	return ""
}
