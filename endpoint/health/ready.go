package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmixer/oas/api"

	"github.com/charmixer/golang-api-template/app"
	"github.com/charmixer/golang-api-template/endpoint"
	"github.com/charmixer/golang-api-template/endpoint/problem"
	hc "github.com/charmixer/golang-api-template/health"

	"go.opentelemetry.io/otel"
)

var healthChecker *hc.HealthChecker

func init() {
	healthChecker = hc.New(
		hc.WithDescription("Readiness of the service"),
	)

	healthChecker.AddCheck(
		hc.WithUptimeCheck("host-uptime"),
		hc.WithMemTotalAllocCheck("mem-total-alloc"),
		hc.WithMemObtainedCheck("mem-obtained"),
		hc.WithNumGcCheck("mem-gc-cycles"),
		hc.WithCpuCheck("cpu-usage"),

		hc.WithBuildNameCheck("build-name"),
		hc.WithBuildTagCheck("build-tag"),
		hc.WithBuildCommitCheck("build-commit"),
		hc.WithBuildEnvironmentCheck("build-environment"),
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
	ctx, span := otel.Tracer("request").Start(r.Context(), fmt.Sprintf("%s handler", r.URL.Path))
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

	// Has to be set after init, since env is not available until then
	healthChecker.SetOption(
		hc.WithVersion(app.Env.Build.Version),
		hc.WithReleaseId(app.Env.Build.Commit),
	)

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
