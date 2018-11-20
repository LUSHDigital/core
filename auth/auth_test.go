package auth_test

import (
	"log"
	"os"
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/auth"
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
)

func TestMain(m *testing.M) {
	if tokeniser, err = auth.NewTokeniser(testPrivateKey, testPublicKey, "testing"); err != nil {
		log.Fatalf("could not parse JWT keys: %s", err)
	}

	os.Exit(m.Run())
}
