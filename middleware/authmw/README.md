# Auth Middleware
The package `core/middleware/authmw` is used to attach authentication information to requests and responses for REST and gRPC. To learn more about how to use auth inside of your application you should read the [documentation for the **core/auth** package](https://github.com/LUSHDigital/core/tree/master/auth#auth).

## Examples

### Attach gRPC auth middlewares to server

```go
server := grpc.NewServer(
    authmw.NewStreamServerInterceptor(broker),
    authmw.NewUnaryServerInterceptor(broker),
)
```
