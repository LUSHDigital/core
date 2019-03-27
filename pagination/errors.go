package pagination

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrCalculateOffset is used when the pagination offset could not be calculated.
	ErrCalculateOffset = errors.New("cannot calculate offset: insufficient data")

	// ErrCalculateLastPage is used when the pagination last page could not be calculated.
	ErrCalculateLastPage = errors.New("cannot calculate last page: insufficient data")

	// ErrMetadataMissing happens when there is no metadata with the request
	ErrMetadataMissing = status.Error(codes.InvalidArgument, "metadata missing")

	// ErrMetadataInvalid happens when a metadata key is invalid or missing
	ErrMetadataInvalid = func(key string, err error) error {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("invalid or missing [%s]: %v", key, err))
	}
)
