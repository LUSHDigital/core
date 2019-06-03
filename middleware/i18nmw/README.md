# Internationalisation Middleware
The package `core/middleware/i18nmw` is used to capture internationalisation parameters from HTTP headers or GRPC metadata.

### gRPC server middleware

```go
server = grpc.NewServer(
    i18nmw.NewStreamServerInterceptor(),
    i18nmw.NewUnaryServerInterceptor(),
)
```

### HTTP server middleware
Using gorilla mux.

```go
r := mux.NewRouter()
r.Use(i18nmw.ParseLocaleMiddleware)
```

Using standard `net/http` library.

```go
handler = i18nmw.ParseLocaleHandler(func(w http.ResponseWriter, r *http.Request) {
    locale := i18n.LocaleFromContext(r.Context())
    fmt.Fprintln(w, locale)
})
```
