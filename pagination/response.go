package pagination

// Response represents a pagination response.
type Response struct {
	Total       int  `json:"total"`        // The total number of items.
	PerPage     int  `json:"per_page"`     //  Number of items displayed per page.
	CurrentPage int  `json:"current_page"` // The current page number.
	LastPage    int  `json:"last_page"`    // The number of the last possible page.
	NextPage    *int `json:"next_page"`    // The number of the next page (if possible).
	PrevPage    *int `json:"prev_page"`    // The number of the previous page (if possible).
}

// newResponse returns a new paginated Response for a microservice endpoint.
func newResponse(total, perPage, currentPage, lastPage int) *Response {
	r := &Response{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: currentPage,
		LastPage:    lastPage,
	}

	// Set the next page.
	if r.LastPage > r.CurrentPage {
		nextPage := r.CurrentPage + 1
		r.NextPage = &nextPage
	}

	// Set the previous page.
	if r.CurrentPage > 1 {
		prevPage := r.CurrentPage - 1
		r.PrevPage = &prevPage
	}

	return r
}
