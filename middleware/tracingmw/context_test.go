package tracingmw_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/middleware/tracingmw"
)

func TestContext(t *testing.T) {
	ctx := tracingmw.ContextWithRequestID(context.Background(), "1234")
	req := tracingmw.RequestIDFromContext(ctx)
	equals(t, "1234", req)
}
