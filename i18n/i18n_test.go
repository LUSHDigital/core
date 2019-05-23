package i18n_test

import (
	"fmt"
	"testing"

	"github.com/LUSHDigital/core/i18n"
	"github.com/LUSHDigital/core/test"
)

func ExampleParseLocale() {
	locale, _ := i18n.ParseLocale("EN_gb")
	fmt.Println(locale)
}

func TestParseLocale(t *testing.T) {
	locale, err := i18n.ParseLocale("EN")
	test.Equals(t, nil, err)
	test.Equals(t, "en", locale)
	locale, err = i18n.ParseLocale("EN_gB")
	test.Equals(t, nil, err)
	test.Equals(t, "en-GB", locale)
}
