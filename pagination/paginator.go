package pagination

import (
	"errors"
	"math"
)

// Paginator - Manages pagination of a data set.
type Paginator struct {
	perPage  int // The number of items per page.
	page     int // Which page are we on?
	offset   int // The current offset to pass to the query.
	total    int // The total number of items
	lastPage int // The number of the last possible page.
}

// calculateOffset - Calculate the offset based on the current values.
//
// Return:
//     error - An error if it occurred.
func (p *Paginator) calculateOffset() error {
	if p.page == 0 || p.perPage == 0 {
		return errors.New("cannot calculate offset: insufficient data")
	}

	p.offset = (p.page - 1) * p.perPage
	return nil
}

// calculateLastPage - Calculate the last page based on the current values.
//
// Return:
//     error - An error if it occurred.
func (p *Paginator) calculateLastPage() error {
	if p.total == 0 || p.perPage == 0 {
		return errors.New("cannot calculate last page: insufficient data")
	}

	p.lastPage = int(math.Ceil(float64(p.total) / float64(p.perPage)))

	return nil
}

// GetPerPage
//
// Return:
//     int - Number of items per page.
func (p *Paginator) GetPerPage() int {
	return p.perPage
}

// SetPerPage
//
// Params:
//     perPage int - Number of items per page.
//
// Return:
//     error - An error if it occurred.
func (p *Paginator) SetPerPage(perPage int) error {
	p.perPage = perPage

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
	}

	lastPageErr := p.calculateLastPage()
	if lastPageErr != nil {
		return lastPageErr
	}
	return nil
}

// GetPage
//
// Return:
//     int - Which page are we on?
func (p *Paginator) GetPage() int {
	return p.page
}

// SetPage
//
// Params:
//     page int - Which page are we on?
//
// Return:
//     error - An error if it occurred.
func (p *Paginator) SetPage(page int) error {
	p.page = page

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
	}

	lastPageErr := p.calculateLastPage()
	if lastPageErr != nil {
		return lastPageErr
	}
	return nil
}

// GetOffset
//
// Return:
//     int - The current offset to pass to the query.
func (p *Paginator) GetOffset() int {
	return p.offset
}

// GetTotal
//
// Return:
//     int - The total number of items.
func (p *Paginator) GetTotal() int {
	return p.total
}

// GetTotal
//
// Params:
//     total int - The total number of items.
//
// Return:
//     error - An error if it occurred.
func (p *Paginator) SetTotal(total int) error {
	p.total = total

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
	}

	lastPageErr := p.calculateLastPage()
	if lastPageErr != nil {
		return lastPageErr
	}
	return nil
}

// GetLastPage
//
// Return:
//     int - The number of the last possible page
func (p *Paginator) GetLastPage() int {
	return p.lastPage
}

// PrepareResponse - Prepare the pagination response.
//
// Return:
//     *response - The pagination response.
func (p *Paginator) PrepareResponse() *Response {
	return newResponse(p.total, p.perPage, p.page, p.lastPage)
}

// NewPaginator - Instantiate a new paginator.
//
// Params:
//     perPage int - Number of items to display per page.
//     page int - Which page are we on?
//     total int - Number of items there are in total.
//
// Return:
//     *Paginator - The instantiated paginator object.
//     error - An error if it occurred
func NewPaginator(perPage, page, total int) (*Paginator, error) {
	// Create the paginator.
	p := Paginator{
		perPage: perPage,
		page:    page,
		total:   total,
	}

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return nil, offsetErr
	}

	lastPageErr := p.calculateLastPage()
	if lastPageErr != nil {
		return nil, lastPageErr
	}

	return &p, nil
}
