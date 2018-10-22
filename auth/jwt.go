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

// JWTClaims represents the claims within the JWT.
type JWTClaims struct {
	Consumer Consumer `json:"consumer"`
	jwt.StandardClaims
}

// HandlerValidateJWT takes a JWT from the request headers, attempts validation and returns a http handler.
func HandlerValidateJWT(publicKey *rsa.PublicKey, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)

		// Parse the JWT token.
		authToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(aToken *jwt.Token) (interface{}, error) {
			// Ensure the signing method was not changed.
			if _, ok := aToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", aToken.Header["alg"])
			}

			return publicKey, nil
		})

		// Bail out if the token could not be parsed.
		if err != nil {
			if _, ok := err.(*jwt.ValidationError); ok {
				// Handle any token specific errors.
				var errorMessage string

				if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorMalformed != 0 {
					errorMessage = errorMessageMalformed
				} else if err.(*jwt.ValidationError).Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					errorMessage = errorMessageExpired
				} else {
					errorMessage = errorMessageInvalid
				}

				response.New(http.StatusUnauthorized, errorMessage, nil).WriteTo(w)
				return
			}

			response.New(http.StatusUnauthorized, errorMessageInvalid, nil).WriteTo(w)
		}

		// Check the claims and token are valid.
		if _, ok := authToken.Claims.(*JWTClaims); !ok || !authToken.Valid {
			response.New(http.StatusUnauthorized, errorMessageClaimsInvalid, nil).WriteTo(w)
			return
		}

		// Add the customer ID to the request context.
		ctx := ContextWithConsumer(r.Context(), authToken.Claims.(*JWTClaims).Consumer)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
