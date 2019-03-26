# Auth
The `core/auth` package provides functions for services to issue and sign api consumer tokens.
It contains several middlewares for HTTP and GRPC to aid streamlining the authentication process.

## Examples

### HTTP middleware

```go
r := mux.NewRouter()
r.Handle("/users", auth.HandlerValidateJWT(broker, func(w http.ResponseWriter, r *http.Request) {
	consumer := auth.ConsumerFromContext(r.Context())
	if !consumer.HasAnyGrant("users.read") {
		http.Error(w, "access denied", http.StatusUnauthorized)
	}
}))
```

### gRPC middleware

```go
srv := grpc.NewServer(
	grpc.StreamInterceptor(auth.StreamServerInterceptor(broker)),
	grpc.UnaryInterceptor(auth.UnaryServerInterceptor(broker)),
)

l, err := net.Listen("tpc", ":50051")
if err != nil {
	log.Fatalln(err)
}
log.Fatalln(srv.Serve(l))
```