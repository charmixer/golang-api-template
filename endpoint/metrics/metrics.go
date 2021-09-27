package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/middleware"

	"github.com/charmixer/oas/api"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Metrics", Description:""},
	}
)

type GetMetricsRequest struct {}
type GetMetricsResponse struct {}

// https://golang.org/doc/effective_go#embedding
type GetMetricsEndpoint struct {
	endpoint.Endpoint
	Request GetMetricsRequest
	Response GetMetricsResponse
}
func (ep *GetMetricsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).(http.HandlerFunc)
	t.ServeHTTP(w, r)
}

func NewGetMetricsEndpoint() (endpoint.EndpointHandler) {
	ep := GetMetricsEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary: "Get metrics from the application",
			Description: `Get metrics from the application`,
			Tags: OPENAPI_TAGS,

			Request: api.Request{
				Description: `Request metrics`,
				//Schema: GetMetricsRequest{},
			},

			Responses: []api.Response{{
				Description: `Metrics from prometheus`,
				Code: 200,
				ContentType: []string{"application/text"},
				//Schema: GetMetricsResponse{},
			}},
		}),

		endpoint.WithMiddleware(
			middleware.WithRequestParser(&ep.Request),
		),
	)

	// Must be pointer to allow ServeHTTP method to be used with *Endpoint
	return &ep
}