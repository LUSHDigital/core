package auth

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// UnexpectedSigningMethodError when JWT parsing encounters an unexpected signature method
type UnexpectedSigningMethodError struct {
	alg interface{}
}

func (e UnexpectedSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", e.alg)
}

// TokenInvalidError happens when a token could not be validated because of an unknown reason
type TokenInvalidError struct{ error }

// TokenSignatureError happens when the signature could not be verified with the given public key
type TokenSignatureError struct{ error }

// TokenExpiredError happens when the token has expired or is not yet valid
type TokenExpiredError struct{ error }

// TokenMalformedError happens when the token is not the correct format
type TokenMalformedError struct{ error }

var (
	// ErrTokenInvalid happens when a token could not be validated because of an unknown reason
	ErrTokenInvalid = TokenInvalidError{fmt.Errorf("token invalid")}
)

// Parser represents a set of methods for parsing and validating a JWT against a public key
type Parser struct {
	publicKey *rsa.PublicKey
}

// NewParser returns a new parser with a public key.
func NewParser(pk *rsa.PublicKey) *Parser {
	return &Parser{publicKey: pk}
}

// NewParserFromPublicKeyPEM parses a public key to
func NewParserFromPublicKeyPEM(pkb []byte) (*Parser, error) {
	pk, err := jwt.ParseRSAPublicKeyFromPEM(pkb)
	if err != nil {
		return nil, err
	}
	return &Parser{publicKey: pk}, nil
}

// Token takes a string and returns a valid jwt token
func (p *Parser) Token(raw string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, UnexpectedSigningMethodError{token.Header["alg"]}
		}
		return p.publicKey, nil
	})
	if err != nil {
		switch err := err.(type) {
		case *jwt.ValidationError:
			switch {
			case err.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformedError{fmt.Errorf("token malformed: %v", err)}
			case err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				return nil, TokenExpiredError{fmt.Errorf("%v", err)}
			case err.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, TokenSignatureError{fmt.Errorf("token signature invalid: %v", err)}
			default:
				return nil, TokenInvalidError{fmt.Errorf("token invalid: %v", err)}
			}
		}
		return nil, err
	}
	// Check the claims and token are valid
	if _, ok := token.Claims.(*Claims); !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return token, nil
}

// Claims returns the consumer details for a given auth token.
func (p *Parser) Claims(raw string) (*Claims, error) {
	token, err := p.Token(raw)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}
