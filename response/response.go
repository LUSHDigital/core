// Package response defines the how the default microservice response must look and behave like.
package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/LUSHDigital/core/pagination"
)

// Responder defines the behaviour of a response for JSON over HTTP.
type Responder interface {
	WriteTo(w http.ResponseWriter) error
}

// Response defines a JSON response body over HTTP.
type Response struct {
	Code       int                  `json:"code"`                 // Any valid HTTP response code
	Message    string               `json:"message"`              // Any relevant message (optional)
	Data       *Data                `json:"data,omitempty"`       // Data to pass along to the response (optional)
	Pagination *pagination.Response `json:"pagination,omitempty"` // Pagination data
}

// WriteTo writes a JSON response to a HTTP writer.
func (r *Response) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	// Don't attempt to write a body for 204s.
	if r.Code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(r)
}

// Data represents the collection data the the response will return to the consumer.
// Type ends up being the name of the key containing the collection of Content
type Data struct {
	Type    string
	Content interface{}
}

// UnmarshalJSON implements the Unmarshaler interface
// this implementation will fill the type in the case we're been provided a valid single collection
// and set the content to the contents of said collection.
// for every other options, it behaves like normal.
// Despite the fact that we are not supposed to marshal without a type set,
// this is purposefully left open to unmarshal without a collection name set, in case you may want to set it later,
// and for interop with other systems which may not send the collection properly.
func (d *Data) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &d.Content); err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}

	data, ok := d.Content.(map[string]interface{})
	if !ok {
		return nil
	}
	// count how many collections were provided
	var count int
	for _, value := range data {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			count++
		}
	}
	if count > 1 {
		// we can stop there since this is not a single collection
		return nil
	}
	for key, value := range data {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			d.Type = key
			d.Content = data[key]
		}
	}

	return nil
}

var invalidDataError = errors.New("invalid data provided")
// MarshalJSON implements the Marshaler interface and is there to ensure the output
// is correct when we return data to the consumer
func (d *Data) MarshalJSON() ([]byte, error) {
	if d.Type == "" || strings.Contains(d.Type, " ") {
		return nil, invalidDataError
	}

	m := map[string]interface{}{
		d.Type: d.Content,
	}
	return json.Marshal(m)
}
