package paginationmw

import (
	"context"
	"log"
	"strconv"

	"github.com/LUSHDigital/core/pagination"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Debug is used to turn on debug logging for this package.
var Debug = false

// InterceptServerRequest returns a new Response instance from the provided
// context, or returns an error if it fails, or finds the context to be invalid.
func InterceptServerRequest(ctx context.Context) (pagination.Request, error) {
	var err error
	req := pagination.Request{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return req, nil
	}
	extract := func(key string, md metadata.MD) (uint64, error) {
		val := md.Get(key)
		if len(val) < 1 {
			if Debug {
				log.Printf("grpc pagination: tried to access %q meta data key but it didn't have any values", key)
			}
			return 0, nil
		}
		n, err := strconv.ParseUint(val[0], 10, 64)
		if err != nil {
			if Debug {
				log.Printf("grpc pagination: could not parse %q key: %v", key, err)
			}
			return 0, pagination.ErrMetadataInvalid(key, err)
		}
		return n, nil
	}

	req.PerPage, err = extract("per-page", md)
	if err != nil {
		return req, err
	}

	req.Page, err = extract("page", md)
	if err != nil {
		return req, err
	}

	return req, nil
}

// UnaryServerInterceptor is a gRPC server-side interceptor that checks that
// pagination is provided is valid for unary procedures
func UnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	pr, err := InterceptServerRequest(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := handler(pagination.ContextWithRequest(ctx, pr), req)
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that checks that
// pagination provided is valid for streaming procedures
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	pr, err := InterceptServerRequest(ss.Context())
	if err != nil {
		return err
	}
	err = handler(srv, &paginatedServerStream{ss, pr})
	return err
}

type paginatedServerStream struct {
	grpc.ServerStream
	pr pagination.Request
}

func (s *paginatedServerStream) Context() context.Context {
	return pagination.ContextWithRequest(s.ServerStream.Context(), s.pr)
}
