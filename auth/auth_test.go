package auth_test

import (
	"crypto/rsa"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/core/auth"
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
	testIncorrectPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC2BcJ5OatQY653s21PKAt6ttNxfC0iurY1foZL41540viPJ75U
tuMd0o/II9NJnyBOHyn6U1tu3f6saE+eqHJqh3iO35qWLMaFUwIoiJU93mtx966O
WOENchowmZR3QRVxLQkFtAhYRJdsJ+ICwNYJsXOIDutwVMoU0E/a+flLrQIDAQAB
AoGABQCJhJ1SGOZ0X/O9WESIdDnb+61m7CJnaXbtp946tWVO0VhNQbS47xPfQafC
Ya6Oy7bNh4SM6bIOEpzXO0vzwO/ULPYqV59HQIfe9zdMi3f4aocDz1WY+0lrX0Q/
ZDK8kM1CJQDkXUhS/FcJaPx16KVpJOl3gV3t/FljHkpiwkECQQDcZqEBICqq7jlv
17BaPBjeu37E0n1rlTdNiFECnnFTW0WzD/MW6Xsz+h6MglB8t+hPFh1qvPG0ZG0F
7SbAM/nxAkEA02w4Tbuxo0Q11nHQIU4XIz1P0Jc4Gb3+dFxUfEll22ygDxa+UU5Q
h4zcvfM838eMzJNimC6pEMhDYmmSm7NRfQJBAI9KagLBVvwqRU1hfVYtHD4yyAhO
kRwQtxPBPGnneOYowPfZtsF+qorwYkwXrRxotLA2QInUrZAKepcPx9HN+QECQQDL
r9BSu4h5dga0UjQlUhmifrg9iuKmkk/qhOV0VDZIfs95mfzDUkLtRL2KVyQHqDWz
Bi+P1CxXmcipsHJphQn1AkBScr2xUw59XPNESXd9YDeNyBJclxtn8zL8l0Uc51iu
FhROYlX5k/W/y4zIynkukBSMa8jmKbF+ie4genISrx9Q
-----END RSA PRIVATE KEY-----`
	testIncorrectPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDx6dqtEuyEf7Mpviqa/rYl316f
OoPozRgG8msH03tC9+esdfseftddfsefes8LnYF6cA7j4lBwnjOJ3Omts5CXwtVS
2rsFqvITfh0XNQq6W1JB2igTzezpybvpY3M157NImF0ijRPMcxP2qAjP7YgWjuDX
W+kIFfkbaZVWbkUYAwIDAQAB
-----END PUBLIC KEY-----`

	issuer        *auth.Issuer
	expiredIssuer *auth.Issuer
	invalidIssuer *auth.Issuer
	futureIssuer  *auth.Issuer

	parser *auth.Parser

	correctPK   *rsa.PublicKey
	incorrectPK *rsa.PublicKey

	now       time.Time
	then      time.Time
	at        time.Time
	validTime time.Duration
)

func TestMain(m *testing.M) {
	var err error
	now = time.Now()
	then = now.Add(-(76 * time.Hour))
	at = now.Add(76 * time.Hour)
	validTime = time.Hour
	jwt.TimeFunc = func() time.Time { return now }

	issuer, err = auth.NewIssuerFromPrivateKeyPEM(auth.IssuerConfig{
		Name:        "test",
		TimeFunc:    func() time.Time { return now },
		ValidPeriod: validTime,
	}, []byte(testPrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	expiredIssuer, err = auth.NewIssuerFromPrivateKeyPEM(auth.IssuerConfig{
		Name:        "test expired",
		TimeFunc:    func() time.Time { return then },
		ValidPeriod: validTime,
	}, []byte(testPrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	futureIssuer, err = auth.NewIssuerFromPrivateKeyPEM(auth.IssuerConfig{
		Name:        "test future",
		TimeFunc:    func() time.Time { return at },
		ValidPeriod: validTime,
	}, []byte(testPrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	invalidIssuer, err = auth.NewIssuerFromPrivateKeyPEM(auth.IssuerConfig{
		Name:        "test invalid",
		TimeFunc:    func() time.Time { return now },
		ValidPeriod: validTime,
	}, []byte(testIncorrectPrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	parser = issuer.Parser()
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
