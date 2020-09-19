package vulcain

import (
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"go.uber.org/zap"
)

type openAPI struct {
	swagger *openapi3.Swagger
	router  *openapi3filter.Router
	logger  *zap.Logger
}

func newOpenAPI(file string, logger *zap.Logger) *openAPI {
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(file)
	if err != nil {
		panic(err)
	}

	return &openAPI{
		swagger,
		openapi3filter.NewRouter().WithSwagger(swagger),
		logger,
	}
}

func (o *openAPI) getRoute(url *url.URL) *openapi3filter.Route {
	route, _, err := o.router.FindRoute("GET", url)
	if err != nil {
		o.logger.Debug("route not found in the OpenAPI specification", zap.Stringer("url", url), zap.Error(err))
	}

	return route
}

// TODO: support operationRef in addition to operationId
func (o *openAPI) getRelation(r *openapi3filter.Route, selector, value string) string {
	for code, responseRef := range r.Operation.Responses {
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

func (o *openAPI) generateLink(operationID, parameter, value string) string {
	for path, i := range o.swagger.Paths {
		if op := i.GetOperation("GET"); op != nil && op.OperationID == operationID {
			return strings.ReplaceAll(path, "{"+parameter+"}", value)
		}
	}

	o.logger.Debug("peration not found in the OpenAPI specification", zap.String("operationID", operationID))

	return ""
}
