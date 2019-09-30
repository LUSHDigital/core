package authmock

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

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

// NewECDSAKeyPair will create a key pair.
func NewECDSAKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, nil
}

// MustNewECDSAKeyPair will create a key pair and will panic on failure.
func MustNewECDSAKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	private, public, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}
	return private, public
}

// NewECDSAIssuerAndParser creates a new issuer with a random key pair.
func NewECDSAIssuerAndParser() (*auth.Issuer, *auth.Parser, error) {
	private, public, err := NewECDSAKeyPair()
	if err != nil {
		panic(err)
	}
	issuer, parser := NewECDSAIssuerAndParserFromKeyPair(private, public)
	return issuer, parser, nil
}

// MustNewECDSAIsserAndParser creates a new issuer and parser with a random key pair and will panic on failure.
func MustNewECDSAIsserAndParser() (*auth.Issuer, *auth.Parser) {
	issuer, parser, err := NewECDSAIssuerAndParser()
	if err != nil {
		panic(err)
	}
	return issuer, parser
}

// NewECDSAIssuerAndParserFromKeyPair creates a new issuer and parser from an ecdsa key pair.
func NewECDSAIssuerAndParserFromKeyPair(private *ecdsa.PrivateKey, public *ecdsa.PublicKey) (*auth.Issuer, *auth.Parser) {
	issuer := auth.NewIssuer(private, jwt.SigningMethodES256)
	parser := auth.NewParser(public, ECDSAKeyFunc)
	return issuer, parser
}
