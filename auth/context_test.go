package auth_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/auth"
)

func TestContext(t *testing.T) {
	ctx := auth.ContextWithConsumer(context.Background(), auth.Consumer{
		ID:     999,
		Grants: []string{"foo"},
	})
	consumer := auth.ConsumerFromContext(ctx)
	equals(t, true, consumer.IsUser(999))
}
