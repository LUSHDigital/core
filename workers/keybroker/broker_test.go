package keybroker_test

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/LUSHDigital/core/workers/keybroker"
	"github.com/dgrijalva/jwt-go"
)

var (
	sourceString = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDx6dqtEuyEf7Mpviqa/rYl316f
OoPozRgG8msH03tC9+exMGUNlExmdMZKgY8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS
2rsFqvITfh0XNQq6W1JB2igTzezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDX
W+kIFfkbaZVWbkUYAwIDAQAB
-----END PUBLIC KEY-----`
	ctx context.Context
)

type badSource struct{}

func (s *badSource) Get(ctx context.Context) ([]byte, error) {
	return []byte{}, nil
}

func Example() {
	broker := keybroker.NewRSA(keybroker.Config{
		Source:   keybroker.JWTPublicKeySources,
		Interval: 5 * time.Second,
	})

	// Run the broker
	go broker.Run(ctx, ioutil.Discard)

	// Queue retrieval of new key
	broker.Renew()

	// Copy the current public key held by the broker
	broker.Copy()
}

func TestServer_Run(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	t.Run("good source", func(t *testing.T) {
		pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(sourceString))
		if err != nil {
			t.Fatal(err)
		}

		b1 := keybroker.NewRSA(keybroker.Config{
			Source:   keybroker.StringSource(sourceString),
			Interval: tick,
		})
		go b1.Run(ctx, ioutil.Discard)
		defer b1.Close()

		time.Sleep(10 * time.Millisecond)
		equals(t, *pk, b1.Copy())
	})

	t.Run("bad source", func(t *testing.T) {
		b2 := keybroker.NewRSA(keybroker.Config{
			Source:   &badSource{},
			Interval: tick,
		})
		go b2.Run(ctx, ioutil.Discard)
		defer b2.Close()

		time.Sleep(10 * time.Millisecond)
		equals(t, *keybroker.DefaultRSA, b2.Copy())
	})
}

func equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
