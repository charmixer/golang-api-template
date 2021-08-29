package metrics

import (
	"net/http"

	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type GetMetricsRequest struct {

}
type GetMetricsResponse struct {
	cpu_usage int
	mem_usage int
}


func GetMetrics() (http.HandlerFunc) {
	return promhttp.Handler().(http.HandlerFunc)
}
