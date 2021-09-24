package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/charmixer/oas/api"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Metrics", Description:""},
	}
)

type GetMetricsRequest struct {}
type GetMetricsResponse struct {}

func (req GetMetricsRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).(http.HandlerFunc)
	t.ServeHTTP(w, r)
}
func (req GetMetricsRequest) Specification() (api.Path) {
	return api.Path{
		Summary: "Get metrics from the application",
		Description: `Get metrics from the application`,
		Tags: OPENAPI_TAGS,

		Request: api.Request{
			Description: `Request metrics`,
			Schema: GetMetricsRequest{},
		},

		Responses: []api.Response{{
			Description: `Metrics from prometheus`,
			Code: 200,
			ContentType: []string{"application/text"},
			Schema: GetMetricsResponse{},
		}},
	}
}
