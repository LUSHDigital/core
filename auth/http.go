package auth

import (
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-core-golang/keys"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/gofrs/uuid"
)

const (
	authHeader               = "Authorization"
	authHeaderPrefix         = "Bearer "
	msgMissingRequiredGrants = "missing required grants"
)

// HandlerValidateJWT takes a JWT from the request headers, attempts validation and returns a http handler.
func HandlerValidateJWT(brk keys.RSAPublicKeyCopierRenewer, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)
		pk := brk.Copy()
		tokeniser := &Tokeniser{publicKey: &pk}
		token, err := tokeniser.ParseToken(raw)
		if err != nil {
			switch err.(type) {
			case TokenSignatureError:
				brk.Renew() // Renew the public key if there's an error validating the token signature
			}
			response.New(http.StatusUnauthorized, err.Error(), nil).WriteTo(w)
			return
		}
		ctx := ContextWithConsumer(r.Context(), token.Claims.(*Claims).Consumer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// EnsureRequestID will create a Request ID header if one is not found.
// It will then place the request ID into the request's context.
func EnsureRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-request-id") == "" {
			requestID, err := uuid.NewV4()
			if err != nil {
				response.InternalError(err).WriteTo(w)
				return
			}
			r.Header.Add("x-request-id", requestID.String())
		}
		ctxWithReqID := NewContextWithRequestID(r.Context(), r)
		next.ServeHTTP(w, r.WithContext(ctxWithReqID))
	})
}

// HandlerGrants is an HTTP handler to check that the consumer in the request context has the required grants.
func HandlerGrants(grants []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := ConsumerFromContext(r.Context())
		if consumer.HasAnyGrant(grants...) {
			next.ServeHTTP(w, r)
			return
		}
		response.New(http.StatusUnauthorized, msgMissingRequiredGrants, nil).WriteTo(w)
	})

}
