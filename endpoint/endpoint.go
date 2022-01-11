package endpoint

import (
	"context"
	"net/http"

	"github.com/charmixer/golang-api-template/endpoint/problem"
	"github.com/charmixer/golang-api-template/middleware"
	"github.com/charmixer/oas/api"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type EndpointHandler interface {
	http.Handler
	Specification() api.Path
	HandleInternalError(context.Context, trace.Span, error) *problem.ProblemDetails
	Middleware() []middleware.MiddlewareHandler
}
type Endpoint struct {
	EndpointHandler
	specification api.Path
	middleware    []middleware.MiddlewareHandler
}
type EndpointOption func(e *Endpoint)

func (ep *Endpoint) Setup(options ...EndpointOption) {
	for _, opt := range options {
		opt(ep)
	}
}

func (ep Endpoint) HandleInternalError(ctx context.Context, span trace.Span, err error) *problem.ProblemDetails {
	log.Error().Err(err)
	span.SetStatus(http.StatusInternalServerError, err.Error())
	return problem.New(http.StatusInternalServerError).WithErr(err)
}

func (ep Endpoint) Specification() api.Path {
	return ep.specification
}
func (ep Endpoint) Middleware() []middleware.MiddlewareHandler {
	return ep.middleware
}

func WithSpecification(spec api.Path) EndpointOption {
	return func(e *Endpoint) {
		e.specification = spec
	}
}

func WithMiddleware(handlers ...middleware.MiddlewareHandler) EndpointOption {
	return func(e *Endpoint) {
		e.middleware = append(e.middleware, handlers...)
	}
}

// https://pkg.go.dev/github.com/alexliesenfeld/health#WithCacheDuration
// https://github.com/alexliesenfeld/health/blob/v0.6.0/config.go#L159
