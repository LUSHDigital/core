package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/microservice-core-golang/auth"
	"github.com/LUSHDigital/microservice-core-golang/response"
)

var signingKey = []byte("this is my secret key, shhhh")

func TestHandlerValidateJWT(t *testing.T) {
	tt := []struct {
		name                 string
		signingSecret        []byte
		verifiyingSecret     []byte
		issuedAt             int64
		notBefore            int64
		expiresAt            int64
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			name:                 "token is good",
			signingSecret:        signingKey,
			verifiyingSecret:     signingKey,
			issuedAt:             time.Now().Add(-2 * time.Hour).Unix(),
			notBefore:            time.Now().Add(-1 * time.Hour).Unix(),
			expiresAt:            time.Now().Add(1 * time.Hour).Unix(),
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		{
			name:                 "token has expired",
			signingSecret:        signingKey,
			verifiyingSecret:     signingKey,
			issuedAt:             time.Now().Add(-2 * time.Hour).Unix(),
			notBefore:            time.Now().Add(-1 * time.Hour).Unix(),
			expiresAt:            time.Now().Add(-1 * time.Minute).Unix(),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "token expired or not yet valid",
		},
		{
			name:                 "token is not ready yet",
			signingSecret:        signingKey,
			verifiyingSecret:     signingKey,
			issuedAt:             time.Now().Add(-2 * time.Hour).Unix(),
			notBefore:            time.Now().Add(1 * time.Minute).Unix(),
			expiresAt:            time.Now().Add(1 * time.Hour).Unix(),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "token expired or not yet valid",
		},
		{
			name:                 "issuedAt is in the future",
			signingSecret:        signingKey,
			verifiyingSecret:     signingKey,
			issuedAt:             time.Now().Add(1 * time.Hour).Unix(),
			notBefore:            time.Now().Add(1 * time.Minute).Unix(),
			expiresAt:            time.Now().Add(1 * time.Hour).Unix(),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "token expired or not yet valid",
		},
		{
			name:                 "token not signed with matching key",
			signingSecret:        signingKey,
			verifiyingSecret:     []byte("not my key"),
			issuedAt:             time.Now().Add(-2 * time.Hour).Unix(),
			notBefore:            time.Now().Add(-1 * time.Hour).Unix(),
			expiresAt:            time.Now().Add(1 * time.Hour).Unix(),
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "invalid token",
		},
	}

	for _, tc := range tt {
		// create our test consumer
		consumer := auth.Consumer{
			ID:     5,
			Grants: []auth.Grant{"test.grant"},
		}

		// create a JWT for the consumer
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  tc.issuedAt,
				NotBefore: tc.notBefore,
				ExpiresAt: tc.expiresAt,
			},
			Consumer: consumer,
		})

		// sign the JWT
		signedToken, err := token.SignedString(tc.signingSecret)
		if err != nil {
			t.Fatalf("Test '%s' failed with %v", tc.name, err)
		}

		// make the request, the verb and path a irrelevant
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Test '%s' failed with %v", tc.name, err)
		}

		// add our token to the request headers
		req.Header.Add("Authorization", "Bearer "+signedToken)

		// make a response writer that will record return status codes and things
		rr := httptest.NewRecorder()

		// call the handler
		handler := auth.HandlerValidateJWT(tc.verifiyingSecret, okHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.expectedStatusCode {
			t.Errorf("Test '%s' failed with: handler did not return %d: got %d", tc.name, tc.expectedStatusCode, status)
		}

		// unmarshall body to response
		var responseBody response.Response
		if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("Test '%s' failed with %v", tc.name, err)
		}

		// check message response
		if message := responseBody.Code; message != tc.expectedStatusCode {
			t.Errorf("Test '%s' failed with: handler did not return %d: got %d", tc.name, tc.expectedStatusCode, responseBody.Code)
		}

		if message := responseBody.Message; message != tc.expectedErrorMessage {
			t.Errorf("Test '%s' failed with: handler did not return %s: got %s", tc.name, tc.expectedErrorMessage, message)
		}

		if tc.expectedStatusCode == http.StatusOK {
			var returnedConsumer auth.Consumer
			err := responseBody.ExtractData("consumer", &returnedConsumer)
			if err != nil {
				t.Fatalf("Test '%s' failed with %v", tc.name, err)
			}

			if returnedConsumer.ID != consumer.ID {
				t.Errorf("Test '%s' failed with: consumer incorrect in response body: got %v", tc.name, returnedConsumer)
			}
		}
	}
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	consumer := auth.ConsumerFromContext(r.Context())

	response.New(http.StatusOK, "", &response.Data{Type: "consumer", Content: consumer}).WriteTo(w)
}
