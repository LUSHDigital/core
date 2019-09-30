package authmock

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/dgrijalva/jwt-go"
)

const (
	bitSize = 2048
)

// RSAKeyFunc represents the keyfunc for the mock parser.
func RSAKeyFunc(pk crypto.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return pk, fmt.Errorf("unknown algorithm: %v", token.Header["alg"])
		}
		return pk, nil
	}
}

// NewRSAIssuerAndParser creates a new issuer with a random key pair.
func NewRSAIssuerAndParser() (*auth.Issuer, *auth.Parser, error) {
	private, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	name, err := os.Hostname()
	if err != nil {
		return nil, nil, err
	}
	issuer := auth.NewIssuer(private, auth.IssuerConfig{
		Name:          name,
		ValidPeriod:   30 * time.Minute,
		SigningMethod: jwt.SigningMethodRS256,
	})
	parser := auth.NewParser(&private.PublicKey, RSAKeyFunc)
	return issuer, parser, nil
}
