package test_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/LUSHDigital/core/test"
	"github.com/google/go-cmp/cmp"
)

type TestError struct{ error }

func TestErrorTypeComparator(t *testing.T) {
	testErr := TestError{errors.New("oops")}
	someErr := errors.New("ouch")
	opts := cmp.Options{
		test.ErrorTypeComparer,
	}
	if cmp.Equal(testErr, someErr, opts) {
		t.Fatalf("expected %T and %T to be of different types\n", someErr, testErr)
	}
}

func TestErrorReporter(t *testing.T) {
	testErr := TestError{errors.New("oops")}
	someErr := errors.New("ouch")

	var e test.ErrorReporter
	opts := cmp.Options{
		test.ErrorTypeComparer,
		cmp.Reporter(&e),
	}
	if !cmp.Equal(testErr, someErr, opts) {
		if e.String() == "" {
			t.Fatal("reporter should produce output")
		}
	}
}

func ExampleErrorTypeComparer() {
	testErr := TestError{errors.New("oops")}
	someErr := errors.New("ouch")

	var e test.ErrorReporter
	opts := cmp.Options{
		test.ErrorTypeComparer,
		cmp.Reporter(&e),
	}
	if !cmp.Equal(testErr, someErr, opts) {
		fmt.Println(e.String())
	}
	// Output:
	// error type mismatch:
	// 	expected: test_test.TestError
	// 	got: *errors.errorString

}
