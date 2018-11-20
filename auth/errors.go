package auth

import (
	"errors"
	"fmt"
)

var (
	// ErrTokenMalformed is the message to return for a malformed token.
	ErrTokenMalformed = errors.New("token malformed")

	// ErrTokenExpired is the message to return for an expired token.
	ErrTokenExpired = errors.New("token expired or not yet valid")

	// ErrTokenInvalid is the message to return for an invalid token.
	ErrTokenInvalid = errors.New("invalid token")
)

// ErrUnexpectedSigningMethod is thrown when parsing a JWT encounters an
// unexpected signature method.
type ErrUnexpectedSigningMethod struct {
	alg interface{}
}

func (e *ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", e.alg)
}

// ErrAssertClaims is thrown when asserting the type of claims
type ErrAssertClaims struct {
	claims interface{}
}

func (e *ErrAssertClaims) Error() string {
	return fmt.Sprintf("cannot assert claims for type %T", e.claims)
}
