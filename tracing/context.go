package tracing

import (
	"context"
)

type key int

const (
	requestIDKey key = iota
)

// ContextWithRequestID takes a context and an *http.Request and returns a new context with the RequestID.
func ContextWithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIDKey, rid)
}

// RequestIDFromContext extracts the RequestID from the supplied context.
func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}
