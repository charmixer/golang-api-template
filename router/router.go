package router

import (
	"net/http"

	_ "github.com/charmixer/golang-api-template/endpoints/metrics"
	"github.com/charmixer/golang-api-template/endpoints/health"
	"github.com/charmixer/golang-api-template/endpoints/openapi"

	"github.com/charmixer/oas/api"

	"github.com/julienschmidt/httprouter"
	// "github.com/gorilla/mux"
)

type Route interface {
	http.Handler
	Specification() api.Path
}

type Router struct {
	Mux *httprouter.Router
	OpenAPI api.Api
}

func (r *Router) NewRoute(method string, uri string, route Route) {
	r.OpenAPI.NewEndpoint(method, uri, route.Specification())
	r.Mux.Handler(method, uri, route)
}

func NewRouter() (Router) {
	r := Router{
		OpenAPI: api.Api{
			Title: "Golang api template",
			Description: `Gives a simple blueprint for creating new api's`,
			Version: "0.0.0",
		},
		Mux: httprouter.New(),
	}

	// definere din routes og informere oas
	//oas.NewEndpoint("GET",  "/metrics", metrics.GetMetricsSpec(), []api.Tag{healthTag})
	//r.Handler("GET",  "/metrics", metrics.GetMetrics)

	// NewRoute(r, "GET", "/metrics", metrics.GetMetrics{})
	r.NewRoute("GET", "/health", health.PostHealthRequest{})

	r.NewRoute("GET", "/docs", openapi.GetDocsRequest{})
	r.NewRoute("GET", "/docs/openapi", openapi.GetOpenapiRequest{})


	/*
	//	docsTag := api.Tag{Name: "Docs", Description: "Documentation stuff"}
		healthTag := api.Tag{Name: "Health", Description: "Health stuff"}

		oas.NewPath("GET",  "/docs", openapi.GetOpenapiDocs, openapi.GetOpenapiDocsSpec(), []api.Tag{docsTag})
		oas.NewPath("GET",  "/docs/openapi", openapi.GetOpenapi, openapi.GetOpenapiSpec(), []api.Tag{docsTag})

		oas.NewPath("GET",  "/health", health.GetHealth, health.GetHealthSpec(), []api.Tag{healthTag})
		oas.NewPath("POST", "/health", health.PostHealth, health.PostHealthSpec(), []api.Tag{healthTag})

		//oas.NewPath("GET",  "/metrics", metrics.GetMetrics, metrics.GetMetricsSpec(), []api.Tag{healthTag})

		a := health.PostHealth{}

		oas.NewEndpoint("POST",  "/health", a.Specification(), []api.Tag{healthTag})
		//r.Handler()
	*/

	//oas.NewEndpoint("GET",  "/metrics", h.Specification(), []api.Tag{healthTag})
	//r.Handler("GET",  "/metrics", h)


	/*for _,e := range oas.Paths {
		r.Handler(e.Method, e.Url, health.PostHealth{})
	}*/

	/*
	# github.com/charmixer/golang-api-template/router
	router/router.go:34:55: cannot use health.PostHealth{} (type health.PostHealth) as type *api.Endpoint in argument to oas.NewEndpoint:
		*api.Endpoint is pointer to interface, not interface
	router/router.go:43:40: invalid type assertion: e.Endpoint.(health.PostHealth) (non-interface type *api.Endpoint on left)
	*/

	return r
}
/*
func NewRouter(oas api.Api) (r *mux.Router) {
	r = mux.NewRouter()

	for _,p := range oas.GetPaths() {
		r.HandleFunc(p.Url, p.Handler).Methods(p.Method)
	}

	return r
}
*/
