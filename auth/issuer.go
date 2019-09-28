package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// DefaultValidPeriod is the default amount of minutes a token is valid
	DefaultValidPeriod = time.Duration(60 * time.Minute)
)

// IssuerConfig is a set of data to configure an issuer
type IssuerConfig struct {
	Name          string
	ValidPeriod   time.Duration
	SigningMethod jwt.SigningMethod
}

// Issuer represents a set of methods for generating a JWT with a private key
type Issuer struct {
	method  jwt.SigningMethod
	private crypto.PrivateKey
	name    string
	valid   time.Duration
}

type encodingFunc func(der []byte) (interface{}, error)

// NewIssuerFromPEM will take a private key PEM and derive the private key from it.
func NewIssuerFromPEM(key []byte, c IssuerConfig) (*Issuer, error) {
	private, err := PrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}
	return NewIssuer(private, c), nil
}

// NewIssuerFromPEMWithPassword will take a private key PEM with a password and derive the private key from it.
func NewIssuerFromPEMWithPassword(key []byte, password string, c IssuerConfig) (*Issuer, error) {
	private, err := PrivateKeyFromPEMWithPassword(key, password)
	if err != nil {
		return nil, err
	}
	return NewIssuer(private, c), nil
}

// NewIssuer creates a new issuer.
func NewIssuer(private crypto.PrivateKey, c IssuerConfig) *Issuer {
	method := c.SigningMethod
	if method == nil {
		switch private.(type) {
		case *rsa.PrivateKey:
			method = jwt.SigningMethodRS256
		case *ecdsa.PrivateKey:
			method = jwt.SigningMethodES256
		}
	}
	if c.ValidPeriod < time.Nanosecond {
		c.ValidPeriod = DefaultValidPeriod
	}
	return &Issuer{
		method:  method,
		private: private,
		name:    c.Name,
		valid:   c.ValidPeriod,
	}
}

// Issue will sign a JWT and return its string representation.
func (i *Issuer) Issue(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(i.method, claims)
	return token.SignedString(i.private)
}
