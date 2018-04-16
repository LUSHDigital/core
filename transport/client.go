package microservicetransport

import (
	"net/http"
	"time"
)

// DefaultHTTPClient - returns a default http.Client implementation
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}
