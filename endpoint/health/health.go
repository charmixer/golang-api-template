package health

import (
	"net/http"
	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/middleware"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Health", Description:"Endpoints reporting the health of the application"},
	}
)

type GetHealthRequest struct {}
type GetHealthResponse struct {
	Alive bool `json:"alive_json" xml:"alive_xml" desc:"Tells if bla"`
	Ready bool `json:"ready_json" xml:"ready_xml"`
}

// https://golang.org/doc/effective_go#embedding
type GetHealthEndpoint struct {
	endpoint.Endpoint
	Request GetHealthRequest
	Response GetHealthResponse
}
func (ep *GetHealthEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ep.Response = GetHealthResponse{
		Alive: true,
	}
}

func NewGetHealthEndpoint() (endpoint.EndpointHandler) {
	ep := GetHealthEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary: "Test 2",
			Description: `Testing 2`,
			Tags: OPENAPI_TAGS,

			Request: api.Request{
				Description: `Testing Request`,
				Schema: GetHealthRequest{},
			},

			Responses: []api.Response{{
				Description: `Testing OK Response`,
				Code: 200,
				Schema: GetHealthResponse{},
			}},
		}),

		endpoint.WithMiddleware(
			middleware.WithRequestParser(&ep.Request),
			middleware.WithRequestValidation(&ep.Request),

			middleware.WithResponseValidation(&ep.Response),
			middleware.WithResponseWriter(&ep.Response),
		),

	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return &ep
}
