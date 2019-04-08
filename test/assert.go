package test

import (
	"reflect"
	"testing"
)

// Equals performs a deep equal comparison against two values and fails if they are not the same.
func Equals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}

// NotEquals performs a deep equal comparison against two values and fails if they are the same.
func NotEquals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()
	if reflect.DeepEqual(expected, actual) {
		tb.Fatalf("\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n", expected, actual)
	}
}
