package auth

import (
	"context"
	"net/http"
)

type key int

const (
	consumerKey key = iota
	requestIDKey
)

// ContextWithConsumer takes a context and a service consumer and returns a new context with the consumer embedded.
func ContextWithConsumer(parent context.Context, consumer Consumer) context.Context {
	return context.WithValue(parent, consumerKey, consumer)
}

// ConsumerFromContext extracts the consumer from the supplied context.
func ConsumerFromContext(c context.Context) Consumer {
	return c.Value(consumerKey).(Consumer)
}

// NewContextWithRequestID takes a context and an *http.Request and returns a new context with the RequestID.
func NewContextWithRequestID(c context.Context, r *http.Request) context.Context {
	return context.WithValue(c, requestIDKey, r.Header.Get("x-request-id"))
}

// RequestIDFromContext extracts the RequestID from the supplied context.
func RequestIDFromContext(c context.Context) string {
	return c.Value(requestIDKey).(string)
}
