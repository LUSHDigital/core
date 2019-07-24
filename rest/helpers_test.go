package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/rest"
	"github.com/LUSHDigital/core/test"
)

type TestPayload struct {
	Message string `json:"message"`
}

type Envelope struct {
	Message string `json:"message"`
}

var (
	p    = uint64(1)
	page = &pagination.Response{
		PerPage:     1,
		Offset:      1,
		Total:       1,
		LastPage:    1,
		CurrentPage: 1,
		NextPage:    &p,
		PrevPage:    &p,
	}
	payload = &rest.Data{
		Type: "test",
		Content: &TestPayload{
			Message: "hello",
		},
	}
)

func TestOKResponse(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.OKResponse(payload, page).WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusOK, res.StatusCode)
	tpl := &TestPayload{}
	err = rest.UnmarshalJSONResponse(req.Body.Bytes(), &tpl)
	test.Equals(t, nil, err)
	test.Equals(t, "hello", tpl.Message)
}

func TestCreatedResponse(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.CreatedResponse(payload, page).WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusCreated, res.StatusCode)
	tpl := &TestPayload{}
	err = rest.UnmarshalJSONResponse(req.Body.Bytes(), &tpl)
	test.Equals(t, nil, err)
	test.Equals(t, "hello", tpl.Message)
}

func TestNoContentError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.NoContentResponse().WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusNoContent, res.StatusCode)
	test.Equals(t, 0, req.Body.Len())
}

func TestJSONError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.JSONError("cannot marshal").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusUnprocessableEntity, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "json error: cannot marshal", e.Message)
}

func TestParameterError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.ParameterError("name").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusUnprocessableEntity, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "invalid or missing parameter: name", e.Message)
}

func TestValidationError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.ValidationError("user", "name cannot be blank").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusUnprocessableEntity, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "validation error on user: name cannot be blank", e.Message)
}

func TestNotFoundError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.NotFoundError("user does not exist").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusNotFound, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "resource not found: user does not exist", e.Message)
}

func TestConflictError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.ConflictError("user already exists").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusConflict, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "resource conflict: user already exists", e.Message)
}

func TestInternalError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.InternalError("cannot connect to database").WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusInternalServerError, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "internal server error: cannot connect to database", e.Message)
}

func TestUnauthorizedError(t *testing.T) {
	req := httptest.NewRecorder()
	err := rest.UnauthorizedError().WriteTo(req)
	test.Equals(t, nil, err)
	res := req.Result()
	test.Equals(t, http.StatusUnauthorized, res.StatusCode)
	e := &Envelope{}
	json.Unmarshal(req.Body.Bytes(), e)
	test.Equals(t, "unauthorized", e.Message)
}
