package middleware

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
	"encoding/json"
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

func WithRequestParser(request interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("With Request Parser")

			next.ServeHTTP(w, r)
		})
	}
}

func WithRequestValidation(request interface{}) MiddlewareHandler{
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("With Request Validation")

			next.ServeHTTP(w, r)
		})
	}
}

func WithResponseWriter(response interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			contentType := w.Header().Get("Content-type")
			switch (contentType) {
			case "application/json":
				json.NewEncoder(w).Encode(response)
				break;
			case "application/x-yaml":
				yaml.NewEncoder(w).Encode(response)
				break;
			default:
				panic("Missing Content-Type header")
			}

		})
	}
}

func WithResponseValidation(response interface{}) MiddlewareHandler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "error", http.StatusBadRequest)

			next.ServeHTTP(w, r)
		})
	}
}
