package auth

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWTResponder defines the behaviour of validating a JWT
type JWTResponder interface {
	OnUnauthorizedErr(err error)
	OnComplete(token *jwt.Token)
}

// RespondToJWT takes the raw JWT and the public RSA key
func RespondToJWT(pk *rsa.PublicKey, raw string, responder JWTResponder) {
	token, err := ParseJWT(pk, raw)
	if err != nil {
		responder.OnUnauthorizedErr(err)
		return
	}
	responder.OnComplete(token)
}
