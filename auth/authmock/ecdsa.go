package authmock

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/dgrijalva/jwt-go"
)

// ECDSAKeyFunc represents the keyfunc for the mock parser.
func ECDSAKeyFunc(pk crypto.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return pk, fmt.Errorf("unknown algorithm: %v", token.Header["alg"])
		}
		return pk, nil
	}
}

// NewECDSAIssuerAndParser creates a new issuer with a random key pair.
func NewECDSAIssuerAndParser() (*auth.Issuer, *auth.Parser, error) {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
		SigningMethod: jwt.SigningMethodES256,
	})
	parser := auth.NewParser(&private.PublicKey, ECDSAKeyFunc)
	return issuer, parser, nil
}
