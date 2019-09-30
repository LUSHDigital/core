package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Issuer represents a set of methods for generating a JWT with a private key
type Issuer struct {
	method  jwt.SigningMethod
	private crypto.PrivateKey
	name    string
	valid   time.Duration
}

type encodingFunc func(der []byte) (interface{}, error)

// NewIssuerFromPEM will take a private key PEM and derive the private key from it.
func NewIssuerFromPEM(key []byte, method jwt.SigningMethod) (*Issuer, error) {
	private, err := PrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}
	return NewIssuer(private, method), nil
}

// NewIssuerFromPEMWithPassword will take a private key PEM with a password and derive the private key from it.
func NewIssuerFromPEMWithPassword(key []byte, password string, method jwt.SigningMethod) (*Issuer, error) {
	private, err := PrivateKeyFromPEMWithPassword(key, password)
	if err != nil {
		return nil, err
	}
	return NewIssuer(private, method), nil
}

// NewIssuer creates a new issuer.
func NewIssuer(private crypto.PrivateKey, method jwt.SigningMethod) *Issuer {
	if method == nil {
		switch private.(type) {
		case *rsa.PrivateKey:
			method = jwt.SigningMethodRS256
		case *ecdsa.PrivateKey:
			method = jwt.SigningMethodES256
		}
	}
	return &Issuer{
		method:  method,
		private: private,
	}
}

// Issue will sign a JWT and return its string representation.
func (i *Issuer) Issue(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(i.method, claims)
	return token.SignedString(i.private)
}
