package middleware

import (
	"net/http"

	// "github.com/gorilla/mux"

//"github.com/julienschmidt/httprouter"

	//"github.com/rs/zerolog/log"
)

func WithInputValidation() (MiddlewareHandler) {
	return func(next http.Handler) http.Handler {
	  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	  	next.ServeHTTP(w, r)
	  })
	}
}

func WithOutputValidation() (MiddlewareHandler) {
  return func (next http.Handler) http.Handler {
  	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  		next.ServeHTTP(w, r)
  	})
  }
}
