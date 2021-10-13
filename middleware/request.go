package middleware

import (
	"net/http"

	"encoding/json"

	"go.opentelemetry.io/otel"
)

// TODO move this to better place
func WithJsonRequestParser(request interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tr := otel.Tracer("request")
			ctx, span := tr.Start(ctx, "middleware.request-parser")
			defer span.End()

			// Try to decode the request body into the struct. If there is an error,
			// respond to the client with the error message and a 400 status code.
			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
