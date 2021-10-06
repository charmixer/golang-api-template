package middleware

import (
	"fmt"
	"net/http"
  "encoding/json"

	"gopkg.in/yaml.v2"

  "go.opentelemetry.io/otel"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func WithJsonResponseWriter(response interface{}, status ...int) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context()
    	tr := otel.Tracer("request")
    	ctx, span := tr.Start(ctx, "middleware.json-response-writer")
    	defer span.End()

    	next.ServeHTTP(w, r)

			// Only write if written statuscode is in given list
			if contains(status, w.(*responseWriter).Status) {
				ctx, span = tr.Start(ctx, "write response")
				defer span.End()

				d, err := json.Marshal(response)
				if err != nil {
					panic(err) // TODO FIXME
				}
				w.Write(d)
			}
		})
	}
}

func WithYamlResponseWriter(response interface{}, status ...int) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context()
    	tr := otel.Tracer("request")
    	ctx, span := tr.Start(ctx, "middleware.yaml-response-writer")
    	defer span.End()

    	next.ServeHTTP(w, r)

			// Only write if written statuscode is in given list
			if contains(status, w.(*responseWriter).Status) {
				ctx, span = tr.Start(ctx, "write response")
				defer span.End()

				d, err := yaml.Marshal(response)
				if err != nil {
					panic(err) // TODO FIXME
				}
				w.Write(d)
      }
		})
	}
}

func WithResponseWriter(tp *string, response interface{}, status ...int) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      ctx := r.Context()
    	tr := otel.Tracer("request")
    	ctx, span := tr.Start(ctx, "middleware.response-writer")
    	defer span.End()

    	next.ServeHTTP(w, r)

			// Only write if written statuscode is in given list
			if contains(status, w.(*responseWriter).Status) {
				ctx, span = tr.Start(ctx, "write response")
				defer span.End()

				switch (*tp) {
				case "json":
					d, err := json.Marshal(response)
					if err != nil {
						panic(err) // TODO FIXME
					}
					w.Write(d)
					break;
				case "yaml":
					d, err := yaml.Marshal(response)
					if err != nil {
						panic(err) // TODO FIXME
					}
					w.Write(d)
					break;
				default:
					panic(fmt.Sprintf("Unknown response type given, %s", *tp))
				}
			}

		})
	}
}
