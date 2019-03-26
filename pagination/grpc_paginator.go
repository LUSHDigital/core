// Package pagination defines a paginator able to return formatted responses
// enabling the API consumer to retrieve data in defined chunks
package pagination

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
)

const consumerKey = iota

// ErrMetadataMissing happens when there is no metadata with the request
var ErrMetadataMissing = status.Error(codes.InvalidArgument, "metadata missing")

// ErrMetadataInvalid happens when
var ErrMetadataInvalid = func(key string, err error) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid or missing [%s]: %v", key, err))
}

// InterceptServerPagination returns a new Response instance from the provided
// context, or returns an error if it fails, or finds the context to be invalid.
func InterceptServerPagination(ctx context.Context) (paginator Request, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return paginator, ErrMetadataMissing
	}

	var (
		perPage, page uint64
	)

	extract := func(key string, md metadata.MD) (uint64, error) {
		val := md.Get(key)
		if len(val) < 1 {
			return 0, ErrMetadataInvalid(key, errors.New("invalid length"))
		}
		n, err := strconv.ParseUint(val[0], 10, 64)
		if err != nil {
			return 0, ErrMetadataInvalid(key, err)
		}
		return n, nil
	}

	perPage, err = extract("per_page", md)
	if err != nil {
		return paginator, err
	}
	page, err = extract("page", md)
	if err != nil {
		return paginator, err
	}
	return Request{PerPage: perPage, Page: page}, nil
}

// UnaryServerInterceptor is a gRPC server-side interceptor that checks that JWT
// provided is valid for unary procedures
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	paginator, err := InterceptServerPagination(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := handler(ContextWithPaginator(ctx, paginator), req)
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that checks that JWT provided is valid for streaming procedures
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	paginator, err := InterceptServerPagination(ss.Context())
	if err != nil {
		return err
	}
	err = handler(srv, &paginatedServerStream{ss, paginator})
	return err
}

// ContextWithPaginator takes a context and a service consumer and returns a new context with the consumer embedded.
func ContextWithPaginator(parent context.Context, req Request) context.Context {
	return context.WithValue(parent, consumerKey, req)
}

// PaginatorFromContext extracts the consumer from the supplied context.
func PaginatorFromContext(ctx context.Context) Response {
	if p, ok := ctx.Value(consumerKey).(Response); ok {
		return p
	}
	return Response{}
}

type paginatedServerStream struct {
	grpc.ServerStream
	paginator Request
}

func (s *paginatedServerStream) Context() context.Context {
	return ContextWithPaginator(s.ServerStream.Context(), s.paginator)
}
