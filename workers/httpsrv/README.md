# HTTP Server
The package `core/workers/httpsrv` provides a default set of configuration for hosting a http server in a service.

## Examples

### Starting server and expose the service

```go
srv := httpsrv.New(handler)
srv.Run(ctx, os.Stdout)
```

### Starting server with a custom http server

```
srv := httpsrv.New(handler, &http.Server{
    ReadTimeout: 1 * time.Second,
})
srv.Run(ctx, os.Stdout)
```
