package auth_test

import (
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/uuid"
	"github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/auth/authmock"
)

type Consumer struct {
	ID string
}

type Claims struct {
	jwt.StandardClaims
	Consumer
}

var (
	issuer, invalidIssuer *auth.Issuer
	parser, invalidParser *auth.Parser

	claims, expiredClaims, futureClaims, invalidClaims Claims

	now       time.Time
	then      time.Time
	at        time.Time
	validTime time.Duration
)

func TestMain(m *testing.M) {
	now = time.Now()
	then = now.Add(-(76 * time.Hour))
	at = now.Add(76 * time.Hour)
	validTime = time.Hour
	jwt.TimeFunc = func() time.Time { return now }
	issuer, parser = authmock.MustNewRSAIsserAndParser()
	invalidIssuer, invalidParser = authmock.MustNewRSAIsserAndParser()
	claims = Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.Must(uuid.NewV4()).String(),
			Issuer:    "Auth Test",
			ExpiresAt: at.Unix(),
			IssuedAt:  then.Unix(),
			NotBefore: then.Unix(),
		},
		Consumer: Consumer{
			ID: uuid.Must(uuid.NewV4()).String(),
		},
	}
	os.Exit(m.Run())
}
