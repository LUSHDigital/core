package keybrokermock_test

import (
	"testing"

	"github.com/LUSHDigital/core/auth/authmock"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/keybroker/keybrokermock"
)

func Test_MockRSAPublicKey(t *testing.T) {
	_, public := authmock.MustNewRSAKeyPair()
	mock := keybrokermock.MockRSAPublicKey(public)
	test.Equals(t, *public, mock.Copy())
}

func Test_MockECDSAPublicKey(t *testing.T) {
	_, public := authmock.MustNewECDSAKeyPair()
	mock := keybrokermock.MockECDSAPublicKey(public)
	test.Equals(t, *public, mock.Copy())
}
