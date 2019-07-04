package auth_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/uuid"
)

func mustIssue(jwt string, err error) string {
	if err != nil {
		panic(err)
	}
	return jwt
}

func TestParser_Claims(t *testing.T) {
	consumer := &auth.Consumer{
		ID:        1,
		UUID:      uuid.Must(uuid.NewV4()).String(),
		FirstName: "John",
		LastName:  "Doe",
		Language:  "en",
		Grants:    []string{},
		Roles:     []string{},
		Needs:     []string{},
	}
	cases := []struct {
		name              string
		jwt               string
		expectedErr       error
		expectedIssuedAt  int64
		expectedExpiresAt int64
		expectedNotBefore int64
	}{
		{
			name:              "valid token",
			jwt:               mustIssue(issuer.Issue(consumer)),
			expectedIssuedAt:  now.Unix(),
			expectedNotBefore: now.Unix(),
			expectedExpiresAt: now.Add(60 * time.Minute).Unix(),
		},
		{
			name:              "incorrect signature",
			jwt:               mustIssue(invalidIssuer.Issue(consumer)),
			expectedErr:       auth.TokenSignatureError{Err: fmt.Errorf("crypto/rsa: verification error")},
			expectedIssuedAt:  now.Unix(),
			expectedNotBefore: now.Unix(),
			expectedExpiresAt: now.Add(60 * time.Minute).Unix(),
		},
		{
			name:              "expired token",
			jwt:               mustIssue(expiredIssuer.Issue(consumer)),
			expectedErr:       auth.TokenExpiredError{Err: fmt.Errorf("token is expired by 75h0m0s")},
			expectedIssuedAt:  then.Unix(),
			expectedNotBefore: then.Unix(),
			expectedExpiresAt: then.Add(60 * time.Minute).Unix(),
		},
		{
			name:              "token not valid yet",
			jwt:               mustIssue(futureIssuer.Issue(consumer)),
			expectedErr:       auth.TokenExpiredError{Err: fmt.Errorf("token is not valid yet")},
			expectedIssuedAt:  at.Unix(),
			expectedNotBefore: at.Unix(),
			expectedExpiresAt: at.Add(60 * time.Minute).Unix(),
		},
		{
			name:        "unexpected signing method",
			jwt:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectedErr: auth.UnexpectedSigningMethodError{"HS256"},
		},
		{
			name:        "malformed token",
			jwt:         ".eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwiZXhwIjoxNTE4NjAzNzIwLCJqdGkiOiIyNTAwYjk3MS0wNTcxLTQ4Y2UtYmUzOS1jNWJhNGQwZmU0MGIiLCJpc3MiOiJ0ZXN0aW5nIn0.",
			expectedErr: auth.TokenMalformedError{Err: fmt.Errorf("unexpected end of JSON input")},
		},
		{
			name:        "missing token",
			expectedErr: auth.TokenMalformedError{Err: fmt.Errorf("token contains an invalid number of segments")},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			claims, err := parser.Claims(c.jwt)
			test.Equals(t, c.expectedErr, err)
			test.Equals(t, c.expectedIssuedAt, claims.IssuedAt().Unix())
			test.Equals(t, c.expectedExpiresAt, claims.ExpiresAt().Unix())
			test.Equals(t, c.expectedExpiresAt, claims.ExpiresAt().Unix())
			test.Equals(t, c.expectedNotBefore, claims.NotBefore().Unix())
		})
	}
	t.Run("parse expired and extract data for token refresh", func(t *testing.T) {
		token := mustIssue(expiredIssuer.Issue(consumer))
		claims, err := parser.Claims(token)
		test.NotEquals(t, nil, err)
		switch err.(type) {
		case auth.TokenExpiredError:
			test.Equals(t, then.Add(60*time.Minute).Unix(), claims.ExpiresAt().Unix())
		default:
			test.Equals(t, nil, err)
		}
	})
}
