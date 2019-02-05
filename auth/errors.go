package auth

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenMalformed is the message to return for a malformed token.
	ErrTokenMalformed = errors.New("token malformed")

	// ErrTokenExpired is the message to return for an expired token.
	ErrTokenExpired = errors.New("token expired or not yet valid")

	// ErrTokenInvalid is the message to return for an invalid token.
	ErrTokenInvalid = errors.New("invalid token")
)

// UnexpectedSigningMethodError is thrown when parsing a JWT encounters an
// unexpected signature method.
type UnexpectedSigningMethodError struct {
	alg interface{}
}

func (e UnexpectedSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", e.alg)
}

// ErrAssertClaims is thrown when asserting the type of claims
type ErrAssertClaims struct {
	claims interface{}
}

func (e *ErrAssertClaims) Error() string {
	return fmt.Sprintf("cannot assert claims for type %T", e.claims)
}

func handleParseErr(err error) error {
	if _, ok := err.(*jwt.ValidationError); ok {
		// Handle any token specific errors.
		if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorMalformed != 0 {
			err = ErrTokenMalformed
		} else if err.(*jwt.ValidationError).Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			err = ErrTokenExpired
		} else {
			err = ErrTokenInvalid
		}
	}
	return err
}
