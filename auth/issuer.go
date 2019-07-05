package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"time"

	"github.com/LUSHDigital/uuid"
	"github.com/dgrijalva/jwt-go"
)

const (
	// DefaultValidPeriod is the default amount of minutes a token is valid
	DefaultValidPeriod = time.Duration(60 * time.Minute)
)

// IssuerConfig is a set of data to configure an issuer
type IssuerConfig struct {
	Name        string
	ValidPeriod time.Duration
	TimeFunc    func() time.Time
}

// Issuer represents a set of methods for generating a JWT with a private key
type Issuer struct {
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	name        string
	validPeriod time.Duration
	timeFunc    func() time.Time
}

// NewIssuerFromPrivateKeyPEM will take a private key PEM file and return a token issuer
func NewIssuerFromPrivateKeyPEM(cfg IssuerConfig, pem []byte) (*Issuer, error) {
	pk, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return nil, err
	}
	return NewIssuer(cfg, pk), nil
}

// NewIssuer returns a new JWT instance
func NewIssuer(cfg IssuerConfig, privateKey *rsa.PrivateKey) *Issuer {
	if cfg.ValidPeriod < time.Nanosecond {
		cfg.ValidPeriod = DefaultValidPeriod
	}
	now := cfg.TimeFunc
	if now == nil {
		now = time.Now
	}
	return &Issuer{
		privateKey:  privateKey,
		publicKey:   &privateKey.PublicKey,
		name:        cfg.Name,
		validPeriod: cfg.ValidPeriod,
		timeFunc:    now,
	}
}

// NewMockIssuer creates a new issuer with a random key pair.
func NewMockIssuer() (*Issuer, error) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return NewIssuer(IssuerConfig{
		Name: name,
	}, privateKey), nil
}

// NewMockIssuerWithTime creates a new issuer with a random key pair.
func NewMockIssuerWithTime(now func() time.Time) (*Issuer, error) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return NewIssuer(IssuerConfig{
		Name:     name,
		TimeFunc: now,
	}, privateKey), nil
}

// Issue generates and returns a JWT authentication token for a private key
func (i *Issuer) Issue(consumer *Consumer) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	claims := Claims{
		Consumer: *consumer,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: i.timeFunc().Add(i.validPeriod).Unix(),
			IssuedAt:  i.timeFunc().Unix(),
			NotBefore: i.timeFunc().Unix(),
			Issuer:    i.name,
			Id:        id.String(),
		},
	}
	return i.IssueWithClaims(claims)
}

// IssueWithClaims overrides the default claims and issues a JWT token for the a private key
func (i *Issuer) IssueWithClaims(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(i.privateKey)
}

// Parser returns a parser based on the issuers private key's public counterpart
func (i *Issuer) Parser() *Parser {
	return &Parser{publicKey: i.publicKey}
}
