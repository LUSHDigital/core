package auth_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/LUSHDigital/core/test"
	"github.com/dgrijalva/jwt-go"
)

var (
	err error
)

func ExampleIssuer_Issue() {
	claims := jwt.StandardClaims{
		Id:        "1234",
		Issuer:    "Tests",
		Audience:  "Developers",
		Subject:   "Example",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}
	raw, err := issuer.Issue(&claims)
	if err != nil {
		return
	}
	fmt.Println(raw)
}

func TestIssuer_Issue(t *testing.T) {
	raw, err := issuer.Issue(&claims)
	if err != nil {
		t.Error(err)
	}
	var parsed Claims
	err = parser.Parse(raw, &parsed)
	if err != nil {
		t.Error(err)
	}
	test.Equals(t, claims.Consumer.ID, parsed.Consumer.ID)
}
