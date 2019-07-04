package auth

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// UnexpectedSigningMethodError when JWT parsing encounters an unexpected signature method
type UnexpectedSigningMethodError struct {
	Algorithm interface{}
}

func (e UnexpectedSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method: %v", e.Algorithm)
}

// TokenInvalidError happens when a token could not be validated because of an unknown reason
type TokenInvalidError struct{ Err error }

func (e TokenInvalidError) Error() string {
	return fmt.Sprintf("token invalid: %v", e.Err)
}

// TokenSignatureError happens when the signature could not be verified with the given public key
type TokenSignatureError struct{ Err error }

func (e TokenSignatureError) Error() string {
	return fmt.Sprintf("token signature invalid: %v", e.Err)
}

// TokenExpiredError happens when the token has expired or is not yet valid
type TokenExpiredError struct{ Err error }

func (e TokenExpiredError) Error() string {
	return e.Err.Error()
}

// TokenMalformedError happens when the token is not the correct format
type TokenMalformedError struct{ Err error }

func (e TokenMalformedError) Error() string {
	return fmt.Sprintf("token malformed: %v", e.Err)
}

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
			return p.publicKey, UnexpectedSigningMethodError{token.Header["alg"]}
		}
		return p.publicKey, nil
	})
	if err != nil {
		switch err := err.(type) {
		case *jwt.ValidationError:
			switch {
			case err.Errors&jwt.ValidationErrorMalformed != 0:
				return token, TokenMalformedError{err}
			case err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				return token, TokenExpiredError{err}
			case err.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return token, TokenSignatureError{err}
			default:
				switch err.Inner.(type) {
				case UnexpectedSigningMethodError:
					return token, err
				}
				return token, TokenInvalidError{err}
			}
		default:
			return token, err
		}
	}
	return token, nil
}

// Claims returns the consumer details for a given auth token.
func (p *Parser) Claims(raw string) (*Claims, error) {
	var (
		token *jwt.Token
		err   error
	)
	token, err = p.Token(raw)
	if err != nil {
		switch err.(type) {
		case
			TokenExpiredError,
			TokenSignatureError,
			UnexpectedSigningMethodError,
			TokenInvalidError:
		default:
			return &Claims{
				StandardClaims: jwt.StandardClaims{},
				Consumer:       Consumer{},
			}, err
		}
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return claims, ErrTokenInvalid
	}
	return claims, err
}
