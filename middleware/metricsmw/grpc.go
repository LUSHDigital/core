package metricsmw

import (
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// DefaultServerMetrics is the default instance of ServerMetrics. It is
// intended to be used in conjunction the default Prometheus metrics
// registry.
var DefaultServerMetrics = grpcprometheus.DefaultServerMetrics

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
var StreamServerInterceptor = grpcprometheus.StreamServerInterceptor

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
var UnaryServerInterceptor = grpcprometheus.UnaryServerInterceptor

// StreamClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
var StreamClientInterceptor = grpcprometheus.StreamClientInterceptor

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
var UnaryClientInterceptor = grpcprometheus.UnaryClientInterceptor

func init() {
	// buckets are the values in seconds that represent the upper limit of the histogram buckets.
	// These are set to focus the collection around 1 second whilst providing scope for large values.
	buckets := []float64{0.2,0.4,0.6,0.8,1,1.5,2,5,10,15}
	grpcprometheus.EnableHandlingTimeHistogram(grpcprometheus.WithHistogramBuckets(buckets))
}
