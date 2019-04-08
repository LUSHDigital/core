package test_test

import (
	"testing"

	"github.com/LUSHDigital/core/test"
)

func TestEquals(t *testing.T) {
	testCases := []struct {
		actual   interface{}
		expected interface{}
	}{
		{
			actual:   "",
			expected: "",
		},
		{
			actual:   0,
			expected: 0,
		},
		{
			actual:   true,
			expected: true,
		},
		{
			actual:   '0',
			expected: '0',
		},
		{
			actual:   struct{}{},
			expected: struct{}{},
		},
	}
	for n, tC := range testCases {
		t.Run(string(n), func(t *testing.T) {
			test.Equals(t, tC.expected, tC.actual)
		})
	}
}

func TestNotEquals(t *testing.T) {
	testCases := []struct {
		actual   interface{}
		expected interface{}
	}{
		{
			actual:   "",
			expected: false,
		},
		{
			actual:   0,
			expected: false,
		},
		{
			actual:   true,
			expected: false,
		},
		{
			actual:   '0',
			expected: false,
		},
		{
			actual:   struct{}{},
			expected: false,
		},
	}
	for n, tC := range testCases {
		t.Run(string(n), func(t *testing.T) {
			test.NotEquals(t, tC.expected, tC.actual)
		})
	}
}
