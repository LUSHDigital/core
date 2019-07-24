package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LUSHDigital/core/pagination"
)

// OKResponse returns a prepared 200 OK response.
func OKResponse(data *Data, page *pagination.Response) *Response {
	return &Response{
		Code:       http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Data:       data,
		Pagination: page,
	}
}

// CreatedResponse returns a prepared 201 Created response.
func CreatedResponse(data *Data, page *pagination.Response) *Response {
	return &Response{
		Code:       http.StatusCreated,
		Message:    http.StatusText(http.StatusOK),
		Data:       data,
		Pagination: page,
	}
}

// Errorf returns a prepared error response using the provided code and formatted message.
func Errorf(code int, format string, a ...interface{}) *Response {
	return &Response{Code: code, Message: fmt.Sprintf(format, a...)}
}

// NoContentResponse returns a prepared 204 No Content response.
func NoContentResponse() *EmptyResponse {
	return &EmptyResponse{}
}

// JSONError returns a prepared 422 Unprocessable Entity response if the JSON is found to contain syntax errors, or invalid values for types.
func JSONError(msg interface{}) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("json error: %v", msg)}
}

// ParameterError returns a prepared 422 Unprocessable Entity response, including the name of the failing parameter in the message field of the response object.
func ParameterError(parameter string) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("invalid or missing parameter: %v", parameter)}
}

// ValidationError returns a prepared 422 Unprocessable Entity response, including the name of the failing validation/validator in the message field of the response object.
func ValidationError(resource string, msg interface{}) *Response {
	return &Response{Code: http.StatusUnprocessableEntity, Message: fmt.Sprintf("validation error on %s: %v", resource, msg)}
}

// NotFoundError returns a prepared 404 Not Found response, including the message passed by the user in the message field of the response object.
func NotFoundError(msg interface{}) *Response {
	return &Response{Code: http.StatusNotFound, Message: fmt.Sprintf("resource not found: %v", msg)}
}

// ConflictError returns a prepared 409 Conflict response, including the message passed by the user in the message field of the response object.
func ConflictError(msg interface{}) *Response {
	return &Response{Code: http.StatusConflict, Message: fmt.Sprintf("resource conflict: %v", msg)}
}

// InternalError returns a prepared 500 Internal Server Error, including the error message in the message field of the response object.
func InternalError(msg interface{}) *Response {
	return &Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("internal server error: %v", msg)}
}

// UnauthorizedError returns a prepared 401 Unauthorized error.
func UnauthorizedError() *Response {
	return &Response{Code: http.StatusUnauthorized, Message: "unauthorized"}
}

// !!! DEPRECATED FUNCTIONS !!!

// NotFoundErr returns a prepared 404 Not Found response, including the message passed by the user in the message field of the response object.
// DEPRECATED: Use NotFoundError rather than NotFoundErr
// TODO: Remove in version 1.x
func NotFoundErr(msg interface{}) *Response {
	log.Println("DEPRECATED: Use NotFoundError rather than NotFoundErr")
	return NotFoundError(msg)
}

// ConflictErr returns a prepared 409 Conflict response, including the message passed by the user in the message field of the response object.
// DEPRECATED: Use ConflicError rather than ConflictErr
// TODO: Remove in version 1.x
func ConflictErr(msg interface{}) *Response {
	log.Println("DEPRECATED: Use ConflicError rather than ConflictErr")
	return ConflictError(msg)
}

// ParamError returns a prepared 422 Unprocessable Entity response, including the name of the failing parameter in the message field of the response object.
// DEPRECATED: Use ParameterError rather than ParamError
// TODO: Remove in version 1.x
func ParamError(parameter string) *Response {
	log.Println("DEPRECATED: Use ParameterError rather than ParamError")
	return ParameterError(parameter)
}

// Unauthorized returns a prepared 401 Unauthorized error.
// DEPRECATED: Use UnauthorizedError rather than Unauthorized
// TODO: Remove in version 1.x
func Unauthorized() *Response {
	log.Println("DEPRECATED: Use UnauthorizedError rather than Unauthorized")
	return UnauthorizedError()
}

// DBError returns a prepared 500 Internal Server Error response.
// DEPRECATED: Use InternalError rather than DBError
// TODO: Remove in version 1.x
func DBError(msg interface{}) *Response {
	log.Println("DEPRECATED: Use InternalError rather than DBError")
	return &Response{Code: http.StatusInternalServerError, Message: fmt.Sprintf("db error: %v", msg)}
}

// DBErrorf returns a prepared 500 Internal Server Error response, using the user provided formatted message.
// DEPRECATED: Use InternalError rather than DBErrorf
// TODO: Remove in version 1.x
func DBErrorf(format string, err error) *Response {
	log.Println("DEPRECATED: Use InternalError rather than DBErrorf")
	var msg string
	switch format {
	case "":
		msg = fmt.Sprintf("db error: %v", err)
	default:
		msg = fmt.Sprintf(format, err)
	}
	return &Response{Code: http.StatusInternalServerError, Message: msg}
}
