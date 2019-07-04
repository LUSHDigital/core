package i18n_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/i18n"
	"github.com/LUSHDigital/core/test"
)

var (
	ctx    context.Context
	locale string
)

func ExampleContextWithLocale() {
	ctx = i18n.ContextWithLocale(context.Background(), "sv")
}

func ExampleLocaleFromContext() {
	locale = i18n.LocaleFromContext(ctx)
}

func TestContext(t *testing.T) {
	ctx = i18n.ContextWithLocale(context.Background(), "sv")
	locale = i18n.LocaleFromContext(ctx)
	test.Equals(t, "sv", locale)
	i18n.DefaultLocale = "es"
	locale = i18n.LocaleFromContext(context.Background())
	test.Equals(t, "es", locale)
	i18n.DefaultLocale = "en"
	locale = i18n.LocaleFromContext(context.Background())
	test.Equals(t, "en", locale)
}
