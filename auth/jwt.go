package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-core-golang/response"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	authHeader       = "Authorization"
	authHeaderPrefix = "Bearer "

	errorMessageMalformed     = "token malformed"
	errorMessageExpired       = "token expired or not yet valid"
	errorMessageInvalid       = "invalid token"
	errorMessageClaimsInvalid = "invalid token claims"
)

// ConsumerFor derives the Consumer from the JWT claims
func ConsumerFor(token *jwt.Token) *Consumer {
	return &token.Claims.(*JWTClaims).Consumer
}

// JWTClaims represents the claims within the JWT.
type JWTClaims struct {
	Consumer Consumer `json:"consumer"`
	jwt.StandardClaims
}

// JWTResponder defines the behaviour of validating a JWT
type JWTResponder interface {
	OnUnauthorizedErr(err error)
	OnComplete(token *jwt.Token)
}

type jwtHTTPResponder struct {
	w    http.ResponseWriter
	r    *http.Request
	next http.HandlerFunc
}

func (r *jwtHTTPResponder) OnUnauthorizedErr(err error) {
	response.New(http.StatusUnauthorized, err.Error(), nil).WriteTo(r.w)
}

func (r *jwtHTTPResponder) OnComplete(token *jwt.Token) {
	ctx := ContextWithConsumer(r.r.Context(), token.Claims.(*JWTClaims).Consumer)
	r.next.ServeHTTP(r.w, r.r.WithContext(ctx))
}

// ValidateJWT takes the raw JWT and the public RSA key
func ValidateJWT(raw string, publicKey *rsa.PublicKey, responder JWTResponder) {
	// Parse the JWT token
	token, err := jwt.ParseWithClaims(raw, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method was not changed
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	// Bail out if the token could not be parsed
	if err != nil {
		if _, ok := err.(*jwt.ValidationError); ok {
			// Handle any token specific errors
			var errorMessage string
			if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorMalformed != 0 {
				errorMessage = errorMessageMalformed
			} else if err.(*jwt.ValidationError).Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				errorMessage = errorMessageExpired
			} else {
				errorMessage = errorMessageInvalid
			}
			responder.OnUnauthorizedErr(fmt.Errorf(errorMessage))
			return
		}
		responder.OnUnauthorizedErr(fmt.Errorf(errorMessageInvalid))
		return
	}

	// Check the claims and token are valid
	if _, ok := token.Claims.(*JWTClaims); !ok || !token.Valid {
		responder.OnUnauthorizedErr(fmt.Errorf(errorMessageClaimsInvalid))
		return
	}

	responder.OnComplete(token)
}

// HandlerValidateJWT takes a JWT from the request headers, attempts validation and returns a http handler.
func HandlerValidateJWT(publicKey *rsa.PublicKey, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)
		responder := &jwtHTTPResponder{w, r, next}
		ValidateJWT(token, publicKey, responder)
	})
}
