package auth

import (
	"context"

	"github.com/LUSHDigital/microservice-core-golang/keys"

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

// ContextWithJWTMetadata will add a JWT to the client outgoing context metadata
func ContextWithJWTMetadata(ctx context.Context, jwt string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metaAuthTokenKey, jwt)
}

// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryClientInterceptor(jwt string) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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
func InterceptServerJWT(ctx context.Context, brk keys.RSAPublicKeyCopierRenewer) (Consumer, error) {
	var consumer Consumer
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return consumer, ErrMetadataMissing
	}
	tokens := md.Get(metaAuthTokenKey)
	if len(tokens) < 1 {
		return consumer, ErrAuthTokenMissing
	}
	raw := tokens[0]
	pk := brk.Copy()

	parser := Parser{publicKey: &pk}
	claims, err := parser.Claims(raw)
	if err != nil {
		if err != nil {
			switch err.(type) {
			case TokenMalformedError:
				return consumer, status.Error(codes.InvalidArgument, err.Error())
			case TokenSignatureError:
				brk.Renew() // Renew the public key if there's an error validating the token signature
			}
			return consumer, status.Error(codes.Unauthenticated, err.Error())
		}
	}
	return claims.Consumer, nil
}

// UnaryServerInterceptor is a gRPC server-side interceptor that checks that JWT provided is valid for unary procedures
func UnaryServerInterceptor(brk keys.RSAPublicKeyCopierRenewer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		consumer, err := InterceptServerJWT(ctx, brk)
		if err != nil {
			return nil, err
		}
		resp, err := handler(ContextWithConsumer(ctx, consumer), req)
		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that checks that JWT provided is valid for streaming procedures
func StreamServerInterceptor(brk keys.RSAPublicKeyCopierRenewer) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		consumer, err := InterceptServerJWT(ss.Context(), brk)
		if err != nil {
			return err
		}
		err = handler(srv, &authenticatedServerStream{ss, consumer})
		return err
	}
}

type authenticatedServerStream struct {
	grpc.ServerStream
	consumer Consumer
}

func (s *authenticatedServerStream) Context() context.Context {
	return ContextWithConsumer(s.ServerStream.Context(), s.consumer)
}
