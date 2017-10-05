// Package routing defines the basic structure a microservice route must abide by.
package routing

import "net/http"

// Route defines an HTTP route
type Route struct {
	Path    string                                   `json:"uri"`
	Method  string                                   `json:"method"`
	Handler func(http.ResponseWriter, *http.Request) `json:"-"`
}
