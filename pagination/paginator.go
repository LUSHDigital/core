// Package pagination defines a paginator able to return formatted responses
// enabling the API consumer to retrieve data in defined chunks
package pagination

import (
	"math"
)

// Paginator manages pagination of a data set.
type Paginator struct {
	perPage  int // The number of items per page.
	page     int // Which page are we on?
	offset   int // The current offset to pass to the query.
	total    int // The total number of items
	lastPage int // The number of the last possible page.
}

// calculateOffset sets the offset field based on the current values.
func (p *Paginator) calculateOffset() error {
	if p.page == 0 || p.perPage == 0 {
		return ErrCalculateOffset
	}

	p.offset = (p.page - 1) * p.perPage
	return nil
}

// calculateLastPage sets the lastPage field based on the current values and
// returns an error if it fails.
func (p *Paginator) calculateLastPage() error {
	// If there are no items to paginate return early to prevent divide-by-zero.
	// An error is too heavy-handed here and makes the library difficult to use
	// when total >= 0.
	if p.total == 0 {
		return nil
	}

	// A per-page value of zero on the other hand is a bit crazy, so an error is
	// an acceptable response.
	if p.perPage == 0 {
		return ErrCalculateLastPage
	}

	p.lastPage = int(math.Ceil(float64(p.total) / float64(p.perPage)))

	return nil
}

// GetPerPage returns the number of items per page.
func (p *Paginator) GetPerPage() int {
	return p.perPage
}

// SetPerPage defines how many items per page the paginator will return, based
// on the supplied parameter, and will return an error if anything fails.
func (p *Paginator) SetPerPage(perPage int) error {
	p.perPage = perPage

	if err := p.calculateOffset(); err != nil {
		return err
	}

	return p.calculateLastPage()
}

// GetPage returns the current page index
func (p *Paginator) GetPage() int {
	return p.page
}

// SetPage sets the page field to the provided value and returns an error if
// anything fails
func (p *Paginator) SetPage(page int) error {
	p.page = page

	if err := p.calculateOffset(); err != nil {
		return err
	}

	return p.calculateLastPage()
}

// GetOffset returns the current offset of the paginator.
func (p *Paginator) GetOffset() int {
	return p.offset
}

// GetTotal returns the total number of items in the paginator.
func (p *Paginator) GetTotal() int {
	return p.total
}

// SetTotal sets the total number of items in the paginator to the provided
// value and returns an error if it fails.
func (p *Paginator) SetTotal(total int) error {
	p.total = total

	if err := p.calculateOffset(); err != nil {
		return err
	}

	return p.calculateLastPage()
}

// GetLastPage returns the last possible page number.
func (p *Paginator) GetLastPage() int {
	return p.lastPage
}

// PrepareResponse returns a prepared pagination response.
func (p *Paginator) PrepareResponse() *Response {
	return newResponse(p.total, p.perPage, p.page, p.lastPage)
}

// NewPaginator returns a new Paginator instance with the provided
// parameters set and returns an error if it fails.
func NewPaginator(perPage, page, total int) (paginator *Paginator, err error) {
	// Create the paginator.
	paginator = &Paginator{
		perPage: perPage,
		page:    page,
		total:   total,
	}

	if err = paginator.calculateOffset(); err != nil {
		return nil, err
	}

	if err = paginator.calculateLastPage(); err != nil {
		return nil, err
	}
	return
}
