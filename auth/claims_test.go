package auth_test

import (
	"testing"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/test"
	"github.com/dgrijalva/jwt-go"
)

func Test_Claims_ExpiredAt(t *testing.T) {
	ts := time.Date(2005, 01, 11, 1, 1, 1, 0, time.UTC)
	claims := &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ts.Unix(),
		},
		Consumer: auth.Consumer{},
	}
	test.Equals(t, claims.StandardClaims.ExpiresAt, claims.ExpiresAt().Unix())
}
func Test_Claims_IssuedAt(t *testing.T) {
	ts := time.Date(2005, 01, 11, 1, 1, 1, 0, time.UTC)
	claims := &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: ts.Unix(),
		},
		Consumer: auth.Consumer{},
	}
	test.Equals(t, claims.StandardClaims.IssuedAt, claims.IssuedAt().Unix())
}
func Test_Claims_NotBefore(t *testing.T) {
	ts := time.Date(2005, 01, 11, 1, 1, 1, 0, time.UTC)
	claims := &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			NotBefore: ts.Unix(),
		},
		Consumer: auth.Consumer{},
	}
	test.Equals(t, claims.StandardClaims.NotBefore, claims.NotBefore().Unix())
}
