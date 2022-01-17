package docs

import (
	"fmt"
	"net/http"

	"github.com/charmixer/golang-api-template/env"
	"github.com/charmixer/oas/api"
	"github.com/charmixer/oas/exporter"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/endpoint/problem"

	"go.opentelemetry.io/otel"

	_ "github.com/rs/zerolog/log"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Documentation", Description: ""},
	}
)

type GetOpenapiRequest struct {
	Format string `json:"format" query:"format" query:"format" description:"Format returned by the endpoint, eg. json"`
}
type GetOpenapiResponse exporter.Openapi

// https://golang.org/doc/effective_go#embedding
type GetOpenapiEndpoint struct {
	endpoint.Endpoint
}

func (ep GetOpenapiEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("request").Start(r.Context(), fmt.Sprintf("%s handler", r.URL.Path))
	defer span.End()

	request := GetOpenapiRequest{}
	if err := endpoint.WithRequestQueryParser(ctx, r, &request); err != nil {
		problem.MustWrite(w, err)
		return
	}

	response := env.Env.OpenAPI

	responseType := ""
	if request.Format == "json" {
		w.Header().Set("Content-Type", "application/json")
		responseType = "json"
	} else {
		w.Header().Set("Content-Type", "text/plain; application/yaml; charset=utf-8")
		responseType = "yaml"
	}

	if err := endpoint.WithResponseWriter(ctx, w, responseType, response); err != nil {
		problem.MustWrite(w, err)
		return
	}
}

func NewGetOpenapiEndpoint() endpoint.EndpointHandler {
	ep := GetOpenapiEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary:     "OpenAPI specification",
			Description: ``,
			Tags:        OPENAPI_TAGS,

			Request: api.Request{
				Description: ``,
				Schema:      GetOpenapiRequest{},
			},

			Responses: []api.Response{{
				Description: `Returns openapi spec in given format`,
				Code:        200,
				ContentType: []string{"application/json", "application/yaml"},
				Schema:      exporter.Openapi{},
			}, /*{
				Description: `error ...`,
				Code: 400,
				ContentType: []string{"application/json"},
				Schema: GetOpenapiEndpoint.BadRequest{},
			}*/},
		}),
	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return ep
}
