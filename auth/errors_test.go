package auth

import (
	"reflect"
	"testing"
)

func TestUnexpectedSigningMethodError_Error(t *testing.T) {
	e := UnexpectedSigningMethodError{alg: "test"}
	deepEqual(t, "unexpected signing method: test", e.Error())
}

func TestErrAssertClaims_Error(t *testing.T) {
	e := ErrAssertClaims{claims: &Claims{}}
	deepEqual(t, "cannot assert claims for type *auth.Claims", e.Error())
}

func deepEqual(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
