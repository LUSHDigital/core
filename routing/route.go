package routing

import "net/http"

// Route - A HTTP route.
type Route struct {
	Path    string                                   `json:"uri"`
	Method  string                                   `json:"method"`
	Handler func(http.ResponseWriter, *http.Request) `json:"-"`
}
