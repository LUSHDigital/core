package auth

import (
	"testing"

	"gitlab.platformserviceaccount.com/lush-soa/service/markets/service/test"
)

func TestErrUnexpectedSigningMethod_Error(t *testing.T) {
	e := ErrUnexpectedSigningMethod{alg: "test"}
	test.Equals(t, "unexpected signing method: test", e.Error())
}

func TestErrAssertClaims_Error(t *testing.T) {
	e := ErrAssertClaims{claims: &Claims{}}
	test.Equals(t, "cannot assert claims for type *auth.Claims", e.Error())
}
