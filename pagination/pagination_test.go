package pagination

import (
	"reflect"
	"testing"
)

// TestNewPaginator - Test basic creation of the paginator.
func TestNewPaginator(t *testing.T) {
	tt := []struct {
		name             string
		perPage          int
		page             int
		total            int
		expectedOffset   int
		expectedLastPage int
	}{
		{
			name:             "100 items. 10 per page. Page 1.",
			perPage:          10,
			page:             1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 10,
		},
		{
			name:             "100 items. 10 per page. Page 2.",
			perPage:          10,
			page:             2,
			total:            100,
			expectedOffset:   10,
			expectedLastPage: 10,
		},
		{
			name:             "100 items. 7 per page. Page 1.",
			perPage:          7,
			page:             1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 15,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.perPage, tc.page, tc.total)
			if err != nil {
				t.Errorf("failed to create paginator: %s", err)
			}

			if paginator.GetOffset() != tc.expectedOffset {
				t.Errorf("offset: want: %v\ngot: %v", tc.expectedOffset, paginator.GetOffset())
			}

			if paginator.GetLastPage() != tc.expectedLastPage {
				t.Errorf("last page: want: %v\ngot: %v", tc.expectedLastPage, paginator.GetLastPage())
			}
		})
	}
}

// TestNewPaginator - Test changing the page of a paginator.
func TestPaginator_SetPage(t *testing.T) {
	tt := []struct {
		name             string
		perPage          int
		page             int
		changePage       int
		total            int
		expectedOffset   int
		expectedLastPage int
	}{
		{
			name:             "100 items. 10 per page. Page 1.",
			perPage:          10,
			page:             1,
			changePage:       1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 10,
		},
		{
			name:             "100 items. 10 per page. Page 2.",
			perPage:          10,
			page:             1,
			changePage:       2,
			total:            100,
			expectedOffset:   10,
			expectedLastPage: 10,
		},
		{
			name:             "100 items. 10 per page. Page 3.",
			perPage:          10,
			page:             1,
			changePage:       3,
			total:            100,
			expectedOffset:   20,
			expectedLastPage: 10,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.perPage, tc.page, tc.total)
			if err != nil {
				t.Errorf("failed to create paginator: %s", err)
			}

			paginator.SetPage(tc.changePage)

			if paginator.GetOffset() != tc.expectedOffset {
				t.Errorf("%s: offset: want: %v\ngot: %v", tc.name, tc.expectedOffset, paginator.GetOffset())
			}

			if paginator.GetLastPage() != tc.expectedLastPage {
				t.Errorf("%s: last page: want: %v\ngot: %v", tc.name, tc.expectedLastPage, paginator.GetLastPage())
			}
		})
	}
}

// TestPaginator_SetPerPage - Test changing the page of a paginator.
func TestPaginator_SetPerPage(t *testing.T) {
	tt := []struct {
		name             string
		perPage          int
		changePerPage    int
		page             int
		total            int
		expectedOffset   int
		expectedLastPage int
	}{
		{
			name:             "100 items. 20 per page. Page 1.",
			perPage:          20,
			changePerPage:    20,
			page:             1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 5,
		},
		{
			name:             "100 items. 30 per page. Page 1.",
			perPage:          30,
			changePerPage:    30,
			page:             1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 4,
		},
		{
			name:             "100 items. 40 per page. Page 1.",
			perPage:          40,
			changePerPage:    40,
			page:             1,
			total:            100,
			expectedOffset:   0,
			expectedLastPage: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.perPage, tc.page, tc.total)
			if err != nil {
				t.Errorf("failed to create paginator: %s", err)
			}

			paginator.SetPerPage(tc.changePerPage)

			if paginator.GetOffset() != tc.expectedOffset {
				t.Errorf("%s: offset: want: %v\ngot: %v", tc.name, tc.expectedOffset, paginator.GetOffset())
			}

			if paginator.GetLastPage() != tc.expectedLastPage {
				t.Errorf("%s: last page: want: %v\ngot: %v", tc.name, tc.expectedLastPage, paginator.GetLastPage())
			}
		})
	}
}

// TestPaginator_SetTotal - Test changing the page of a paginator.
func TestPaginator_SetTotal(t *testing.T) {
	tt := []struct {
		name             string
		perPage          int
		page             int
		total            int
		changeTotal      int
		expectedOffset   int
		expectedLastPage int
	}{
		{
			name:             "100 items. 10 per page. Page 1.",
			perPage:          10,
			page:             1,
			total:            100,
			changeTotal:      100,
			expectedOffset:   0,
			expectedLastPage: 10,
		},
		{
			name:             "20 items. 10 per page. Page 1.",
			perPage:          10,
			page:             1,
			total:            100,
			changeTotal:      20,
			expectedOffset:   0,
			expectedLastPage: 2,
		},
		{
			name:             "8 items. 10 per page. Page 1.",
			perPage:          10,
			page:             1,
			total:            100,
			changeTotal:      8,
			expectedOffset:   0,
			expectedLastPage: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.perPage, tc.page, tc.total)
			if err != nil {
				t.Errorf("failed to create paginator: %s", err)
			}

			paginator.SetTotal(tc.changeTotal)

			if paginator.GetOffset() != tc.expectedOffset {
				t.Errorf("%s: offset: want: %v\ngot: %v", tc.name, tc.expectedOffset, paginator.GetOffset())
			}

			if paginator.GetLastPage() != tc.expectedLastPage {
				t.Errorf("%s: last page: want: %v\ngot: %v", tc.name, tc.expectedLastPage, paginator.GetLastPage())
			}
		})
	}
}

// TestPaginator_PrepareResponse - Test the paginator response.
func TestPaginator_PrepareResponse(t *testing.T) {
	tt := []struct {
		name     string
		perPage  int
		page     int
		total    int
		response Response
	}{
		{
			name:    "100 items. 10 per page. Page 1.",
			perPage: 10,
			page:    1,
			total:   100,
			response: Response{
				Total:       100,
				PerPage:     10,
				CurrentPage: 1,
				LastPage:    10,
				NextPage:    func(i int) *int { return &i }(2),
				PrevPage:    nil,
			},
		},
		{
			name:    "100 items. 10 per page. Page 2.",
			perPage: 10,
			page:    2,
			total:   100,
			response: Response{
				Total:       100,
				PerPage:     10,
				CurrentPage: 2,
				LastPage:    10,
				NextPage:    func(i int) *int { return &i }(3),
				PrevPage:    func(i int) *int { return &i }(1),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.perPage, tc.page, tc.total)
			if err != nil {
				t.Errorf("failed to create paginator: %s", err)
			}

			if !reflect.DeepEqual(paginator.PrepareResponse(), &tc.response) {
				t.Errorf("TestPaginator_PrepareResponse: %s: expected %v got %v", tc.name, &tc.response, paginator.PrepareResponse())
			}
		})
	}
}
