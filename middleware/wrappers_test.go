package middleware_test

import (
	"context"
	"testing"

	"github.com/LUSHDigital/core/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LUSHDigital/core/middleware"
)

func TestWrapServerStream(t *testing.T) {
	ctx := context.WithValue(context.TODO(), "something", 1)
	fake := &fakeServerStream{ctx: ctx}
	wrapped := middleware.WrapServerStream(fake)
	t.Run("values from fake must propagate to wrapper", func(t *testing.T) {
		test.NotEquals(t, nil, wrapped.Context().Value("something"))
	})
	wrapped.WrappedContext = context.WithValue(wrapped.Context(), "other", 2)
	t.Run("values from wrapper must be set", func(t *testing.T) {
		test.NotEquals(t, nil, wrapped.Context().Value("other"))
	})
}

type fakeServerStream struct {
	grpc.ServerStream
	ctx         context.Context
	recvMessage interface{}
	sentMessage interface{}
}

func (f *fakeServerStream) Context() context.Context {
	return f.ctx
}

func (f *fakeServerStream) SendMsg(m interface{}) error {
	if f.sentMessage != nil {
		return status.Errorf(codes.AlreadyExists, "fakeServerStream only takes one message, sorry")
	}
	f.sentMessage = m
	return nil
}

func (f *fakeServerStream) RecvMsg(m interface{}) error {
	if f.recvMessage == nil {
		return status.Errorf(codes.NotFound, "fakeServerStream has no message, sorry")
	}
	return nil
}

type fakeClientStream struct {
	grpc.ClientStream
}
