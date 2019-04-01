package pagination_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/pagination/testdata"
	"github.com/davecgh/go-spew/spew"
	context2 "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestContext(t *testing.T) {
	ctx := pagination.ContextWithRequest(context.Background(), pagination.Request{
		Page:    1,
		PerPage: 10,
	})

	req := pagination.RequestFromContext(ctx)

	equals(t, uint64(1), req.Page)
	equals(t, uint64(10), req.PerPage)
}

func TestInterceptServerRequest(t *testing.T) {
	t.Run("valid incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"per_page": "10",
			"page":     "1",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		req, err := pagination.InterceptServerRequest(ctx)
		if err != nil {
			t.Fatal(err)
		}
		expected := pagination.Request{
			PerPage: 10,
			Page:    1,
		}
		equals(t, expected, req)
	})
	t.Run("invalid per page incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"per_page": "abc",
			"page":     "1",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		_, err := pagination.InterceptServerRequest(ctx)
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
		_, err := pagination.InterceptServerRequest(ctx)
		if err == nil {
			t.Fatal("expected an error but got none")
		}
	})
}

func TestGRPCInterceptor(t *testing.T) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(pagination.UnaryServerInterceptor),
		grpc.StreamInterceptor(pagination.StreamServerInterceptor),
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
	ctx := metadata.NewOutgoingContext(context.Background(), nil)
	ctx = pagination.ContextWithRequest(ctx, pagination.Request{
		PerPage: 10,
		Page:    1,
	})
	spew.Dump(ctx)
	if _, err := client.SayHello(ctx, &greeter.Empty{}); err != nil {
		t.Fatal(err)
	}

}

type GreeterServer struct{}

func (*GreeterServer) SayHello(ctx context2.Context, _ *greeter.Empty) (*greeter.Empty, error) {
	spew.Dump(ctx)
	req := pagination.RequestFromContext(ctx)
	if req.Page != 1 && req.PerPage != 10 {
		return &greeter.Empty{}, errors.New("failed to intercept")
	}

	return &greeter.Empty{}, nil
}
