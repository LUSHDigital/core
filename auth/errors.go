package auth

import (
	"fmt"
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
