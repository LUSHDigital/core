package keybroker_test

import (
	"context"
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
	privSourceString = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDx6dqtEuyEf7Mpviqa/rYl316fOoPozRgG8msH03tC9+exMGUN
lExmdMZKgY8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS2rsFqvITfh0XNQq6W1JB2igT
zezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDXW+kIFfkbaZVWbkUYAwIDAQAB
AoGAQZnU/xIeqV+nyi4Th6yC4IpOMoe/taXIWjnq5FhpGKP5ZIdnH+OTREVucE3p
6JBxyC4TG6EHh0KfX0dU5xHGp5ncts8QOhzZ3uJNwKsG6OAaNXI9pkhxty8EHhC9
GPP+fZdAmEtQhzpN2wfMxO1Q6vub6c7HmAkFh7cYFHlwWcECQQD/z8LOR2G6G7PR
AWBcyML4nWPPFagf9Rl37hoHd75Vy9wXKQOW6b9lkg2XjETj7dR+/Aha0xy28f/x
A/v77ECJAkEA8hd48l1Ec3WT/dKrIw2I4xgfQtwi4H/qH0hKtWsqWFnU8T2TssvO
yKMx3uExS4yn3eWJiO4t+Dah1C88Hgn5KwJATv3LrMDUB5D4VKi1JdGEixqmsFKC
qOOZarQma3npVzrtCmXKyvYA+Q9BjTNuLmfJPzD6L3mTG1bc7oKJdAA+6QJAFzdz
DOMu5le3SpdCfEkXAJhWnyXXXmpF/JwFNiLB29k5l60NFg9/lDQ6WyKDhLhHfPs/
VldpJy2uFVg2TrcsIwJAJw/vz25NbLxibdJ6cqZKF30411tgufDIjgzVr9MQX3np
2elP0lxdJ9FzNP+q6BV4J48/yrDJZLtSGJkFExr2fA==
-----END RSA PRIVATE KEY-----`
	ctx context.Context
)

type badSource struct{}

func (s *badSource) Get(ctx context.Context) ([]byte, error) {
	return []byte{}, nil
}

func Example() {
	broker := keybroker.NewPublicRSA(&keybroker.Config{
		Source:   keybroker.JWTPublicKeySources,
		Interval: 5 * time.Second,
	})

	// Run the broker
	go broker.Run(ctx)

	// Queue retrieval of new key
	broker.Renew()

	// Copy the current public key held by the broker
	broker.Copy()
}

func TestRSAPublicKeyBroker_Run(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	t.Run("good source", func(t *testing.T) {
		pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(sourceString))
		if err != nil {
			t.Fatal(err)
		}

		b1 := keybroker.NewPublicRSA(&keybroker.Config{
			Source:   keybroker.StringSource(sourceString),
			Interval: tick,
		})
		go b1.Run(ctx)
		defer b1.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *pk, b1.Copy())
	})

	t.Run("bad source", func(t *testing.T) {
		b2 := keybroker.NewPublicRSA(&keybroker.Config{
			Source:   &badSource{},
			Interval: tick,
		})
		go b2.Run(ctx)
		defer b2.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *keybroker.DefaultRSA, b2.Copy())
	})
}

func TestRSAPublicKeyBroker_Check(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	b1 := keybroker.NewPublicRSA(&keybroker.Config{
		Source:   keybroker.StringSource(sourceString),
		Interval: tick,
	})

	go b1.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, true, ok)
		test.Equals(t, "rsa public key broker has retrieved key of size 128", messages[0])
	})

	b1.Close()
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, stopped", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "rsa public key broker is not yet running", messages[0])
	})

	b2 := keybroker.NewPublicRSA(&keybroker.Config{
		Source:   &badSource{},
		Interval: tick,
	})
	go b2.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b2.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "rsa public key broker has not yet retrieved a key", messages[0])
	})
}

func TestRSAPrivateKeyBroker_Run(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	t.Run("good source", func(t *testing.T) {
		pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privSourceString))
		if err != nil {
			t.Fatal(err)
		}
		b1 := keybroker.NewPrivateRSA(&keybroker.Config{
			Source:   keybroker.StringSource(privSourceString),
			Interval: tick,
		})
		go b1.Run(ctx)
		defer b1.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *pk, b1.Copy())
	})

	t.Run("bad source", func(t *testing.T) {
		b2 := keybroker.NewPrivateRSA(&keybroker.Config{
			Source:   &badSource{},
			Interval: tick,
		})
		go b2.Run(ctx)
		defer b2.Close()

		time.Sleep(10 * time.Millisecond)
		test.Equals(t, *keybroker.DefaultPrivateRSA, b2.Copy())
	})
}

func TestRSAPrivateKeyBroker_Check(t *testing.T) {
	ctx = context.Background()
	tick := 5 * time.Millisecond

	b1 := keybroker.NewPrivateRSA(&keybroker.Config{
		Source:   keybroker.StringSource(privSourceString),
		Interval: tick,
	})

	go b1.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, true, ok)
		test.Equals(t, "rsa private key broker has retrieved key of size 128", messages[0])
	})

	b1.Close()
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, stopped", func(t *testing.T) {
		messages, ok := b1.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "rsa private key broker is not yet running", messages[0])
	})

	b2 := keybroker.NewPrivateRSA(&keybroker.Config{
		Source:   &badSource{},
		Interval: tick,
	})
	go b2.Run(ctx)
	time.Sleep(10 * time.Millisecond)

	t.Run("good source, started", func(t *testing.T) {
		messages, ok := b2.Check()
		test.Equals(t, false, ok)
		test.Equals(t, "rsa private key broker has not yet retrieved a key", messages[0])
	})
}
