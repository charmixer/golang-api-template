package middleware

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	Status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, Status: http.StatusOK}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		wrapped := wrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		log.Info().
			Str("type", "access").
			Str("request_id", r.Context().Value("req_id").(string)).
			Str("remote_ip", r.Context().Value("remote_ip").(string)).
			Str("user_agent", r.Context().Value("user_agent").(string)).
			Str("referer", r.Context().Value("referer").(string)).
			Str("method", r.Method).
			Str("duration", time.Since(start).String()).
			Int("status", wrapped.Status).
			Stringer("url", r.URL).
			Msg("")
	})
}

func Logging2(c alice.Chain) alice.Chain {
	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(log.Logger))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to that handler, all our logs will come with some prepopulated fields.
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("type", "access").
			Str("request_id", r.Context().Value("req_id").(string)).
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))

	return c
}
