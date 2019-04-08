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