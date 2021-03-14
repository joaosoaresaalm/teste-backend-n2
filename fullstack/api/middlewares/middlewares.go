package middlewares

import (
	"errors"
	"net/http"

	"github.com/joaosoaresa/fullstack/api/auth"
	"github.com/joaosoaresa/fullstack/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidarToken(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Não Authorizado"))
			return
		}
		next(w, r)
	}
}