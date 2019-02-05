package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

const (
	// TokenValidPeriod is the default amount of minutes a token is valid
	TokenValidPeriod = 60
)

var (
	// ErrTokenInvalid happens when a token could not be validated because of an unknown reason
	ErrTokenInvalid = TokenInvalidError{fmt.Errorf("token invalid")}
)

// TokenInvalidError happens when a token could not be validated because of an unknown reason
type TokenInvalidError struct{ error }

// TokenSignatureError happens when the signature could not be verified with the given public key
type TokenSignatureError struct{ error }

// TokenExpiredError happens when the token has expired or is not yet valid
type TokenExpiredError struct{ error }

// TokenMalformedError happens when the token is not the correct format
type TokenMalformedError struct{ error }

// Tokeniser is the auth tokeniser for JSON Web Tokens
type Tokeniser struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	authIssuer string
}

// NewMockTokeniser creates a new tokeniser with a random key pair
func NewMockTokeniser() (*Tokeniser, error) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey
	issuer, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return NewTokeniser(privateKey, publicKey, issuer), nil
}

// NewTokeniserFromPublicKey parses a public key to
func NewTokeniserFromPublicKey(pkb []byte) (*Tokeniser, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pkb)
	if err != nil {
		return nil, err
	}
	return &Tokeniser{publicKey: publicKey}, nil
}

// NewTokeniserFromKeyPair parses a public key to
func NewTokeniserFromKeyPair(privateKeyB, publicKeyB []byte, issuer string) (*Tokeniser, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyB)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyB)
	if err != nil {
		return nil, err
	}
	return NewTokeniser(privateKey, publicKey, issuer), nil
}

// NewTokeniser returns a new JWT instance
func NewTokeniser(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, issuer string) *Tokeniser {
	return &Tokeniser{
		privateKey: privateKey,
		publicKey:  publicKey,
		authIssuer: issuer,
	}
}

// GenerateToken generates and returns an authentication token.
func (t *Tokeniser) GenerateToken(consumer *Consumer) (string, error) {
	// Create our claims
	// NOTE: The consumer is sanitised
	claims := Claims{
		Consumer: *consumer,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(TokenValidPeriod) * time.Minute).Unix(),
			Issuer:    t.authIssuer,
			Id:        uuid.New(),
		},
	}
	// Create the token
	newToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token
	return newToken.SignedString(t.privateKey)
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

func handleParseErr(err error) error {
	if _, ok := err.(*jwt.ValidationError); ok {
		// Handle any token specific errors.
		if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorMalformed != 0 {
			return TokenMalformedError{fmt.Errorf("token malformed: %v", err)}
		} else if err.(*jwt.ValidationError).Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return TokenExpiredError{fmt.Errorf("%v", err)}
		} else if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return TokenSignatureError{fmt.Errorf("token signature invalid: %v", err)}
		} else {
			return TokenInvalidError{fmt.Errorf("token invalid: %v", err)}
		}
	}
	return err
}
