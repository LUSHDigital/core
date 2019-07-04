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
	switch exp := expected.(type) {
	case error:
		act, ok := actual.(error)
		if !ok {
			tb.Fatalf("actual value should be of type %T", actual)
		}
		if !reflect.DeepEqual(exp.Error(), act.Error()) {
			tb.Fatalf(tmplEQErr, expected, actual)
		}
	default:
		if !reflect.DeepEqual(expected, actual) {
			tb.Fatalf(tmplEQ, expected, actual)
		}
	}
}

// NotEquals performs a deep equal comparison against two values and fails if they are the same.
func NotEquals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	switch exp := expected.(type) {
	case error:
		act, ok := actual.(error)
		if !ok {
			tb.Fatalf("actual value should be of type %T", actual)
		}
		if reflect.DeepEqual(exp.Error(), act.Error()) {
			tb.Fatalf(tmplEQErr, expected, actual)
		}
	default:
		if reflect.DeepEqual(expected, actual) {
			tb.Fatalf(tmplEQ, expected, actual)
		}
	}
}
