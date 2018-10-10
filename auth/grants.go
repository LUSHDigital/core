package auth

import (
	"github.com/LUSHDigital/microservice-core-golang/response"
	"net/http"
)

const errorMissingRequiredGrants = "missing required grants"

// Grant is a type alised string to represent permissions grants.
type Grant string

// HandlerGrants is an HTTP handler to check that the consumer in the request context has the required grants.
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
