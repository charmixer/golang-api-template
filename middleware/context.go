package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/gofrs/uuid"
)

func WithContext() MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for incoming header, use it if exists
			reqID := r.Header.Get("X-Request-Id")

			// Create request id with UUID4
			if reqID == "" {
				uuid4, _ := uuid.NewV4()
				reqID = uuid4.String()
			}

			r.Header.Set("X-Request-Id", reqID)

			ctx := r.Context()

			ctx = context.WithValue(ctx, "req_id", reqID)

			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			ctx = context.WithValue(ctx, "remote_ip", host)

			ua := r.Header.Get("User-Agent")
			ctx = context.WithValue(ctx, "user_agent", ua)

			ref := r.Header.Get("Referer")
			ctx = context.WithValue(ctx, "referer", ref)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
