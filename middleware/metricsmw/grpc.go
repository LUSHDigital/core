package metricsmw

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// DefaultServerMetrics is the default instance of ServerMetrics. It is
// intended to be used in conjunction the default Prometheus metrics
// registry.
var DefaultServerMetrics = grpc_prometheus.DefaultServerMetrics

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
var StreamServerInterceptor = grpc_prometheus.StreamServerInterceptor

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
var UnaryServerInterceptor = grpc_prometheus.UnaryServerInterceptor

// StreamClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
var StreamClientInterceptor = grpc_prometheus.StreamClientInterceptor

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
var UnaryClientInterceptor = grpc_prometheus.UnaryClientInterceptor
