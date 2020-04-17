# Metrics Middleware
The package `core/workers/metricsmw` is used to record and expose metrics for an application. The metrics server is be provided over HTTP using the prometheus extraction protocol.

You can read more about [using prometheus in go on the their offical website](https://prometheus.io/docs/guides/go-application/).


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

### HTTP server metrics
Using gorilla mux.

```go
r := mux.NewRouter()
r.Use(mux.MiddlewareFunc(metrics.MeasureRequestsMiddleware))
```

Using standard net/http library.

```go
http.Handle("/check", metrics.MeasureRequests(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
}))
```

### Time series buckets
The collection of histogram metrics must put data into predefined buckets.
These have been set at: `0.2, 0.4, 0.6, 0.8, 1.0, 1.5, 2.0, 5.0, 10.0, 15.0` seconds.
more information can be found here: [Prometheus Documentation](https://prometheus.io/docs/practices/histograms/#errors-of-quantile-estimation)
