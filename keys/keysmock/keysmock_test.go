package keysmock_test

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"

	"github.com/LUSHDigital/core/keys/keysmock"
)

func Test_MockRSAPublicKey(t *testing.T) {
	private, err := rsa.GenerateKey(rand.Reader, 128)
	if err != nil {
		t.Fatal(err)
	}
	public := private.PublicKey
	mock := keysmock.MockRSAPublicKey(public)
	if err != nil {
		t.Fatal(err)
	}
	equals(t, public, mock.Copy())
}

func equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
