package docs

import (
	"net/http"

	"github.com/charmixer/oas/api"
	"github.com/charmixer/oas/exporter"
	"github.com/charmixer/golang-api-template/app"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/middleware"


	_ "github.com/rs/zerolog/log"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Documentation", Description:""},
	}
)

type Test struct {
	Something string `oas-desc:"Format returned by the endpoint, eg. json"`
}

type GetOpenapiRequest struct {
	Mode string `oas-query:"mode"`
	Format string `oas-query:"format" oas-desc:"Format returned by the endpoint, eg. json"`
	XOverrideMethodHeader string `oas-header:"x-override-method-header" oas-desc:"My description goes here"`
	Debug string `oas-cookie:"debug"`
	Test Test `oas-query:"test"`
}

// https://golang.org/doc/effective_go#embedding
type GetOpenapiEndpoint struct {
	endpoint.Endpoint
	Request GetOpenapiRequest
	Response exporter.Openapi
}
func (ep *GetOpenapiEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("format")

	ep.Response = app.Env.OpenAPI

  if t == "json" {
		w.Header().Set("Content-Type", "application/json")
  } else {
		w.Header().Set("Content-Type", "application/x-yaml")
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
				ContentType: []string{"application/json", "application/x-yaml"},
				Schema: exporter.Openapi{},
			}},
		}),

		endpoint.WithMiddleware(
			middleware.WithRequestParser(&ep.Request),
			middleware.WithResponseWriter(&ep.Response),
		),
	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return &ep
}
