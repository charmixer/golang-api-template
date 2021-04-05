package metrics

import (
	"net/http"
)

type GetMetricsRequest struct {

}
type GetMetricsResponse struct {
	cpu_usage int
	mem_usage int
}


func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get /metrics\n"))
}
