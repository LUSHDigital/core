package tracing_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/tracing"
)

func TestContext(t *testing.T) {
	ctx := tracing.ContextWithRequestID(context.Background(), "1234")
	req := tracing.RequestIDFromContext(ctx)
	equals(t, "1234", req)
}
