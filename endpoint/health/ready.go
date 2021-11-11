package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/endpoint/problem"
	hc "github.com/charmixer/golang-api-template/health"

	"go.opentelemetry.io/otel"
)

var (
	OPENAPI_TAGS = []api.Tag{
		{Name: "Health", Description: "Endpoints reporting the health of the application"},
	}
)

var healthChecker *hc.HealthChecker

func init() {
	healthChecker = hc.New(
		hc.WithVersion("0.0.6"),
	)

	healthChecker.AddCheck(
		hc.WithUptimeCheck("host-uptime"),
		hc.WithMemTotalAllocCheck("mem-total-alloc"),
		hc.WithMemObtainedCheck("mem-obtained"),
		hc.WithNumGcCheck("mem-gc-cycles"),
		hc.WithCpuCheck("cpu-usage"),
	)

	ctx := context.Background()

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	isRunning := false
	go func() {
		// Init first run asap
		healthChecker.Check(ctx)
		for {
			select {
			case <-ticker.C:
				if !isRunning {
					isRunning = true
					healthChecker.Check(ctx)
					isRunning = false
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

type GetHealthReadyRequest struct{}
type GetHealthReadyResponse hc.Health

// https://golang.org/doc/effective_go#embedding
type GetHealthReadyEndpoint struct {
	endpoint.Endpoint
}

func (ep GetHealthReadyEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("request")
	ctx, span := tr.Start(ctx, fmt.Sprintf("%s execution", r.URL.Path))
	defer span.End()

	request := GetHealthReadyRequest{}
	if err := endpoint.WithRequestQueryParser(ctx, r, &request); err != nil {
		problem.MustWrite(w, err)
		return
	}

	if err := endpoint.WithRequestValidation(ctx, &request); err != nil {
		problem.MustWrite(w, err)
		return
	}

	if !healthChecker.IsAvailable() {
		prop := problem.New(http.StatusServiceUnavailable).WithDetail("Service is warming up")
		problem.MustWrite(w, prop)
		return
	}

	response := healthChecker.Health()
	if response.Status == hc.Fail {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := endpoint.WithResponseValidation(ctx, response); err != nil {
		problem.MustWrite(w, err)
		return
	}

	if err := endpoint.WithJsonResponseWriter(ctx, w, response); err != nil {
		problem.MustWrite(w, err)
		return
	}

}

func NewGetHealthReadyEndpoint() endpoint.EndpointHandler {
	ep := GetHealthReadyEndpoint{}

	ep.Setup(
		endpoint.WithSpecification(api.Path{
			Summary:     "Get health information about the service",
			Description: ``,
			Tags:        OPENAPI_TAGS,

			Request: api.Request{
				Description: ``,
				Schema:      GetHealthReadyRequest{},
			},

			Responses: []api.Response{{
				Description: http.StatusText(http.StatusOK),
				Code:        http.StatusOK,
				Schema:      GetHealthReadyResponse{},
			}, {
				Description: http.StatusText(http.StatusBadRequest),
				Code:        http.StatusBadRequest,
				Schema:      problem.ValidationProblem{},
			}, {
				Description: http.StatusText(http.StatusServiceUnavailable),
				Code:        http.StatusServiceUnavailable,
				Schema:      problem.ProblemDetails{},
			}},
		}),
	)

	return ep
}
