package paginationmw_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/LUSHDigital/core/middleware/internal/greeter"
	"github.com/LUSHDigital/core/middleware/paginationmw"
	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInterceptServerRequest(t *testing.T) {
	t.Run("valid incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"per_page": "10",
			"page":     "1",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		req, err := paginationmw.InterceptServerRequest(ctx)
		if err != nil {
			t.Fatal(err)
		}
		expected := pagination.Request{
			PerPage: 10,
			Page:    1,
		}
		test.Equals(t, expected, req)
	})
	t.Run("invalid per page incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"per_page": "abc",
			"page":     "1",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		_, err := paginationmw.InterceptServerRequest(ctx)
		if err == nil {
			t.Fatal("expected an error but got none")
		}

	})
	t.Run("invalid page incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"per_page": "10",
			"page":     "abc",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		_, err := paginationmw.InterceptServerRequest(ctx)
		if err == nil {
			t.Fatal("expected an error but got none")
		}
	})
}

func TestGRPCInterceptor(t *testing.T) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(paginationmw.UnaryServerInterceptor),
		grpc.StreamInterceptor(paginationmw.StreamServerInterceptor),
	)
	greeter.RegisterGreeterServer(server, &GreeterServer{})
	listener, err := net.Listen("tcp", "")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()
	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := greeter.NewGreeterClient(conn)
	md := metadata.New(map[string]string{
		"per_page": "10",
		"page":     "1",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	if _, err := client.SayHello(ctx, &greeter.Empty{}); err != nil {
		t.Fatal(err)
	}
}

type GreeterServer struct{}

func (*GreeterServer) SayHello(ctx context.Context, _ *greeter.Empty) (*greeter.Empty, error) {
	req := pagination.RequestFromContext(ctx)
	if req.Page != 1 || req.PerPage != 10 {
		return &greeter.Empty{}, errors.New("failed to intercept")
	}
	return &greeter.Empty{}, nil
}
