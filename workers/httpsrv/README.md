# HTTP Server
The package `core/workers/httpsrv` provides a default set of configuration for hosting a http server in a service.

## Examples

### Starting default server and expose the service

```go
srv := httpsrv.NewDefault(handler)
srv.Run(ctx, os.Stdout)
```

### Starting server with a custom http server

```
srv := httpsrv.New&http.Server{
    Handler: Handler,
    ReadTimeout: 1 * time.Second,
})
srv.Run(ctx, os.Stdout)
```
