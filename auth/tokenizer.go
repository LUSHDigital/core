package auth

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
	"gitlab.com/LUSHDigital/soa/first-class/api-gateway/service/config"
)

// Tokeniser is the auth tokeniser for JSON Web Tokens
type Tokeniser struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	authIssuer string
}

// NewTokeniser returns a new JWT instance
func NewTokeniser(bPrivateKey, bPublicKey, issuer string) (*Tokeniser, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(bPrivateKey))
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(bPublicKey))
	if err != nil {
		return nil, err
	}
	return &Tokeniser{
		privateKey: privateKey,
		publicKey:  publicKey,
		authIssuer: issuer,
	}, nil
}

// GenerateToken generates and returns an authentication token.
func (t *Tokeniser) GenerateToken(consumer *Consumer) (token string, err error) {
	// Create our claims
	// NOTE: The consumer is sanitised
	claims := Claims{
		Consumer: *consumer,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.TokenValidPeriod) * time.Minute).Unix(),
			Issuer:    t.authIssuer,
			Id:        uuid.New(),
		},
	}
	// Create the token
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Sign the token
	token, err = newToken.SignedString(t.privateKey)
	return
}

// ValidateToken validates an authentication token and returns true/false
// based upon the result.
func (t *Tokeniser) ValidateToken(raw string) (bool, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, checkSignatureFunc(t.publicKey))
	// Bail out if the token could not be parsed
	if err != nil {
		return false, handleParseErr(err)
	}
	// Check the claims and token are valid.
	if _, ok := token.Claims.(*Claims); !ok || !token.Valid {
		return false, ErrTokenInvalid
	}

	return true, nil
}

// ParseToken takes a string and returns a valid jwt token
func (t *Tokeniser) ParseToken(raw string) (*jwt.Token, error) {
	// Parse the JWT token
	token, err := jwt.ParseWithClaims(raw, &Claims{}, checkSignatureFunc(t.publicKey))
	// Bail out if the token could not be parsed
	if err != nil {
		return nil, handleParseErr(err)
	}
	// Check the claims and token are valid
	if _, ok := token.Claims.(*Claims); !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return token, nil
}

// GetTokenConsumer returns the consumer details for a given auth token.
func (t *Tokeniser) GetTokenConsumer(raw string) *Consumer {
	token, _ := jwt.ParseWithClaims(raw, &Claims{}, checkSignatureFunc(t.publicKey))
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil
	}
	return &claims.Consumer
}

// GetTokenExpiry returns the expiry date for a given auth token.
func (t *Tokeniser) GetTokenExpiry(raw string) time.Time {
	var expiry time.Time
	token, _ := jwt.ParseWithClaims(raw, &Claims{}, checkSignatureFunc(t.publicKey))
	if claims, ok := token.Claims.(*Claims); ok {
		expiry = time.Unix(claims.ExpiresAt, 0)
	}
	return expiry
}
