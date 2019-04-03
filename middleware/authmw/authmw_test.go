package authmw_test

import (
	"crypto/rsa"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/middleware/authmw"
	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
)

var (
	testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
	testPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDx6dqtEuyEf7Mpviqa/rYl316f
OoPozRgG8msH03tC9+exMGUNlExmdMZKgY8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS
2rsFqvITfh0XNQq6W1JB2igTzezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDX
W+kIFfkbaZVWbkUYAwIDAQAB
-----END PUBLIC KEY-----`
	testIncorrectPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDx6dqtEuyEf7Mpviqa/rYl316f
OoPozRgG8msH03tC9+esdfseftddfsefes8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS
2rsFqvITfh0XNQq6W1JB2igTzezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDX
W+kIFfkbaZVWbkUYAwIDAQAB
-----END PUBLIC KEY-----`

	issuer      *auth.Issuer
	correctPK   *rsa.PublicKey
	incorrectPK *rsa.PublicKey
	broker      auth.RSAPublicKeyCopierRenewer
)

func TestMain(m *testing.M) {
	var err error
	issuer, err = auth.NewIssuerFromPrivateKeyPEM(auth.IssuerConfig{
		Name: "test",
	}, []byte(testPrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	correctPK, err = jwt.ParseRSAPublicKeyFromPEM([]byte(testPublicKey))
	if err != nil {
		log.Fatalln(err)
	}
	incorrectPK, err = jwt.ParseRSAPublicKeyFromPEM([]byte(testIncorrectPublicKey))
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func ExampleStreamServerInterceptor() {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(authmw.StreamServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleUnaryServerInterceptor() {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authmw.UnaryServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleHandlerValidateJWT() {
	http.Handle("/users", authmw.HandlerValidateJWT(broker, func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant("users.read") {
			http.Error(w, "access denied", http.StatusUnauthorized)
		}
	}))
}

func equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
