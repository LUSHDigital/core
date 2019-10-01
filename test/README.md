# Test
The `core/test` package contains helpers for aiding testing.

## Examples

### Equals

```go
t.Run("foo equals foo", func(t *testing.T) {
    test.Equals(t, "foo", "foo")
})
```

### Not equals

```go
t.Run("foo does not equal foo", func(t *testing.T) {
    test.NotEquals(t, "foo", "bar")
})
```

### Using the ErrorTypeComparer with `cmp`

This comparer relies on the use of the excellent [go-cmp](https://github.com/google/go-cmp) library.

As it's documentation states, it is intended to be a more powerful and safer alternative to reflect.DeepEqual for comparing whether two values are semantically equal.

See [GoDoc documentation](https://godoc.org/github.com/google/go-cmp/cmp) for more information.

```go
t.Run("some error is not a test error", func(t *testing.T) {
	// some fake errors, for the purpose of this example.
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
})

```
