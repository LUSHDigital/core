package i18nmw

import (
	"context"
	"strings"

	"github.com/LUSHDigital/core/i18n"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	metadataKey = "locale"
)

// NewStreamServerInterceptor will spawn a server option for internationalisation in streaming operations.
func NewStreamServerInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(StreamServerInterceptor)
}

// NewUnaryServerInterceptor will spawn a server option for internationalisation in unary operations.
func NewUnaryServerInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(UnaryServerInterceptor)
}

// InterceptLocale returns the locale from the context.
func InterceptLocale(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil
	}
	locales := md.Get(metadataKey)
	var (
		locale string
		err    error
	)
	for _, l := range locales {
		if locale, err = i18n.ParseLocale(strings.TrimSpace(l)); err == nil {
			break
		}
	}
	if locale == "" {
		locale = i18n.DefaultLocale
	}
	return locale, nil
}

// AppendLocaleToOutgoingContext will inject locale metadata to the outgoing context for use with clients.
func AppendLocaleToOutgoingContext(ctx context.Context, locale string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metadataKey, locale)
}

// UnaryServerInterceptor is a gRPC server-side interceptor that places the locale into context for unary procedures.
func UnaryServerInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	locale, err := InterceptLocale(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := handler(i18n.ContextWithLocale(ctx, locale), req)
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that places the locale into context for streaming procedures.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	locale, err := InterceptLocale(ss.Context())
	if err != nil {
		return err
	}
	err = handler(srv, &paginatedServerStream{ss, locale})
	return err
}

type paginatedServerStream struct {
	grpc.ServerStream
	locale string
}

func (s *paginatedServerStream) Context() context.Context {
	return i18n.ContextWithLocale(s.ServerStream.Context(), s.locale)
}
