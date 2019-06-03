package i18nmw_test

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/LUSHDigital/core/i18n"
	"github.com/LUSHDigital/core/middleware/i18nmw"
	"github.com/LUSHDigital/core/middleware/internal/greeter"
	"github.com/LUSHDigital/core/test"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	server *grpc.Server
	conn   *grpc.ClientConn
	err    error
)

func ExampleNewStreamServerInterceptor() {
	server = grpc.NewServer(
		i18nmw.NewStreamServerInterceptor(),
		i18nmw.NewUnaryServerInterceptor(),
	)
}

func TestInterceptLocale(t *testing.T) {
	t.Run("valid incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"locale": "sv",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		locale, err := i18nmw.InterceptLocale(ctx)
		test.Equals(t, nil, err)
		test.Equals(t, "sv", locale)
	})
	t.Run("invalid locale incoming context", func(t *testing.T) {
		ctx := context.Background()
		md := metadata.New(map[string]string{
			"locale": "ABC-123",
		})
		ctx = metadata.NewIncomingContext(ctx, md)
		locale, err := i18nmw.InterceptLocale(ctx)
		test.Equals(t, nil, err)
		test.Equals(t, "en", locale)
	})
}

func TestGRPCInterceptor(t *testing.T) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(i18nmw.UnaryServerInterceptor),
		grpc.StreamInterceptor(i18nmw.StreamServerInterceptor),
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
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{}))
	ctx = i18nmw.AppendLocaleToOutgoingContext(ctx, "en")
	if _, err := client.SayHello(ctx, &greeter.Empty{}); err != nil {
		t.Fatal(err)
	}
}

type GreeterServer struct{}

func (*GreeterServer) SayHello(ctx context.Context, _ *greeter.Empty) (*greeter.Empty, error) {
	locale := i18n.LocaleFromContext(ctx)
	if locale == "" {
		return &greeter.Empty{}, errors.New("failed to intercept locale")
	}
	return &greeter.Empty{}, nil
}
