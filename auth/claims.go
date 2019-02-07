package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims hold the JWT claims to user for a token
type Claims struct {
	Consumer Consumer `json:"consumer"`
	jwt.StandardClaims
}

// ExpiresAt returns the expiry time for claims
func (c *Claims) ExpiresAt() time.Time {
	return time.Unix(c.StandardClaims.ExpiresAt, 0)
}

// IssuedAt returns the issued time for claims
func (c *Claims) IssuedAt() time.Time {
	return time.Unix(c.StandardClaims.IssuedAt, 0)
}

// NotBefore returns the issued time for claims
func (c *Claims) NotBefore() time.Time {
	return time.Unix(c.StandardClaims.NotBefore, 0)
}
