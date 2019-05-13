package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LUSHDigital/core/pagination"
	"github.com/LUSHDigital/core/test"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestData_MarshalJSON(t *testing.T) {
	preq := pagination.MakeResponse(pagination.Request{
		PerPage: 1,
		Page:    1,
	}, 1)
	cases := []struct {
		name     string
		response *Response
		expected []byte
		wantsErr bool
	}{
		{
			name: "valid response with pagination",
			response: &Response{
				Code:    200,
				Message: "",
				Data: &Data{
					Type:    "test",
					Content: map[string]interface{}{"test": "test"},
				},
				Pagination: &preq,
			},
			expected: []byte(`{"code":200,"message":"","data":{"test":{"test":"test"}},"pagination":{"per_page":1,"offset":0,"total":1,"last_page":1,"current_page":1,"next_page":null,"prev_page":null}}`),
			wantsErr: false,
		},
		{
			name: "valid response without pagination",
			response: &Response{
				Code:    200,
				Message: "",
				Data: &Data{
					Type:    "test",
					Content: map[string]interface{}{"test": "test"},
				},
			},
			expected: []byte(`{"code":200,"message":"","data":{"test":{"test":"test"}}}`),
			wantsErr: false,
		},
		{
			name: "valid response without data",
			response: &Response{
				Code:    200,
				Message: "",
			},
			expected: []byte(`{"code":200,"message":""}`),
			wantsErr: false,
		},
		{
			name: "invalid response with empty type",
			response: &Response{
				Code:    200,
				Message: "",
				Data: &Data{
					Type:    "",
					Content: map[string]interface{}{"test": "test"},
				},
			},
			expected: []byte(`{"code":200,"message":"","data":{"test":{"test":"test"}}}`),
			wantsErr: true,
		},
		{
			name: "invalid response with empty type",
			response: &Response{
				Code:    200,
				Message: "",
				Data: &Data{
					Type:    "ha ha",
					Content: map[string]interface{}{"test": "test"},
				},
			},
			expected: []byte(`{"code":200,"message":"","data":{"test":{"test":"test"}}}`),
			wantsErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := json.Marshal(tt.response)
			if err != nil && !tt.wantsErr || err == nil && tt.wantsErr {
				t.Fatal(err)
			}
			// in this case, we don't need to compare the data,
			// we do want to return early.
			if tt.wantsErr && err != nil {
				return
			}
			test.Equals(t, string(raw), string(tt.expected))
		})
	}
}

func TestData_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		name     string
		json     []byte
		expected string
	}{
		{
			name:     "collection",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":{"collection":{"language":"golang","tests":"ok"}}}`),
			expected: "collection",
		},
		{
			name:     "complex response",
			json:     []byte(`{"status":"success","code":200,"message":"","data":{"endpoints":[{"uri":"/","method":"get","grants":[]},{"uri":"/healthz","method":"get","grants":[]}]}}`),
			expected: "endpoints",
		},
		{
			name:     "doube collection",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":{"collection":{"language":"golang","tests":"ok"},"collection2":{"language":"golang","tests":"ok"}}}`),
			expected: "",
		},
		{
			name:     "object",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":[{"language":"golang","tests":"ok"}]}`),
			expected: "",
		},
		{
			name:     "k/v pairs inside object",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":{"test":"hello", "test2":"hello2"}}`),
			expected: "",
		},
		{
			name:     "double nested objects",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":[{"collection":{"language":"golang","tests":"ok"}},{"collection2":{"language":"golang","tests":"ok"}}]}`),
			expected: "",
		},
		{
			name:     "empty arrays",
			json:     []byte(`{"status":"ok","code":200,"message":"","data":{"obj1":[],"obj2":[],"obj3":[]}}`),
			expected: "",
		},
		{
			name:     "empty json",
			json:     []byte(`{}`),
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var resp *Response
			if err := json.Unmarshal(tc.json, &resp); err != nil {
				t.Fail()
			}
			if resp.Data != nil {
				if resp.Data.Type != tc.expected {
					t.Fail()
				}
			}
		})
	}
}

func TestResponse_WriteTo(t *testing.T) {
	h := httptest.NewRecorder()
	type fields struct {
		Status  string
		Code    int
		Message string
		Data    *Data
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "200 response",
			fields: fields{
				Code:    http.StatusOK,
				Data:    nil,
				Message: "",
				Status:  "ok",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Response{
				Code:    tt.fields.Code,
				Message: tt.fields.Message,
				Data:    tt.fields.Data,
			}.WriteTo(h)
		})
	}
}

func TestResponse_WriteTo204(t *testing.T) {
	r := &Response{
		Code:    http.StatusNoContent,
		Message: "message",
		Data:    &Data{Type: "type", Content: "content"},
	}

	w := httptest.NewRecorder()
	if err := r.WriteTo(w); err != nil {
		t.Fatalf("unexpected error writing to buffer: %v", err)
	}

	if w.Code != r.Code {
		t.Errorf("exp: %v, got: %v", r.Code, w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("exp: %q, got: %q", "", w.Body.String())
	}
}

func TestDBError(t *testing.T) {
	tests := []struct {
		name   string
		format string
		err    error
		want   *Response
	}{
		{
			name: "internal error",
			err:  errors.New("some error"),
			want: &Response{Code: http.StatusInternalServerError, Message: "db error: some error"},
		},
		{
			name:   "internal error errorf",
			format: "oh noes: %v",
			err:    errors.New("some error"),
			want:   &Response{Code: http.StatusInternalServerError, Message: "oh noes: some error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Response
			if tt.format != "" {
				got = DBErrorf(tt.format, tt.err)
			} else {
				got = DBError(tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "syntax error",
			args: args{
				err: &json.SyntaxError{
					Offset: 99,
				},
			},
			want: &Response{Code: http.StatusUnprocessableEntity, Message: "json error: "},
		},
		{
			name: "any other error",
			args: args{err: errors.New("some error")},
			want: &Response{Code: http.StatusUnprocessableEntity, Message: "json error: some error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONError() = %v, want %v", got, tt.want)
			}
		})
	}
}
