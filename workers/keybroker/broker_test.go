package keybroker_test

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/LUSHDigital/core/test"
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
	broker := keybroker.NewRSA(&keybroker.Config{
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

		b1 := keybroker.NewRSA(&keybroker.Config{
			Source:   keybroker.StringSource(sourceString),
			Interval: tick,
		})
		go b1.Run(ctx, ioutil.Discard)
		defer b1.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *pk, b1.Copy())
	})

	t.Run("bad source", func(t *testing.T) {
		b2 := keybroker.NewRSA(&keybroker.Config{
			Source:   &badSource{},
			Interval: tick,
		})
		go b2.Run(ctx, ioutil.Discard)
		defer b2.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *keybroker.DefaultRSA, b2.Copy())
	})
}

func TestServer_Check(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	b1 := keybroker.NewRSA(&keybroker.Config{
		Source:   keybroker.StringSource(sourceString),
		Interval: tick,
	})

	go b1.Run(ctx, ioutil.Discard)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, true, ok)
		test.Equals(t, "broker has retrieved key of size 128", messages[0])
	})

	b1.Close()
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, stopped", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "broker is not yet running", messages[0])
	})

	b2 := keybroker.NewRSA(&keybroker.Config{
		Source:   &badSource{},
		Interval: tick,
	})
	go b2.Run(ctx, ioutil.Discard)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b2.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "broker has not yet retrieved a key", messages[0])
	})
}
