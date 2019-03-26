package keys_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
	"time"

	"github.com/LUSHDigital/core/keys"
	"github.com/LUSHDigital/core/keys/keysmock"
	"github.com/dgrijalva/jwt-go"
)

var (
	sourceString = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDx6dqtEuyEf7Mpviqa/rYl316f
OoPozRgG8msH03tC9+exMGUNlExmdMZKgY8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS
2rsFqvITfh0XNQq6W1JB2igTzezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDX
W+kIFfkbaZVWbkUYAwIDAQAB
-----END PUBLIC KEY-----`
)

type badSource struct{}

func (s *badSource) Get(ctx context.Context) ([]byte, error) {
	return []byte{}, nil
}

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
	deepEqual(t, public, mock.Copy())
}

func Test_BrokerRSAPublicKey(t *testing.T) {
	ctx := context.Background()
	source := keys.StringSource(sourceString)
	tick := 5 * time.Millisecond

	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(sourceString))
	if err != nil {
		t.Fatal(err)
	}

	b1 := keys.BrokerRSAPublicKey(ctx, source, tick)
	defer b1.Close()

	time.Sleep(10 * time.Millisecond)
	deepEqual(t, *pk, b1.Copy())

	b2 := keys.BrokerRSAPublicKey(ctx, &badSource{}, tick)
	defer b2.Close()

	time.Sleep(10 * time.Millisecond)
	deepEqual(t, *keys.DefaultRSA, b2.Copy())
}

func deepEqual(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
