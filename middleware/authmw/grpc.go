package authmw

import (
	"context"
	"log"

	"github.com/LUSHDigital/core/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	metaAuthTokenKey = "auth-token"
)

var (
	// ErrMetadataMissing happens when there is no metadata with the request
	ErrMetadataMissing = status.Error(codes.InvalidArgument, "metadata missing")

	// ErrAuthTokenMissing happens when there is no auth token in the metadata
	ErrAuthTokenMissing = status.Error(codes.InvalidArgument, "metadata missing: auth-token")
)

// NewStreamServerInterceptor creates a grpc server option with your key broker.
func NewStreamServerInterceptor(broker RSAPublicKeyCopierRenewer) grpc.ServerOption {
	return grpc.StreamInterceptor(StreamServerInterceptor(broker))
}

// NewUnaryServerInterceptor creates a unary grpc server option with your key broker.
func NewUnaryServerInterceptor(broker RSAPublicKeyCopierRenewer) grpc.ServerOption {
	return grpc.UnaryInterceptor(UnaryServerInterceptor(broker))
}

// ContextWithJWTMetadata will add a JWT to the client outgoing context metadata
func ContextWithJWTMetadata(ctx context.Context, jwt string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metaAuthTokenKey, jwt)
}

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryClientInterceptor(jwt string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = ContextWithJWTMetadata(ctx, jwt)
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

// StreamClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamClientInterceptor(jwt string) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return &authenticatedClientStream{clientStream, jwt}, nil
	}
}

type authenticatedClientStream struct {
	grpc.ClientStream
	jwt string
}

func (s *authenticatedClientStream) Context() context.Context {
	return ContextWithJWTMetadata(s.ClientStream.Context(), s.jwt)
}

// InterceptServerJWT will check the context metadata for a JWT
func InterceptServerJWT(ctx context.Context, broker RSAPublicKeyCopierRenewer) (auth.Consumer, error) {
	var consumer auth.Consumer
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return consumer, ErrMetadataMissing
	}
	tokens := md.Get(metaAuthTokenKey)
	if len(tokens) < 1 {
		return consumer, ErrAuthTokenMissing
	}
	raw := tokens[0]
	pk := broker.Copy()

	parser := auth.NewParser(&pk)
	claims, err := parser.Claims(raw)
	if err != nil {
		switch err.(type) {
		case auth.TokenMalformedError:
			return consumer, status.Error(codes.InvalidArgument, err.Error())
		case auth.TokenSignatureError:
			broker.Renew() // Renew the public key if there's an error validating the token signature
		}
		return consumer, status.Error(codes.Unauthenticated, err.Error())
	}
	return claims.Consumer, nil
}

func handleInterceptError(err error) {
	switch err {
	case ErrMetadataMissing, ErrAuthTokenMissing:
	default:
		log.Printf("grpc auth middleware error: %v\n", err)
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor that checks that JWT provided is valid for unary procedures
func UnaryServerInterceptor(broker RSAPublicKeyCopierRenewer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		consumer, err := InterceptServerJWT(ctx, broker)
		if err != nil {
			handleInterceptError(err)
		}
		resp, err := handler(auth.ContextWithConsumer(ctx, consumer), req)
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that checks that JWT provided is valid for streaming procedures
func StreamServerInterceptor(broker RSAPublicKeyCopierRenewer) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		consumer, err := InterceptServerJWT(ss.Context(), broker)
		if err != nil {
			handleInterceptError(err)
		}
		err = handler(srv, &authenticatedServerStream{ss, consumer})
		return err
	}
}

type authenticatedServerStream struct {
	grpc.ServerStream
	consumer auth.Consumer
}

func (s *authenticatedServerStream) Context() context.Context {
	return auth.ContextWithConsumer(s.ServerStream.Context(), s.consumer)
}
