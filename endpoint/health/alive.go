package health

import (
	"fmt"
	"net/http"

	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/endpoint/problem"

	"go.opentelemetry.io/otel"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Health", Description: ""},
	}
)

// https://golang.org/doc/effective_go#embedding
type GetHealthAliveEndpoint struct {
	endpoint.Endpoint
}

func (ep GetHealthAliveEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer("request").Start(r.Context(), fmt.Sprintf("%s handler", r.URL.Path))
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func NewGetHealthAliveEndpoint() endpoint.EndpointHandler {
	ep := GetHealthAliveEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary:     "Get health information about the service",
			Description: ``,
			Tags:        OPENAPI_TAGS,

			Request: api.Request{
				Description: ``,
				//Schema:      GetHealthReadyRequest{},
			},

			Responses: []api.Response{{
				Description: http.StatusText(http.StatusOK),
				Code:        http.StatusOK,
				//Schema:      GetHealthReadyResponse{},
			}, {
				Description: http.StatusText(http.StatusServiceUnavailable),
				Code:        http.StatusServiceUnavailable,
				Schema:      problem.ProblemDetails{},
			}},
		}),
	)

	return ep
}
