package tracingmw_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/middleware/tracingmw"
	"github.com/LUSHDigital/core/test"
)

func TestContext(t *testing.T) {
	ctx := tracingmw.ContextWithRequestID(context.Background(), "1234")
	req := tracingmw.RequestIDFromContext(ctx)
	test.Equals(t, "1234", req)
}
