package pagination_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/pagination"
)

func TestContext(t *testing.T) {
	ctx := pagination.ContextWithRequest(context.Background(), pagination.Request{
		Page:    1,
		PerPage: 10,
	})

	req := pagination.RequestFromContext(ctx)

	equals(t, uint64(1), req.Page)
	equals(t, uint64(10), req.PerPage)
}
