// Package response defines the how the default microservice response must look and behave like.
package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/LUSHDigital/core/pagination"
)

type data struct {
	Type    string `json:"type"`
	Content json.RawMessage
}

func (d *data) UnmarshalJSON(b []byte) error {
	var data = make(map[string]interface{})
	if err := json.Unmarshal(b, &data); err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}
	var count int
	for _, value := range data {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			count++
		}
	}
	if count > 1 {
		return nil
	}
	for key, value := range data {
		switch value.(type) {
		case map[string]interface{}, []interface{}:
			b, err := json.Marshal(data[key])
			if err != nil {
				return nil
			}
			d.Type = key
			d.Content = b
		}
	}

	return nil
}

// UnmarshalJSONResponse will unmarshal the data from legacy response envelope.
func UnmarshalJSONResponse(d []byte, dst interface{}) error {
	type envelope struct {
		Data *data `json:"data"`
	}
	var e envelope
	if err := json.Unmarshal(d, &e); err != nil {
		return err
	}
	if e.Data == nil {
		return nil
	}
	err := json.Unmarshal(e.Data.Content, dst)
	if err != nil {
		return err
	}
	return nil
}

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
func (r Response) WriteTo(w http.ResponseWriter) error {
	return WriteTo(r.Code, r, w)
}

// WriteTo writes any JSON response to a HTTP writer.
func WriteTo(code int, i interface{}, w http.ResponseWriter) error {
	w.WriteHeader(code)
	// Don't attempt to write a body for 204s.
	if code == http.StatusNoContent {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(i)
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

var errInvalidKeyName = func(key string) error { return fmt.Errorf("invalid key name: %q", key) }

// MarshalJSON implements the Marshaler interface and is there to ensure the output
// is correct when we return data to the consumer
func (d *Data) MarshalJSON() ([]byte, error) {
	if d.Type == "" || strings.Contains(d.Type, " ") {
		return nil, errInvalidKeyName(d.Type)
	}
	return json.Marshal(map[string]interface{}{
		d.Type: d.Content,
	})
}
