package pagination_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/response"
)

// TestMakeResponse tests that the response is being tested
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
			res := pagination.MakeResponse(pagination.Request{PerPage: tc.perPage, Page: tc.page}, tc.total)
			if res.Offset != tc.expectedOffset {
				t.Fatalf("offset: want: %v\ngot: %v", tc.expectedOffset, res.Offset)
			}

			if res.LastPage != tc.expectedLastPage {
				t.Fatalf("last page: want: %v\ngot: %v", tc.expectedLastPage, res.LastPage)
			}
		})
	}
}

func ExampleMakeResponse() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    1,
	}
	presp := pagination.MakeResponse(preq, 100)
	fmt.Printf("%+v\n", presp)
	// Output: {PerPage:10 Page:1 Offset:0 Total:100 LastPage:10}
}

func ExampleMakeResponse_withOffset() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    2,
	}
	presp := pagination.MakeResponse(preq, 100)
	fmt.Printf("%+v\n", presp)
	// Output: {PerPage:10 Page:2 Offset:10 Total:100 LastPage:10}
}

func ExampleMakeResponse_withinResponse() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    2,
	}
	presp := pagination.MakeResponse(preq, 100)

	resp := response.Response{
		Code:    http.StatusOK,
		Message: "some helpful message",
		Data: &response.Data{
			Type:    "some_data",
			Content: map[string]interface{}{"hello": "world"},
		},
		Pagination: &presp,
	}
	raw, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println(string(raw))
	// Output:
	// {
	// 	"code": 200,
	// 	"message": "some helpful message",
	// 	"data": {
	// 		"some_data": {
	// 			"hello": "world"
	// 		}
	// 	},
	// 	"pagination": {
	// 		"per_page": 10,
	// 		"page": 2,
	// 		"offset": 10,
	// 		"total": 100,
	// 		"last_page": 10
	// 	}
	// }
}
