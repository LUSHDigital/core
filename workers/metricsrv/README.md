# Metric Server
The package `core/workers/metricsrv` provides a default set of configuration for hosting a http prometheus metrics in a service.

## Configuration
The metric server can be configured through the environment to match setup in the infrastructure.

- `PROMETHEUS_ADDR` default: `:5117`
- `PROMETHEUS_PATH` default: `/metrics`

## Examples

### Starting server and exposing metrics

```go
srv := metricsrv.New(&metricsrv.Config{
    Server: &http.Server{
        Addr: "0.0.0.0:5117",
    },
    Path: "/metrics",
})
srv.Run(ctx, ioutil.Discard)
```
