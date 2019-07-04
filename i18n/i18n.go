package i18n

import (
	"golang.org/x/text/language"
)

var (
	// DefaultLocale is the default locale is set to english for simplicity.
	DefaultLocale = "en"
)

// ParseLocale will attempt to read the locale from a string.
func ParseLocale(s string) (string, error) {
	tag, err := language.Parse(s)
	if err != nil {
		return "", err
	}
	return tag.String(), nil
}
