package docs

import (
	"fmt"
	"net/http"

	"github.com/charmixer/oas/api"
	"github.com/charmixer/oas/exporter"
	"github.com/charmixer/golang-api-template/app"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/middleware"

	"go.opentelemetry.io/otel"

	_ "github.com/rs/zerolog/log"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Documentation", Description:""},
	}
)

type GetOpenapiRequest struct {
	Format string `json:"format" oas-query:"format" oas-desc:"Format returned by the endpoint, eg. json"`
}

// https://golang.org/doc/effective_go#embedding
type GetOpenapiEndpoint struct {
	endpoint.Endpoint
	Request GetOpenapiRequest
	Response exporter.Openapi
	responseType string
}
func (ep *GetOpenapiEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("request")
	ctx, span := tr.Start(ctx, fmt.Sprintf("%s execution", r.URL.Path))
	defer span.End()

	t := r.URL.Query().Get("format")

	ep.Response = app.Env.OpenAPI

  if t == "json" {
		w.Header().Set("Content-Type", "application/json")
		ep.responseType = "json"
  } else {
		w.Header().Set("Content-Type", "text/plain; application/yaml; charset=utf-8")
		ep.responseType = "yaml"
	}
}

func NewGetOpenapiEndpoint() (endpoint.EndpointHandler) {
	ep := GetOpenapiEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary: "OpenAPI specification",
			Description: ``,
			Tags: OPENAPI_TAGS,

			Request: api.Request{
				Description: ``,
				Schema: GetOpenapiRequest{},
			},

			Responses: []api.Response{{
				Description: `Returns openapi spec in given format`,
				Code: 200,
				ContentType: []string{"application/json", "application/yaml"},
				Schema: exporter.Openapi{},
			},/*{
				Description: `error ...`,
				Code: 400,
				ContentType: []string{"application/json"},
				Schema: GetOpenapiEndpoint.BadRequest{},
			}*/},
		}),

		endpoint.WithMiddleware(
			middleware.WithRequestParser(&ep.Request),
			middleware.WithRequestValidation(&ep.Request/*, &ep.BadRequest*/),

			middleware.WithResponseWriter(&ep.responseType, &ep.Response),
		),
	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return &ep
}
