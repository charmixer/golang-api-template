package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	t := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}).(http.HandlerFunc)
	t.ServeHTTP(w, r)
}
