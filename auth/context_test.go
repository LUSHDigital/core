package auth_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/auth"
)

var (
	ctx context.Context
)

func ExampleContextWithConsumer() {
	ctx = auth.ContextWithConsumer(context.Background(), auth.Consumer{
		ID:     999,
		Grants: []string{"foo"},
	})
}

func ExampleConsumerFromContext() {
	consumer := auth.ConsumerFromContext(ctx)
	consumer.IsUser(999)
}

func TestContext(t *testing.T) {
	ctx = auth.ContextWithConsumer(context.Background(), auth.Consumer{
		ID:     999,
		Grants: []string{"foo"},
	})
	consumer := auth.ConsumerFromContext(ctx)
	equals(t, true, consumer.IsUser(999))
}
