# Metrics
The package `core/metrics` is used to record and expose metrics for an application.
The metrics server is be provided over HTTP using the prometheus extraction protocol.

You can read more about [using prometheus in go on the their offical website](https://prometheus.io/docs/guides/go-application/).

## Configuration
The metrics server can be configured through the environment to match setup in the infrastructure.

- `PROMETHEUS_INTERFACE` default: `:5117`
- `PROMETHEUS_PATH` default: `/metrics`

## Examples

### Starting server and exposing metrics

```go
go func() { log.Fatal(metrics.ListenAndServe()) }()
```

### gRPC server metrics

```go
server := grpc.NewServer(
    grpc.StreamInterceptor(metrics.StreamServerInterceptor),
    grpc.UnaryInterceptor(metrics.UnaryServerInterceptor),
)
```

### gRPC client metrics

```go
conn, err = grpc.Dial(
    address,
    grpc.WithUnaryInterceptor(metrics.UnaryClientInterceptor),
    grpc.WithStreamInterceptor(metrics.StreamClientInterceptor)
)
```