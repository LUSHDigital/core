package authmock

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

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

// NewRSAKeyPair will create a key pair.
func NewRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	return private, &private.PublicKey, nil
}

// MustNewRSAKeyPair will create a key pair and will panic on failure.
func MustNewRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	private, public, err := NewRSAKeyPair()
	if err != nil {
		panic(err)
	}
	return private, public
}

// NewRSAIssuerAndParser creates a new issuer with a random key pair.
func NewRSAIssuerAndParser() (*auth.Issuer, *auth.Parser, error) {
	private, public, err := NewRSAKeyPair()
	if err != nil {
		return nil, nil, err
	}
	issuer, parser := NewRSAIssuerAndParserFromKeyPair(private, public)
	return issuer, parser, nil
}

// MustNewRSAIsserAndParser creates a new issuer and parser with a random key pair and will panic on failure.
func MustNewRSAIsserAndParser() (*auth.Issuer, *auth.Parser) {
	issuer, parser, err := NewRSAIssuerAndParser()
	if err != nil {
		panic(err)
	}
	return issuer, parser
}

// NewRSAIssuerAndParserFromKeyPair creates a new issuer and parser from an rsa key pair.
func NewRSAIssuerAndParserFromKeyPair(private *rsa.PrivateKey, public *rsa.PublicKey) (*auth.Issuer, *auth.Parser) {
	issuer := auth.NewIssuer(private, jwt.SigningMethodRS256)
	parser := auth.NewParser(public, RSAKeyFunc)
	return issuer, parser
}
