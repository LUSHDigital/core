package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/microservice-core-golang/auth"
)

var (
	signingKey = []byte("this is my secret key, shhhh")
)

func TestHandlerValidateJWT_ValidToken(t *testing.T) {
	// create our test consumer
	consumer := auth.Consumer{
		ID:     5,
		Grants: []auth.Grant{"test.grant"},
	}

	// create a JWT for the consumer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
			NotBefore: time.Now().Add(-1 * time.Hour).Unix(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		Consumer: consumer,
	})

	// sign the JWT
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		t.Fatal(err)
	}

	// make the request, the verb and path a irrelevant
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// add our token to the request headers
	req.Header.Add("Authorization", "Bearer "+signedToken)

	// make a response writer that will record return status codes and things
	rr := httptest.NewRecorder()

	// call the handler
	handler := auth.HandlerValidateJWT(signingKey, okHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler did not return 200 OK: got %v", status)
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
