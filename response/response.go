package response

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Standard response statuses.
const (
	StatusOk   = "ok"
	StatusFail = "fail"
)

// Response - A standardised response format for a microservice.
type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *Data  `json:"data,omitempty"`
}

// Data represents the collection data the the response will return to the consumer
// Type ends up being the name of the key containing the collection of Content
type Data struct {
	Type    string
	Content interface{}
}

// UnmarshalJSON implements the Unmarshaler interface
// this implementation will fill the type in the case we're been provided a valid single collection
// and set the content to the contents of said collection.
// for every other options, it behaves like normal.
// Despite the fact that we are not suposed to marshal without a type set
// This is purposefuly left open to unmarshal without a collection name set, in case you may want to set it later,
// and for interop with other systems which may not send the collection properly.
func (d *Data) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &d.Content); err != nil {
		log.Printf("cannot unmarshal data: %v", err)
	}

	data, ok := d.Content.(map[string]interface{})
	if ok {
		// count how many collections were provided
		var count int
		for _, value := range data {
			if _, ok := value.(map[string]interface{}); ok {
				count++
			}
		}
		if count > 1 {
			// we can stop there since this is not a single collection
			return nil
		}

		for key, value := range data {
			if _, ok := value.(map[string]interface{}); ok {
				d.Type = key
				d.Content = data[key]
			}
		}
	}

	return nil
}

// Valid ensures the Data passed to the response is correct
func (d *Data) Valid() bool {
	if d.Type != "" {
		return true
	}
	return false
}

// MarshalJSON implements the Marshaler interface and is there to ensure the output
// is correct when we return data to the consumer
func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Map())
}

// Map returns a version of the data as a map
func (d *Data) Map() map[string]interface{} {
	d.Type = strings.Replace(strings.ToLower(d.Type), " ", "-", -1)
	if !d.Valid() {
		log.Printf("invalid data: %v", d)
		return nil
	}

	return map[string]interface{}{
		d.Type: d.Content,
	}
}

// New returns a new Response for a microservice endpoint
// This ensures that all API endpoints return data in a standardised format:
//
//    {
//       "status": "ok", - Can contain any string. Usually 'ok', 'error' etc.
//       "code": 200, - A HTTP status code.
//       "message": "", - A message string elaborating on the status.
//       "data": {[ - A collection of return data. Can be omitted in the event an error occurred.
//       ]}
//    }
// Params:
//   - [code] - HTTP status code for the response.
//   - [status] - A short status message. Examples: 'OK', 'Bad Request', 'Not Found' etc...
//   - [message] - A more detailed status message
//   - [data] The data to return. Will always be parsed into a collection.
//
// Return:
//   *Response - The populated response object.
func New(code int, status, message string, data *Data) *Response {
	return &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ExtractData - Extract a particular item of data from the response.
//
// Params:
//     srcKey string - The name of the data item we want from the response.
//     dst interface{} - The interface to extract data into.
//
// Return:
//     error - An error if it occurred.
func (r *Response) ExtractData(srcKey string, dst interface{}) error {
	if !r.Data.Valid() {
		return fmt.Errorf("invalid data provided: %v", r.Data)
	}

	for key, value := range r.Data.Map() {
		if key != srcKey {
			continue
		}

		// Get the raw JSON just for the endpoints.
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}

		// Decode the raw JSON.
		json.Unmarshal(rawJSON, &dst)
	}

	return nil
}
