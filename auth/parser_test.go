package auth_test

import (
	"github.com/dgrijalva/jwt-go"
)

func ExampleParser_Parse() {
	var claims jwt.StandardClaims
	err := parser.Parse(`... jwt key ...`, &claims)
	if err != nil {
		return
	}
}
