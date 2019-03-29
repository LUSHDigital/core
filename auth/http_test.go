package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/keys/keysmock"
	"github.com/LUSHDigital/core/response"
	jwt "github.com/dgrijalva/jwt-go"
)

func TestHandlerValidateJWT(t *testing.T) {
	defaultConsumer := auth.Consumer{
		ID:     5,
		Grants: []string{"test.grant"},
	}

	cases := []struct {
		name                 string
		broker               auth.RSAPublicKeyCopierRenewer
		claims               auth.Claims
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			name:   "token is good",
			broker: keysmock.MockRSAPublicKey(*correctPK),
			claims: auth.Claims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
					NotBefore: time.Now().Add(-1 * time.Hour).Unix(),
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
				Consumer: defaultConsumer,
			},
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		{
			name:   "token has expired",
			broker: keysmock.MockRSAPublicKey(*correctPK),
			claims: auth.Claims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
					NotBefore: time.Now().Add(-1 * time.Hour).Unix(),
					ExpiresAt: time.Now().Add(-1 * time.Minute).Unix(),
				},
				Consumer: defaultConsumer,
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "token is expired by 1m0s",
		},
		{
			name:   "token is not ready yet",
			broker: keysmock.MockRSAPublicKey(*correctPK),
			claims: auth.Claims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
					NotBefore: time.Now().Add(1 * time.Minute).Unix(),
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
				Consumer: defaultConsumer,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:   "issuedAt is in the future",
			broker: keysmock.MockRSAPublicKey(*correctPK),
			claims: auth.Claims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Add(1 * time.Hour).Unix(),
					NotBefore: time.Now().Add(1 * time.Minute).Unix(),
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
				Consumer: defaultConsumer,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:   "token not signed with matching key",
			broker: keysmock.MockRSAPublicKey(*incorrectPK),
			claims: auth.Claims{
				StandardClaims: jwt.StandardClaims{
					IssuedAt:  time.Now().Add(-2 * time.Hour).Unix(),
					NotBefore: time.Now().Add(-1 * time.Hour).Unix(),
					ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
				},
				Consumer: defaultConsumer,
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token, err := issuer.IssueWithClaims(c.claims)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Authorization", "Bearer "+token)

			recorder := httptest.NewRecorder()
			handler := auth.HandlerValidateJWT(c.broker, func(w http.ResponseWriter, r *http.Request) {
				consumer := auth.ConsumerFromContext(r.Context())
				response.Response{Code: http.StatusOK, Message: "", Data: &response.Data{Type: "consumer", Content: consumer}}.WriteTo(w)
			})
			handler.ServeHTTP(recorder, req)
			equals(t, c.expectedStatusCode, recorder.Code)

			if c.expectedStatusCode == http.StatusOK {
				var consumer auth.Consumer
				response.UnmarshalResponseJSON(recorder.Body.Bytes(), &consumer)
				equals(t, c.claims.Consumer.ID, consumer.ID)
			}
		})
	}
}
