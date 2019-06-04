package middleware_test

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/LUSHDigital/core/middleware"
	"github.com/LUSHDigital/core/test"
)

var (
	someServiceName  = "SomeService.StreamMethod"
	parentUnaryInfo  = &grpc.UnaryServerInfo{FullMethod: someServiceName}
	parentStreamInfo = &grpc.StreamServerInfo{
		FullMethod:     someServiceName,
		IsServerStream: true,
	}
	someValue     = 1
	parentContext = context.WithValue(context.TODO(), "parent", someValue)
)

func TestChainUnaryServer(t *testing.T) {
	input := "input"
	output := "output"

	first := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requireContextValue(t, ctx, "parent")
		test.Equals(t, parentUnaryInfo, info)
		ctx = context.WithValue(ctx, "first", 1)
		return handler(ctx, req)
	}
	second := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requireContextValue(t, ctx, "parent")
		requireContextValue(t, ctx, "first")
		test.Equals(t, parentUnaryInfo, info)
		ctx = context.WithValue(ctx, "second", 1)
		return handler(ctx, req)
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		test.Equals(t, input, req)
		requireContextValue(t, ctx, "parent")
		requireContextValue(t, ctx, "first")
		requireContextValue(t, ctx, "second")
		return output, nil
	}

	chain := middleware.ChainUnaryServer(first, second)
	out, _ := chain(parentContext, input, parentUnaryInfo, handler)
	test.Equals(t, output, out)
}

func TestChainStreamServer(t *testing.T) {
	someService := &struct{}{}
	recvMessage := "received"
	sentMessage := "sent"
	outputError := fmt.Errorf("some error")

	first := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		requireContextValue(t, stream.Context(), "parent")
		test.Equals(t, parentStreamInfo, info)
		test.Equals(t, someService, srv)
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = context.WithValue(stream.Context(), "first", 1)
		return handler(srv, wrapped)
	}
	second := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		requireContextValue(t, stream.Context(), "parent")
		requireContextValue(t, stream.Context(), "parent")
		test.Equals(t, parentStreamInfo, info)
		test.Equals(t, someService, srv)
		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = context.WithValue(stream.Context(), "second", 1)
		return handler(srv, wrapped)
	}
	handler := func(srv interface{}, stream grpc.ServerStream) error {
		test.Equals(t, someService, srv)
		requireContextValue(t, stream.Context(), "parent")
		requireContextValue(t, stream.Context(), "first")
		requireContextValue(t, stream.Context(), "second")
		test.Equals(t, nil, stream.RecvMsg(recvMessage))
		test.Equals(t, nil, stream.SendMsg(sentMessage))
		return outputError
	}
	fakeStream := &fakeServerStream{ctx: parentContext, recvMessage: recvMessage}
	chain := middleware.ChainStreamServer(first, second)
	err := chain(someService, fakeStream, parentStreamInfo, handler)
	test.Equals(t, outputError, err)
	test.Equals(t, sentMessage, fakeStream.sentMessage)
}

func TestChainUnaryClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	reqMessage := "request"
	replyMessage := "reply"
	outputError := fmt.Errorf("some error")

	first := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requireContextValue(t, ctx, "parent")
		test.Equals(t, someServiceName, method)
		test.Equals(t, len(opts), 1)
		wrappedCtx := context.WithValue(ctx, "first", 1)
		return invoker(wrappedCtx, method, req, reply, cc, opts...)
	}
	second := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requireContextValue(t, ctx, "parent")
		test.Equals(t, someServiceName, method)
		test.Equals(t, len(opts), 1)
		wrappedCtx := context.WithValue(ctx, "second", 1)
		return invoker(wrappedCtx, method, req, reply, cc, append(opts, &grpc.EmptyCallOption{})...)
	}
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		test.Equals(t, someServiceName, method)
		requireContextValue(t, ctx, "parent")
		requireContextValue(t, ctx, "first")
		requireContextValue(t, ctx, "second")
		test.Equals(t, len(opts), 2)
		return outputError
	}
	chain := middleware.ChainUnaryClient(first, second)
	err := chain(parentContext, someServiceName, reqMessage, replyMessage, nil, invoker, parentOpts...)
	test.Equals(t, outputError, err)
}

func TestChainStreamClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	clientStream := &fakeClientStream{}
	fakeStreamDesc := &grpc.StreamDesc{ClientStreams: true, ServerStreams: true, StreamName: someServiceName}

	first := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requireContextValue(t, ctx, "parent")
		test.Equals(t, someServiceName, method)
		test.Equals(t, len(opts), 1)
		wrappedCtx := context.WithValue(ctx, "first", 1)
		return streamer(wrappedCtx, desc, cc, method, opts...)
	}
	second := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		requireContextValue(t, ctx, "parent")
		test.Equals(t, someServiceName, method)
		test.Equals(t, len(opts), 1)
		wrappedCtx := context.WithValue(ctx, "second", 1)
		return streamer(wrappedCtx, desc, cc, method, append(opts, &grpc.EmptyCallOption{})...)
	}
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		test.Equals(t, someServiceName, method)
		test.Equals(t, fakeStreamDesc, desc)
		requireContextValue(t, ctx, "parent")
		requireContextValue(t, ctx, "first")
		requireContextValue(t, ctx, "second")
		test.Equals(t, len(opts), 2)
		return clientStream, nil
	}
	chain := middleware.ChainStreamClient(first, second)
	someStream, err := chain(parentContext, fakeStreamDesc, nil, someServiceName, streamer, parentOpts...)
	test.Equals(t, nil, err)
	test.Equals(t, clientStream, someStream)
}

func requireContextValue(t *testing.T, ctx context.Context, key string) {
	val := ctx.Value(key)
	test.NotEquals(t, nil, val)
	test.Equals(t, someValue, val)
}
