package pagination

import (
	"context"
)

type key int

const (
	consumerKey key = iota
)

// ContextWithRequest takes a context and a pagination request and returns
// a new context with the pagination request embedded.
func ContextWithRequest(parent context.Context, req Request) context.Context {
	return context.WithValue(parent, consumerKey, req)
}

// RequestFromContext extracts the pagination request from the supplied context.
func RequestFromContext(ctx context.Context) Request {
	if p, ok := ctx.Value(consumerKey).(Request); ok {
		return p
	}
	return Request{}
}
