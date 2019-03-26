package pagination_test

import (
	"testing"

	"github.com/LUSHDigital/core/pagination"
)

func TestMakeResponse(t *testing.T) {
	tt := []struct {
		name             string
		perPage          uint64
		page             uint64
		total            uint64
		expectedOffset   uint64
		expectedLastPage uint64
		expErr           error
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
		{
			name:             "0 items",
			perPage:          5,
			page:             1,
			total:            0,
			expectedOffset:   0,
			expectedLastPage: 0,
		},
		{
			name:    "100 items. 0 per page. Page 1",
			perPage: 0,
			page:    1,
			total:   100,
			expErr:  pagination.ErrCalculateOffset,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			paginator := pagination.MakeResponse(pagination.Request{PerPage: tc.perPage, Page: tc.page}, tc.total)
			if paginator.Offset != tc.expectedOffset {
				t.Fatalf("offset: want: %v\ngot: %v", tc.expectedOffset, paginator.Offset)
			}

			if paginator.LastPage != tc.expectedLastPage {
				t.Fatalf("last page: want: %v\ngot: %v", tc.expectedLastPage, paginator.LastPage)
			}
		})
	}
}

