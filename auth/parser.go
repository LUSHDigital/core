package auth

import (
	"crypto"

	jwt "github.com/dgrijalva/jwt-go"
)

// Parser represents a set of methods for parsing and validating a JWT against a public key
type Parser struct {
	public crypto.PublicKey
	fn     PublicKeyFunc
}

// NewParser returns a new parser with a public key.
func NewParser(pk crypto.PublicKey, fn PublicKeyFunc) *Parser {
	return &Parser{public: pk, fn: fn}
}

// NewParserFromPEM will take a PEM and derive the public key from it and instantiate a parser.
func NewParserFromPEM(key []byte, fn PublicKeyFunc) (*Parser, error) {
	public, err := PublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}
	return NewParser(public, fn), nil
}

// PublicKeyFunc is used to parse tokens using a public key.
type PublicKeyFunc func(crypto.PublicKey) jwt.Keyfunc

// Parse takes a string and returns a valid jwt token
func (p *Parser) Parse(raw string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(raw, claims, p.fn(p.public))
	if err != nil {
		return nil
	}
	return nil
}
