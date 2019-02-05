package auth

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
)

// ConsumerFor derives the Consumer from the JWT claims
// DEPRECATED: This should no longer be used in favour of creating a tokeniser
func ConsumerFor(token *jwt.Token) (*Consumer, error) {
	if claims, ok := token.Claims.(*Claims); ok {
		return &claims.Consumer, nil
	}
	return nil, &ErrAssertClaims{token.Claims}
}

// ParseJWT parses a JWT string and checks its signature validity
// DEPRECATED: This should no longer be used in favour of creating a tokeniser
func ParseJWT(pk *rsa.PublicKey, raw string) (*jwt.Token, error) {
	tokeniser := &Tokeniser{publicKey: pk}
	return tokeniser.ParseToken(raw)
}

func checkSignatureFunc(pk *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method was not changed
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, UnexpectedSigningMethodError{token.Header["alg"]}
		}
		return pk, nil
	}
}
