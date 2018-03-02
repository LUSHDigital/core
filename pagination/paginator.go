// Package pagination defines a paginator able to return formatted responses
// enabling the API consumer to retrieve data in defined chunks
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

// calculateOffset sets the offset field based on the current values.
func (p *Paginator) calculateOffset() error {
	if p.page == 0 || p.perPage == 0 {
		return errors.New("cannot calculate offset: insufficient data")
	}

	p.offset = (p.page - 1) * p.perPage
	return nil
}

// calculateLastPage sets the lastPage field based on the
// current values and returns an error if it fails.
func (p *Paginator) calculateLastPage() error {
	if p.total == 0 || p.perPage == 0 {
		return errors.New("cannot calculate last page: insufficient data")
	}

	p.lastPage = int(math.Ceil(float64(p.total) / float64(p.perPage)))

	return nil
}

// GetPerPage returns the number of items per page.
func (p *Paginator) GetPerPage() int {
	return p.perPage
}

// SetPerPage defines how many items per page the
// paginator will return, based on the supplied parameter,
// and will return an error if anything fails.
func (p *Paginator) SetPerPage(perPage int) error {
	p.perPage = perPage

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
	}

	return p.calculateLastPage()
}

// GetPage returns the current page index
func (p *Paginator) GetPage() int {
	return p.page
}

// SetPage sets the page field to the provided value
// and returns an error if anything fails
func (p *Paginator) SetPage(page int) error {
	p.page = page

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
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

// SetTotal sets the total number of items in the paginator
// to the provided value and returns an error if it fails.
func (p *Paginator) SetTotal(total int) error {
	p.total = total

	offsetErr := p.calculateOffset()
	if offsetErr != nil {
		return offsetErr
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
// parameters set and reutrns an error if it fails.
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
