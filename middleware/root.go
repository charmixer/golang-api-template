package middleware

import (
	"net/http"

	"github.com/justinas/alice"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured from middleware, like logging.
type responseWriter struct {
	http.ResponseWriter
	Status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, Status: http.StatusOK}
}

func ResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
	})
}

func GetChain() alice.Chain {
	return alice.New(ResponseWriter, Context, Metrics, Logging)
}
