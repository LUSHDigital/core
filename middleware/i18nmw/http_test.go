package i18nmw_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LUSHDigital/core/middleware/i18nmw"

	"github.com/LUSHDigital/core/i18n"

	"github.com/LUSHDigital/core/rest"
	"github.com/LUSHDigital/core/test"
)

var (
	handler http.Handler
)

func ExampleParseLocaleHandler() {
	handler = i18nmw.ParseLocaleHandler(func(w http.ResponseWriter, r *http.Request) {
		locale := i18n.LocaleFromContext(r.Context())
		fmt.Fprintln(w, locale)
	})
}

func TestHandlerValidateJWT(t *testing.T) {
	cases := []struct {
		name                 string
		accept               string
		expected             string
		expectedErrorMessage string
	}{
		{
			name:                 "locale is ok",
			accept:               "sv",
			expected:             "sv",
			expectedErrorMessage: "",
		},
		{
			name:                 "composite locale is ok",
			accept:               "sv_gb",
			expected:             "sv-GB",
			expectedErrorMessage: "",
		},
		{
			name:                 "a chain of locales are ok",
			accept:               "12-GB,sv,fi",
			expected:             "sv",
			expectedErrorMessage: "",
		},
		{
			name:     "malformatted locale prefers default locale",
			accept:   "12-GB",
			expected: "en",
		},
	}

	type content struct {
		Locale string
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Accept-Language", c.accept)

			recorder := httptest.NewRecorder()
			handler := i18nmw.ParseLocaleHandler(func(w http.ResponseWriter, r *http.Request) {
				locale := i18n.LocaleFromContext(r.Context())
				rest.Response{Code: http.StatusOK, Message: "", Data: &rest.Data{Type: "locale", Content: content{locale}}}.WriteTo(w)
			})
			handler.ServeHTTP(recorder, req)
			test.Equals(t, http.StatusOK, recorder.Code)
			var res content
			rest.UnmarshalJSONResponse(recorder.Body.Bytes(), &res)
			test.Equals(t, c.expected, res.Locale)
		})
	}
}
