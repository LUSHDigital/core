package i18n

import (
	"context"
)

type key int

const (
	localeKey key = iota
)

// ContextWithLocale takes a context and a locale and returns a new context with the locale embedded.
func ContextWithLocale(parent context.Context, locale string) context.Context {
	return context.WithValue(parent, localeKey, locale)
}

// LocaleFromContext extracts the locale from the supplied context.
func LocaleFromContext(ctx context.Context) string {
	if c, ok := ctx.Value(localeKey).(string); ok {
		return c
	}
	return DefaultLanguage
}
