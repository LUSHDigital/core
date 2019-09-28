package auth_test

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/core/auth"
)

type Consumer struct {
	ID string
}

type Claims struct {
	jwt.StandardClaims
	Consumer
}

var (
	testPrivateKey          = loadBytes("private.rsa.key")
	testPublicKey           = loadBytes("public.rsa.key")
	testIncorrectPrivateKey = loadBytes("incorrect.private.rsa.key")
	testIncorrectPublicKey  = loadBytes("incorrect.public.rsa.key")

	issuer        *auth.Issuer
	expiredIssuer *auth.Issuer
	invalidIssuer *auth.Issuer
	futureIssuer  *auth.Issuer

	parser        *auth.Parser
	invalidParser *auth.Parser

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
	keyFunc := func(pk crypto.PublicKey) jwt.Keyfunc {
		return func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return pk, fmt.Errorf("unknown algorithm: %v", token.Header["alg"])
			}
			return pk, nil
		}
	}

	issuer, err = auth.NewIssuerFromPEM(testPrivateKey, auth.IssuerConfig{
		Name:        "test",
		ValidPeriod: validTime,
	})
	if err != nil {
		log.Fatalln(err)
	}
	expiredIssuer, err = auth.NewIssuerFromPEM(testPrivateKey, auth.IssuerConfig{
		Name:        "test expired",
		ValidPeriod: validTime,
	})
	if err != nil {
		log.Fatalln(err)
	}
	futureIssuer, err = auth.NewIssuerFromPEM(testPrivateKey, auth.IssuerConfig{
		Name:        "test future",
		ValidPeriod: validTime,
	})
	if err != nil {
		log.Fatalln(err)
	}
	invalidIssuer, err = auth.NewIssuerFromPEM(testIncorrectPrivateKey, auth.IssuerConfig{
		Name:        "test invalid",
		ValidPeriod: validTime,
	})
	if err != nil {
		log.Fatalln(err)
	}
	parser, err = auth.NewParserFromPEM(testPublicKey, keyFunc)
	if err != nil {
		log.Fatalln(err)
	}
	invalidParser, err = auth.NewParserFromPEM(testIncorrectPublicKey, keyFunc)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}

func loadBytes(name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
