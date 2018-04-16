// Package response defines the how the default microservice response must look and behave like.
package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"database/sql"

	"github.com/LUSHDigital/microservice-core-golang/pagination"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

// Standard response statuses.
const (
	StatusOk   = "ok"
	StatusFail = "fail"
)

// Responder - Responder for microservice responses.
type Responder interface {
	// ExtractData returns a particular item of data from the response.
	ExtractData(srcKey string, dst interface{}) error

	// GetCode returns the response code.
	GetCode() int
}

// Response - A standardised response format for a microservice.
type Response struct {
	Status  string `json:"status"`         // Can be 'ok' or 'fail'
	Code    int    `json:"code"`           // Any valid HTTP response code
	Message string `json:"message"`        // Any relevant message (optional)
	Data    *Data  `json:"data,omitempty"` // Data to pass along to the response (optional)
}

// New returns a new Response for a microservice endpoint
// This ensures that all API endpoints return data in a standardised format:
//
//    {
//       "status": "ok or fail",
//       "code": any HTTP response code,
//       "message": "any relevant message (optional)",
//       "data": {[
//          ...
//       ]}
//    }
func New(code int, message string, data *Data) *Response {
	var status string
	switch {
	case code >= http.StatusOK && code < http.StatusBadRequest:
		status = StatusOk
	default:
		status = StatusFail
	}
	return &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// SQLError returns a prepared 204 No Content response if the error passed is of type sql.ErrNoRows,
// otherwise, returns a 500 Internal Server Error prepared response.
func SQLError(err error) *Response {
	return SQLErrorf("", err)
}

// SQLErrorf allows a custom error message to be passed to the SQLError function.
func SQLErrorf(format string, err error) *Response {
	if err == sql.ErrNoRows {
		return New(http.StatusNoContent, "no data found", nil)
	}
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
			return New(http.StatusUnprocessableEntity, "duplicate entry.", nil)
		}
	}

	// Use any format message provided by the user, otherwise, just return the error string.
	var message string
	if format == "" {
		message = fmt.Sprintf("db error: %v", err)
	} else {
		message = fmt.Sprintf(format, err)
	}

	return New(http.StatusInternalServerError, message, nil)
}

// JSONError returns a prepared 422 Unprocessable Entity response if the JSON is found to
// contain syntax errors, or invalid values for types.
func JSONError(err error) *Response {
	return New(http.StatusUnprocessableEntity, fmt.Sprintf("json error: %v", err), nil)
}

// ParamError returns a prepared 422 Unprocessable Entity response, including the name of
// the failing parameter in the message field of the response object.
func ParamError(name string) *Response {
	return New(http.StatusUnprocessableEntity, fmt.Sprintf("invalid or missing parameter: %v", name), nil)
}

// ValidationError returns a prepared 422 Unprocessable Entity response, including the name of
// the failing validation/validator in the message field of the response object.
func ValidationError(err error, name string) *Response {
	return New(http.StatusUnprocessableEntity, fmt.Sprintf("validation error on %s: %v", name, err), nil)
}

// NotFoundErr returns a prepared 404 Not Found response, including the message passed by the user
// in the message field of the response object.
func NotFoundErr(msg string) *Response {
	return New(http.StatusNotFound, msg, nil)
}

// ConflictErr returns a prepared 409 Conflict response, including the message passed by the user
// in the message field of the response object.
func ConflictErr(msg string) *Response {
	return New(http.StatusConflict, msg, nil)
}

// InternalError returns a prepared 500 Internal Server Error, including the error
// message in the message field of the response object
func InternalError(err error) *Response {
	return New(http.StatusInternalServerError, fmt.Sprintf("internal server error: %v", err), nil)
}

// WriteTo - pick a response writer to write the default json response to.
func (r *Response) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	j, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}

// ExtractData returns a particular item of data from the response.
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

// GetCode returns the response code.
func (r *Response) GetCode() int {
	return r.Code
}

// PaginatedResponse - A paginated response format for a microservice.
type PaginatedResponse struct {
	Status     string               `json:"status"`         // Can be 'ok' or 'fail'
	Code       int                  `json:"code"`           // Any valid HTTP response code
	Message    string               `json:"message"`        // Any relevant message (optional)
	Data       *Data                `json:"data,omitempty"` // Data to pass along to the response (optional)
	Pagination *pagination.Response `json:"pagination"`     // Pagination data
}

// NewPaginated returns a new PaginatedResponse for a microservice endpoint
func NewPaginated(paginator *pagination.Paginator, code int, message string, data *Data) *PaginatedResponse {
	var status string
	switch {
	case code >= http.StatusOK && code < http.StatusBadRequest:
		status = StatusOk
	default:
		status = StatusFail
	}
	return &PaginatedResponse{
		Code:       code,
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: paginator.PrepareResponse(),
	}
}

// WriteTo - pick a response writer to write the default json response to.
func (p *PaginatedResponse) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.Code)

	j, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}

// ExtractData returns a particular item of data from the response.
func (p *PaginatedResponse) ExtractData(srcKey string, dst interface{}) error {
	if !p.Data.Valid() {
		return fmt.Errorf("invalid data provided: %v", p.Data)
	}
	for key, value := range p.Data.Map() {
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

// GetCode returns the response code.
func (p *PaginatedResponse) GetCode() int {
	return p.Code
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

// Valid ensures the Data passed to the response is correct (it must contain a Type along with the data).
func (d *Data) Valid() bool {
	return d.Type != ""
}

// MarshalJSON implements the Marshaler interface and is there to ensure the output
// is correct when we return data to the consumer
func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Map())
}

// Map returns a version of the data as a map
func (d *Data) Map() map[string]interface{} {
	if !d.Valid() {
		return nil
	}
	d.Type = strings.Replace(strings.ToLower(d.Type), " ", "-", -1)

	return map[string]interface{}{
		d.Type: d.Content,
	}
}
