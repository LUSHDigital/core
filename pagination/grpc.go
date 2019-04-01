package pagination

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InterceptServerRequest returns a new Response instance from the provided
// context, or returns an error if it fails, or finds the context to be invalid.
func InterceptServerRequest(ctx context.Context) (Request, error) {
	var (
		req Request
		err error
	)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return req, nil
	}
	extract := func(key string, md metadata.MD) (uint64, error) {
		val := md.Get(key)
		if len(val) < 1 {
			return 0, nil
		}
		n, err := strconv.ParseUint(val[0], 10, 64)
		if err != nil {
			return 0, ErrMetadataInvalid(key, err)
		}
		return n, nil
	}

	req.PerPage, err = extract("per_page", md)
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
	resp, err := handler(ContextWithRequest(ctx, pr), req)
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
	pr Request
}

func (s *paginatedServerStream) Context() context.Context {
	return ContextWithRequest(s.ServerStream.Context(), s.pr)
}
