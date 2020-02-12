package i18nmw

import (
	"net/http"
	"strings"

	"github.com/LUSHDigital/core/i18n"
)

const (
	acceptLanguageHeader = "Accept-Language"
)

// ParseLocaleHandler will take a language from an http header and attach it to the context.
func ParseLocaleHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get(acceptLanguageHeader))
		locales := strings.Split(raw, ",")
		var (
			locale string
			err    error
		)
		for _, l := range locales {
			if locale, err = i18n.ParseLocale(strings.TrimSpace(l)); err == nil {
				break
			}
		}
		if locale == "" {
			locale = i18n.DefaultLocale
		}
		ctx := i18n.ContextWithLocale(r.Context(), locale)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// MiddlewareFunc represents a middleware func for use with gorilla mux.
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler.ServeHTTP)
}

// ParseLocaleMiddleware wraps the parse locale handler in a gorilla mux middleware.
var ParseLocaleMiddleware = MiddlewareFunc(ParseLocaleHandler)
