# Pagination

The pagination package works with a combination of two types:

the `pagination.Request` and `pagination.Response` types.

The `Request` allows a client to pass in parameters to the server:

- `page`: which page the client wants to see
- `per_page`: how many items the client wants on each page.

From this request, a `Response` can be constructed, once you know how many total
items you are going to respond with. When making a new pagination response,
offset and last page values are calculated automatically.

## Examples

### Create a paginator
```go
func ExamplePagination() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    1,
	}
	presp := pagination.MakeResponse(preq, 100)
	fmt.Printf("%+v\n", presp)
	// Output: {PerPage:10 Page:1 Offset:0 Total:100 LastPage:10}
}
```

### Create a paginator, with an offset
```go
func ExampleOffsetPagination() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    2,
	}
	presp := pagination.MakeResponse(preq, 100)
	fmt.Printf("%+v\n", presp)
	// Output: {PerPage:10 Page:2 Offset:10 Total:100 LastPage:10}
}
```

### Create a paginator and use it in an HTTP response
```go
func ExamplePaginationResponse() {
	preq := pagination.Request{
		PerPage: 10,
		Page:    2,
	}
	presp := pagination.MakeResponse(preq, 100)

	resp := rest.Response{
		Code:    http.StatusOK,
		Message: "some helpful message",
		Data: &rest.Data{
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
```