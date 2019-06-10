package test

import (
	"reflect"
	"testing"
)

const (
	tmplEQ    = "\n\texpected: %#[1]v (%[1]T)\n\t  actual: %#[2]v (%[2]T)\n"
	tmplEQErr = "\n\texpected: %[1]v (%[1]T)\n\t  actual: %[2]v (%[2]T)\n"
)

// Equals performs a deep equal comparison against two values and fails if they are not the same.
func Equals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf(tmpl(actual), expected, actual)
	}
}

// NotEquals performs a deep equal comparison against two values and fails if they are the same.
func NotEquals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if reflect.DeepEqual(expected, actual) {
		tb.Fatalf(tmpl(actual), expected, actual)
	}
}

func tmpl(t interface{}) string {
	switch t.(type) {
	case error:
		return tmplEQErr
	default:
		return tmplEQ
	}
}
