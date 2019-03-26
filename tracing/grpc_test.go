package tracing_test

import (
	"testing"

	"github.com/LUSHDigital/core/tracing"
	"google.golang.org/grpc"
)

func TestGRPCMiddleware(t *testing.T) {
	grpc.StreamInterceptor(tracing.StreamServerInterceptor)
	grpc.UnaryInterceptor(tracing.UnaryServerInterceptor)
}
