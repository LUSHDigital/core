// Package transport defines a set of interface and methods for communicating
// with services and aggregators.
package transport

import "net/http"

// Transport - Interface responsible for communication.
type Transport interface {
	// Call - Do the current service request.
	Call() (*http.Response, error)

	// Dial - Create a request to a service resource.
	Dial(request *Request) error

	// GetName - Get the name of the service
	GetName() string
}
