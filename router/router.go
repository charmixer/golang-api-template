package router

import (
	"net/http"

	"github.com/charmixer/golang-api-template/endpoint"

	"github.com/charmixer/golang-api-template/endpoint/metrics"
	"github.com/charmixer/golang-api-template/endpoint/health"
	"github.com/charmixer/golang-api-template/endpoint/docs"

	"github.com/charmixer/golang-api-template/middleware"

	"github.com/charmixer/oas/api"

	"github.com/julienschmidt/httprouter"
)

/*type Route interface {
	http.Handler
	Specification() api.Path
}*/

type Router struct {
	httprouter.Router
	OpenAPI api.Api
	Middleware []middleware.MiddlewareHandler
}

func (r *Router) NewRoute(method string, uri string, ep endpoint.EndpointHandler, handlers ...middleware.MiddlewareHandler) {
	r.OpenAPI.NewEndpoint(method, uri, ep.Specification())

	middlewareHandlers := append(handlers, ep.Middleware()...)
	r.Handler(method, uri, middleware.New(ep.(http.Handler), middlewareHandlers...))
}
func (r *Router) Use(h ...middleware.MiddlewareHandler) {
	r.Middleware = append(r.Middleware, h...)
}
func (r *Router) Handle() http.Handler {
	return middleware.New(r, r.Middleware...)
}

func NewRouter(name string, description string, version string) (*Router) {
	r := &Router{
		OpenAPI: api.Api{
			Title: name,
			Description: description,
			Version: version,
		},
	}

	// Ordering matters
	r.Use(
		middleware.WithInitialization(),
		middleware.WithContext(),
		middleware.WithTracing(name),
		middleware.WithMetrics(),
		middleware.WithLogging(),

		//middleware.WithAuthentication(),
	)

	r.NewRoute("GET", "/health", health.NewGetHealthEndpoint())

	r.NewRoute("GET", "/docs", docs.NewGetDocsEndpoint())
	r.NewRoute("GET", "/docs/openapi", docs.NewGetOpenapiEndpoint())
	r.NewRoute("POST", "/docs/openapi", docs.NewGetOpenapiEndpoint())

	r.NewRoute("GET", "/metrics", metrics.NewGetMetricsEndpoint())

	return r
}
