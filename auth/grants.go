package auth

import (
	"github.com/LUSHDigital/microservice-core-golang/response"
	"net/http"
)

const errorMissingRequiredGrants = "missing required grants"

type Grant string

func HandlerGrants(grants []Grant, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := ConsumerFromContext(r.Context())
		if consumer.HasGrants(grants...) {
			next.ServeHTTP(w, r)
			return
		}

		response.New(http.StatusUnauthorized, errorMissingRequiredGrants, nil).WriteTo(w)
	})
}
