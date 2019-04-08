package pagination_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/test"
)

func TestContext(t *testing.T) {
	ctx := pagination.ContextWithRequest(context.Background(), pagination.Request{
		Page:    1,
		PerPage: 10,
	})

	req := pagination.RequestFromContext(ctx)

	test.Equals(t, uint64(1), req.Page)
	test.Equals(t, uint64(10), req.PerPage)
}
