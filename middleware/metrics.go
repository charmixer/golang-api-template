package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path", "method", "status"},
)

/*
var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)
*/

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path", "method"})

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// route := mux.CurrentRoute(r)
		// path, _ := route.GetPathTemplate()

		ctx := r.Context()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path, r.Method))

		wrapped := w.(*responseWriter)
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		//responseStatus.WithLabelValues(strconv.Itoa(wrapped.Status)).Inc()
		totalRequests.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(wrapped.Status)).Inc()

		timer.ObserveDuration()
	})
}

func init() {
	prometheus.Register(totalRequests)
	//prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
