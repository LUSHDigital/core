# Metric Server
The package `core/workers/metricsrv` provides a default set of configuration for hosting a http prometheus metrics in a service.

## Configuration
The metric server can be configured through the environment to match setup in the infrastructure.

- `PROMETHEUS_INTERFACE` default: `:5117`
- `PROMETHEUS_PATH` default: `/metrics`

## Examples

### Starting server and exposing metrics

```go
srv := metricsrv.New()
srv.Run(ctx, ioutil.Discard)
```
