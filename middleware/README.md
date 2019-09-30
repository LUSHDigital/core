# Middleware
The `core/middleware` package is used to interact with HTTP & gRPC middlewares.

## Middlewares
- [core/middleware/i18nmw](https://github.com/LUSHDigital/core/tree/master/middleware/i18nmw#internationalisation-middleware)
- [core/middleware/metricsmw](https://github.com/LUSHDigital/core/tree/master/middleware/metricsmw#metrics-middleware)
- [core/middleware/paginationmw](https://github.com/LUSHDigital/core/tree/master/middleware/paginationmw#pagination-middleware)
- [core/middleware/tracingmw](https://github.com/LUSHDigital/core/tree/master/middleware/tracingmw#tracing-middleware)

## Examples

### Chain gRPC middlewares

```go
server = grpcsrv.New(nil, middleware.WithUnaryServerChain(
    metricsmw.UnaryServerInterceptor,
    tracingmw.UnaryServerInterceptor,
    paginationmw.UnaryServerInterceptor,
))
```