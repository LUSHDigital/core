package tracing_test

import (
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/tracing"

	"google.golang.org/grpc"
)

func TestGRPCMiddleware(t *testing.T) {
	grpc.StreamInterceptor(tracing.StreamServerInterceptor)
	grpc.UnaryInterceptor(tracing.UnaryServerInterceptor)
}
