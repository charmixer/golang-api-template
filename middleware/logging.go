package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/rs/zerolog/log"
)

func WithLogging() MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tr := otel.Tracer("request")
			ctx, span := tr.Start(ctx, "middleware.logging")
			defer span.End()

			start := time.Now()

			wrapped := w.(*responseWriter)

			next.ServeHTTP(wrapped, r.WithContext(ctx))

			ctx, span = tr.Start(r.Context(), "write request to log")
			defer span.End()

			log.Info().
				Str("type", "access").
				Str("request_id", ctx.Value("req_id").(string)).
				Str("remote_ip", ctx.Value("remote_ip").(string)).
				Str("user_agent", ctx.Value("user_agent").(string)).
				Str("referer", ctx.Value("referer").(string)).
				Str("method", r.Method).
				Str("duration", time.Since(start).String()).
				Int("status", wrapped.Status).
				Stringer("url", r.URL).
				Msg("")
		})
	}
}
