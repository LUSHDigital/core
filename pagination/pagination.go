package pagination

import "math"

// Request derives the requested pagination data given by the client.
type Request struct {
	PerPage, Page uint64
}

// Offset calculates the offset from the provided pagination request.
func (p Request) Offset() uint64 {
	return (p.Page - 1) * p.PerPage
}

// Response manages pagination of a data set.
type Response struct {
	PerPage  uint64 `json:"per_page"`  // The number of items per page.
	Page     uint64 `json:"page"`      // Which page are we on?
	Offset   uint64 `json:"offset"`    // The current offset to pass to the query.
	Total    uint64 `json:"total"`     // The total number of items
	LastPage uint64 `json:"last_page"` // The number of the last possible page.
}

// MakeResponse returns a new Response with the provided
// parameters set.
func MakeResponse(request Request, total uint64) Response {
	// Avoid division by zero error, which would wrap uint to it's max value.
	var lastPage uint64
	if request.PerPage > 0 {
		lastPage = uint64(math.Ceil(float64(total) / float64(request.PerPage)))
	}

	return Response{
		PerPage:  request.PerPage,
		Page:     request.Page,
		Total:    total,
		Offset:   request.Offset(),
		LastPage: lastPage,
	}
}
