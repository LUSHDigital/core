package auth

import (
	"testing"
)

func TestErrUnexpectedSigningMethod_Error(t *testing.T) {
	e := ErrUnexpectedSigningMethod{alg: "test"}
	deepEqual(t, "unexpected signing method: test", e.Error())
}

func TestErrAssertClaims_Error(t *testing.T) {
	e := ErrAssertClaims{claims: &Claims{}}
	deepEqual(t, "cannot assert claims for type *auth.Claims", e.Error())
}
