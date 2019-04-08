package keybrokermock_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/keybroker/keybrokermock"
)

func Test_MockRSAPublicKey(t *testing.T) {
	private, err := rsa.GenerateKey(rand.Reader, 128)
	if err != nil {
		t.Fatal(err)
	}
	public := private.PublicKey
	mock := keybrokermock.MockRSAPublicKey(public)
	if err != nil {
		t.Fatal(err)
	}
	test.Equals(t, public, mock.Copy())
}
