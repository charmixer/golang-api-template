package router

import (
	"net/http"

	"github.com/charmixer/golang-api-template/endpoints/metrics"
	"github.com/charmixer/golang-api-template/endpoints/health"
	"github.com/charmixer/golang-api-template/endpoints/openapi"

	"github.com/charmixer/golang-api-template/middleware"

	"github.com/charmixer/oas/api"

	"github.com/julienschmidt/httprouter"
)

type Route interface {
	http.Handler
	Specification() api.Path
}

type Router struct {
	httprouter.Router
	OpenAPI api.Api
	Middleware []middleware.MiddlewareHandler
}

func (r *Router) NewRoute(method string, uri string, route Route, handlers ...middleware.MiddlewareHandler) {
	r.OpenAPI.NewEndpoint(method, uri, route.Specification())
	r.Handler(method, uri, middleware.New(route, handlers...))
}
func (r *Router) Use(h ...middleware.MiddlewareHandler) {
	r.Middleware = append(r.Middleware, h...)
}
func (r *Router) Handle() http.Handler {
	return middleware.New(r, r.Middleware...)
}

func NewRouter(appName string) (*Router) {
	r := &Router{
		OpenAPI: api.Api{
			Title: "Golang api template",
			Description: `Gives a simple blueprint for creating new api's`,
			Version: "0.0.0",
		},
	}

	// Ordering matters
	r.Use(
		middleware.WithResponseWriter(),
		middleware.WithContext(),
		middleware.WithTracing(appName),
		middleware.WithMetrics(),
		middleware.WithLogging(),
	)

  /*healthEndpoint := health.HealthEndpoint{
		Method: "GET",
		Path: "/health",

	}
	healthRequest := health.GetHealthRequest{}
	r.NewRoute("GET", "/health", &healthEndpoint,
		middleware.WithInputValidation(healthEndpoint),
		middleware.WithOutputValidation(healthEndpoint),
	)*/

	r.NewRoute("GET", "/health", health.GetHealthRequest{})

	r.NewRoute("GET", "/docs", openapi.GetDocsRequest{})
	r.NewRoute("GET", "/docs/openapi", openapi.GetOpenapiRequest{})

	r.NewRoute("GET", "/metrics", metrics.GetMetricsRequest{})

	return r
}
