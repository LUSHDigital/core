package auth

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-core-golang/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
)

const (
	authHeader       = "Authorization"
	authHeaderPrefix = "Bearer "

	errorMissingRequiredGrants = "missing required grants"
)

type jwtHTTPResponder struct {
	w    http.ResponseWriter
	r    *http.Request
	next http.HandlerFunc
}

func (r *jwtHTTPResponder) OnUnauthorizedErr(err error) {
	response.New(http.StatusUnauthorized, err.Error(), nil).WriteTo(r.w)
}

func (r *jwtHTTPResponder) OnComplete(token *jwt.Token) {
	ctx := ContextWithConsumer(r.r.Context(), token.Claims.(*Claims).Consumer)
	r.next.ServeHTTP(r.w, r.r.WithContext(ctx))
}

// HandlerValidateJWT takes a JWT from the request headers, attempts validation and returns a http handler.
func HandlerValidateJWT(pk *rsa.PublicKey, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)
		responder := &jwtHTTPResponder{w, r, next}
		RespondToJWT(pk, token, responder)
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

		response.New(http.StatusUnauthorized, errorMissingRequiredGrants, nil).WriteTo(w)
	})
}
