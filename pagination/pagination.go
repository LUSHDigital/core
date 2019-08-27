package pagination

import (
	"fmt"
	"math"

	"google.golang.org/grpc/metadata"
)

// Request derives the requested pagination data given by the client.
type Request struct {
	PerPage uint64
	Page    uint64
}

// Metadata returns gRPC metadata for a pagination request.
func (r Request) Metadata() metadata.MD {
	return metadata.New(map[string]string{
		"per_page": fmt.Sprintf("%d", r.PerPage),
		"page":     fmt.Sprintf("%d", r.Page),
	})
}

// Offset calculates the offset from the provided pagination request.
func (r Request) Offset() uint64 {
	return (r.Page - 1) * r.PerPage
}

// Response manages pagination of a data set.
type Response struct {
	PerPage     uint64  `json:"per_page"`     // The number of items per page.
	Offset      uint64  `json:"offset"`       // The current offset to pass to the query.
	Total       uint64  `json:"total"`        // The total number of items
	LastPage    uint64  `json:"last_page"`    // The number of the last possible page.
	CurrentPage uint64  `json:"current_page"` // The current page number.
	NextPage    *uint64 `json:"next_page"`    // The number of the next page (if possible).
	PrevPage    *uint64 `json:"prev_page"`    // The number of the previous page (if possible).
}

// Metadata returns gRPC metadata for a pagination response.
func (r Response) Metadata() metadata.MD {
	md := metadata.New(map[string]string{
		"per_page":     fmt.Sprintf("%d", r.PerPage),
		"offset":       fmt.Sprintf("%d", r.Offset),
		"total":        fmt.Sprintf("%d", r.Total),
		"last_page":    fmt.Sprintf("%d", r.LastPage),
		"current_page": fmt.Sprintf("%d", r.CurrentPage),
	})
	if r.NextPage != nil {
		md.Set("next_page", fmt.Sprintf("%d", *r.NextPage))
	}
	if r.PrevPage != nil {
		md.Set("prev_page", fmt.Sprintf("%d", *r.PrevPage))
	}
	return md
}

// MakeResponse returns a new Response with the provided
// parameters set.
func MakeResponse(request Request, total uint64) Response {
	// Avoid division by zero error, which would wrap uint to it's max value.
	var lastPage uint64
	if request.PerPage > 0 {
		lastPage = uint64(math.Ceil(float64(total) / float64(request.PerPage)))
	}

	resp := Response{
		PerPage:     request.PerPage,
		CurrentPage: request.Page,
		Total:       total,
		Offset:      request.Offset(),
		LastPage:    lastPage,
	}

	if lastPage > request.Page {
		next := request.Page + 1
		resp.NextPage = &next
	}
	if request.Page > 1 {
		prev := request.Page - 1
		resp.PrevPage = &prev
	}

	return resp
}
