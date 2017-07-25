package format

import (
	"encoding/json"
	"net/http"

	"github.com/LUSHDigital/microservice-core-golang/response"
)

// JSONResponseFormatter - Format a microservice response as JSON.
//
// Params:
//     w http.ResponseWriter - The HTTP response writer.
//     response *Response - The microservice response object.
func JSONResponseFormatter(w http.ResponseWriter, response *response.Response) {
	// Set the content type header.
	w.Header().Set("Content-Type", "application/json")

	// Set the status code.
	w.WriteHeader(response.Code)

	// Set the response.
	json.NewEncoder(w).Encode(response)
}
