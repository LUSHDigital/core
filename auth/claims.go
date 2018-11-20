package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// Claims hold the JWT claims to user for a token
type Claims struct {
	Consumer Consumer `json:"consumer"`
	jwt.StandardClaims
}
