package router

import (
	_ "net/http"

	"github.com/charmixer/golang-api-template/endpoints/metrics"
	"github.com/charmixer/golang-api-template/endpoints/health"
	"github.com/charmixer/golang-api-template/endpoints/openapi"

	"github.com/charmixer/oas/api"

	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	// "github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
)

func NewOas() (oas api.Api){
	oas = api.Api{}

	oas.Title = "Golang api template"
	oas.Description = `Gives a simple blueprint for creating new api's`
	oas.Version = "0.0.0"

	docsTag := api.Tag{Name: "Docs", Description: "Documentation stuff"}
	healthTag := api.Tag{Name: "Health", Description: "Health stuff"}

	oas.NewPath("GET",  "/docs", openapi.GetOpenapiDocs, openapi.GetOpenapiDocsSpec(), []api.Tag{docsTag})
	oas.NewPath("GET",  "/docs/openapi.yaml", openapi.GetOpenapi, openapi.GetOpenapiSpec(), []api.Tag{docsTag})

	oas.NewPath("GET",  "/health", health.GetHealth, health.GetHealthSpec(), []api.Tag{healthTag})
	oas.NewPath("POST", "/health", health.PostHealth, health.PostHealthSpec(), []api.Tag{healthTag})

	oas.NewPath("GET",  "/metrics", metrics.GetMetrics(), health.GetHealthSpec(), []api.Tag{healthTag})

	return oas
}


func NewRouter(oas api.Api) (r *mux.Router) {
	r = mux.NewRouter()

	for _,p := range oas.GetPaths() {
		r.HandleFunc(p.Url, p.Handler).Methods(p.Method)
	}

	return r
}
