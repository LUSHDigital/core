package auth

import (
	"context"
)

type key int

const (
	consumerKey key = iota
)

// ContextWithConsumer takes a context and a service consumer and returns a new context with the consumer embedded.
func ContextWithConsumer(parent context.Context, consumer Consumer) context.Context {
	return context.WithValue(parent, consumerKey, consumer)
}

// ConsumerFromContext extracts the consumer from the supplied context.
func ConsumerFromContext(ctx context.Context) Consumer {
	return ctx.Value(consumerKey).(Consumer)
}
