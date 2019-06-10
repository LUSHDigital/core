# Middleware
The `core/middleware` package is used to interact with HTTP & gRPC middlewares.

## Examples

### Chain gRPC middlewares

```go
server = grpcsrv.New(nil, middleware.WithUnaryServerChain(
    metricsmw.UnaryServerInterceptor,
    tracingmw.UnaryServerInterceptor,
    paginationmw.UnaryServerInterceptor,
))
```