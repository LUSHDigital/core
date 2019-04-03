package tracingmw_test

import (
	"testing"

	"github.com/LUSHDigital/core/middleware/tracingmw"

	"google.golang.org/grpc"
)

func TestGRPCMiddleware(t *testing.T) {
	grpc.StreamInterceptor(tracingmw.StreamServerInterceptor)
	grpc.UnaryInterceptor(tracingmw.UnaryServerInterceptor)
}
