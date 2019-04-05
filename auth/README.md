# Auth
The `core/auth` package provides functions for services to issue and sign api consumer tokens.
It contains several middlewares for HTTP and GRPC to aid streamlining the authentication process.

## Examples

### Put consumers through context
Setting the consumer in a context.

```go
ctx = auth.ContextWithConsumer(context.Background(), auth.Consumer{
	ID:     999,
	Grants: []string{"foo"},
})
```

Retreiving a consumer from context.

```go
consumer := auth.ConsumerFromContext(ctx)
consumer.IsUser(999)
```

### Issue new tokens

```go
consumer := &auth.Consumer{
	ID:        999,
	FirstName: "Testy",
	LastName:  "McTest",
	Grants: []string{
		"testing.read",
		"testing.create",
	},
}
raw, err := issuer.Issue(consumer)
if err != nil {
	log.Println(err)
	return
}
fmt.Println(raw)
```