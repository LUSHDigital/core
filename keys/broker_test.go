package keys_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/LUSHDigital/microservice-core-golang/keys"
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

func Test_BrokerRSAPublicKey(t *testing.T) {
	ctx := context.Background()
	source := keys.StringSource(sourceString)
	tick := 5 * time.Millisecond

	pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(sourceString))
	if err != nil {
		t.Fatal(err)
	}

	broker, cancel := keys.BrokerRSAPublicKey(ctx, source, tick)
	defer cancel()

	time.Sleep(10 * time.Millisecond)
	deepEqual(t, *pk, broker.Copy())
}

func deepEqual(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
