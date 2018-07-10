package pagination

import "errors"

var (
	// ErrCalculateOffset is used when the pagination offset could not be calculated.
	ErrCalculateOffset = errors.New("cannot calculate offset: insufficient data")

	// ErrCalculateLastPage is used when the pagination last page could not be calculated.
	ErrCalculateLastPage = errors.New("cannot calculate last page: insufficient data")
)
