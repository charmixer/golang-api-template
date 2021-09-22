package middleware

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/zerolog/log"
)

func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var match mux.RouteMatch
		routeExists := s.Router.Match(r, &match)
		if routeExists && match.Route.GetName(){
				routeName := match.Route.GetName()
		}
		
		log.Info().Msgf("======== %#v", route)

		next.ServeHTTP(w, r)

	})
}

func ValidateResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

	})
}

