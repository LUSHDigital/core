# HTTP Server
The package `core/workers/httpsrv` provides a default set of configuration for hosting a http server in a service.

## Configuration
The HTTP server can be configured through these environment variables:

- `HTTP_ADDR` the HTTP server listener's network address (default: `0.0.0.0:80`)

## Examples

### Starting default server and expose the service

```go
srv := httpsrv.NewDefault(handler)
srv.Run(ctx, os.Stdout)
```

### Starting server with a custom http server

```go
srv := httpsrv.New(&http.Server{
    Handler: Handler,
    ReadTimeout: 1 * time.Second,
})
srv.Run(ctx)
```
