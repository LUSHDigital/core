package rest

import (
	"fmt"
	"net/http"
)

// DBError returns a prepared 500 Internal Server Error response.
func DBError(err error) *Response {
	return DBErrorf("", err)
}

// DBErrorf returns a prepared 500 Internal Server Error response,
// using the user provided formatted message.
func DBErrorf(format string, err error) *Response {
	var msg string
	switch format {
	case "":
		msg = fmt.Sprintf("db error: %v", err)
	default:
		msg = fmt.Sprintf(format, err)
	}
	return &Response{Code: http.StatusInternalServerError, Message: msg}
}

// JSONError returns a prepared 422 Unprocessable Entity response if the JSON is found to
// contain syntax errors, or invalid values for types.
func JSONError(err error) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("json error: %v", err)}
}

// ParamError returns a prepared 422 Unprocessable Entity response, including the name of
// the failing parameter in the message field of the response object.
func ParamError(name string) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("invalid or missing parameter: %v", name)}
}

// ValidationError returns a prepared 422 Unprocessable Entity response, including the name of
// the failing validation/validator in the message field of the response object.
func ValidationError(err error, name string) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("validation error on %s: %v", name, err)}
}

// NotFoundErr returns a prepared 404 Not Found response, including the message passed by the user
// in the message field of the response object.
func NotFoundErr(msg string) *Response {
	return &Response{Code: http.StatusNotFound, Message: msg}
}

// ConflictErr returns a prepared 409 Conflict response, including the message passed by the user
// in the message field of the response object.
func ConflictErr(msg string) *Response {
	return &Response{Code: http.StatusConflict, Message: msg}
}

// InternalError returns a prepared 500 Internal Server Error, including the error
// message in the message field of the response object.
func InternalError(err error) *Response {
	return &Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("internal server error: %v", err)}
}

// Unauthorized returns a prepared 401 Unauthorized error.
func Unauthorized() *Response {
	return &Response{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}
}
