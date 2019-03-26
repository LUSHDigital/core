package tracing

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	metaRequestIDKey = "request-id"
)

var (
	// ErrNewRequestID happens when the application failed to generate a new
	// request id
	ErrNewRequestID = status.Error(codes.InvalidArgument, "request id could not be generated")

	// ErrMetadataMissing happens when there is no metadata with the request
	ErrMetadataMissing = status.Error(codes.InvalidArgument, "metadata missing")
)

// InterceptServerRequestID will derive a request id from the context
func InterceptServerRequestID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMetadataMissing
	}
	rids := md.Get(metaRequestIDKey)
	if len(rids) < 1 {
		id, err := uuid.NewV4()
		if err != nil {
			return "", ErrNewRequestID
		}
		return id.String(), nil
	}
	return rids[0], nil
}

// UnaryServerInterceptor is a gRPC server-side unary interceptor that checks
// that there is a request ID and ensures one gets set
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	rid, err := InterceptServerRequestID(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ContextWithRequestID(ctx, rid), req)
}

// StreamServerInterceptor is a gRPC server-side streaming interceptor that checks
// that there is a request ID and ensures one gets set
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	rid, err := InterceptServerRequestID(ss.Context())
	if err != nil {
		return err
	}
	err = handler(srv, &ridServerStream{ss, rid})
	return err
}

type ridServerStream struct {
	grpc.ServerStream
	rid string
}

func (ss *ridServerStream) Context() context.Context {
	return ContextWithRequestID(ss.ServerStream.Context(), ss.rid)
}
