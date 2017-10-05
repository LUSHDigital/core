// Package format defines any necessary custom formatters for the responses
// presently only JSON is supported.
package format

import (
	"encoding/json"
	"net/http"

	"github.com/LUSHDigital/microservice-core-golang/response"
)

// JSONResponseFormatter formats a microservice response as JSON.
func JSONResponseFormatter(w http.ResponseWriter, response response.ResponseInterface) {
	// Set the content type header.
	w.Header().Set("Content-Type", "application/json")

	// Set the status code.
	w.WriteHeader(response.GetCode())

	// Set the response.
	json.NewEncoder(w).Encode(response)
}
