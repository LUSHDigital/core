# GRPC Server
The package `core/workers/grpcsrv` provides a default set of configuration for hosting a grpc server in a service.

## Examples

### Starting server and exposing the service

```go
srv := grpcsrv.New(&grpcsrv.Config{
    Addr: ":8080",
},
    grpc.StreamInterceptor(paginationmw.StreamServerInterceptor),
    grpc.UnaryInterceptor(paginationmw.UnaryServerInterceptor),
)
srv.Run(ctx, ioutil.Discard)
```
