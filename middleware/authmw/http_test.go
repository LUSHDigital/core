package authmw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/middleware/authmw"
	"github.com/LUSHDigital/core/rest"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/keybroker/keybrokermock"
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
			name:                 "token is missing",
			broker:               keybrokermock.MockRSAPublicKey(*correctPK),
			claims:               auth.Claims{},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "missing token",
		},
		{
			name:   "token is good",
			broker: keybrokermock.MockRSAPublicKey(*correctPK),
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
			broker: keybrokermock.MockRSAPublicKey(*correctPK),
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
			broker: keybrokermock.MockRSAPublicKey(*correctPK),
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
			broker: keybrokermock.MockRSAPublicKey(*correctPK),
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
			broker: keybrokermock.MockRSAPublicKey(*incorrectPK),
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

			// only add the token when a consumer is provided.
			if c.claims.Consumer.IsUser(defaultConsumer.ID) {
				req.Header.Add("Authorization", "Bearer "+token)
			}

			recorder := httptest.NewRecorder()
			handler := authmw.HandlerValidateJWT(c.broker, func(w http.ResponseWriter, r *http.Request) {
				consumer := auth.ConsumerFromContext(r.Context())
				rest.Response{Code: http.StatusOK, Message: "", Data: &rest.Data{Type: "consumer", Content: consumer}}.WriteTo(w)
			})
			handler.ServeHTTP(recorder, req)
			test.Equals(t, c.expectedStatusCode, recorder.Code)

			if c.expectedStatusCode == http.StatusOK {
				var consumer auth.Consumer
				rest.UnmarshalJSONResponse(recorder.Body.Bytes(), &consumer)
				test.Equals(t, c.claims.Consumer.ID, consumer.ID)
			}
		})
	}
}