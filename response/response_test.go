package response

import (
	"reflect"
	"testing"

	"fmt"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// An example data type.
var (
	// An example data set for testing with.
	expectedResponseData = map[string]interface{}{
		"tests":    "ok",
		"language": "golang",
	}
	// example Data struct
	preparedData = &Data{
		Type:    "tests",
		Content: expectedResponseData,
	}

	// An example response object for testing with.
	expectedResponse = &Response{
		Status:  StatusOk,
		Code:    200,
		Message: "",
		Data: &Data{
			Type:    "tests",
			Content: expectedResponseData,
		},
	}

	// An example response object (with no data) for testing with.
	expectedResponseNoData = &Response{
		Status:  StatusOk,
		Code:    200,
		Message: "",
	}

	// the expected error in case type is missing
	errorTypeEmptyWhenDataProvided = "data provided, type cannot be empty"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name     string
		code     int
		status   string
		message  string
		typ      string
		data     *Data
		expected *Response
	}{
		{
			name:     "response valid",
			code:     200,
			status:   StatusOk,
			message:  "",
			data:     preparedData,
			expected: expectedResponse,
		},
		{
			name:     "response no data",
			code:     200,
			status:   StatusOk,
			message:  "",
			data:     nil,
			expected: expectedResponseNoData,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp := New(tc.code, tc.status, tc.message, tc.data)

			if !reflect.DeepEqual(resp, tc.expected) {
				t.Errorf("want: %v\ngot: %v", tc.expected, resp)
			}
		})
	}
}

func TestResponse_ExtractData(t *testing.T) {
	resp := New(200, StatusOk, "", preparedData)
	//
	// Extract the data.
	var dst map[string]interface{}
	extractedData := resp.ExtractData("tests", dst)
	//
	// Compare the data.
	if reflect.DeepEqual(dst, resp.Data.Map()["test"]) {
		t.Errorf("TestExtractData: Expected %v got %v", resp.Data.Map()["tests"], extractedData)
	}

	// test with broken data as well
	resp = New(200, StatusOk, "", &Data{
		Content: expectedResponseData,
	})
	//
	// Extract the data.
	dst = nil
	extractedData = resp.ExtractData("tests", dst)
	//
	// Compare the data.
	if reflect.DeepEqual(dst, nil) {
		t.Errorf("TestExtractData: Expected %v got %v", resp.Data.Map()["tests"], extractedData)
	}
}

func TestData_MarshalJSON(t *testing.T) {
	tt := []struct {
		name string
		data Data
	}{
		{
			name: "correct data",
			data: Data{
				Type:    "testCollection",
				Content: map[string]interface{}{"test": "test"},
			},
		},
		{
			name: "missing data",
			data: Data{
				Type:    "testCollection",
				Content: map[string]interface{}{},
			},
		},
		{
			name: "missing type",
			data: Data{
				Type:    "",
				Content: map[string]interface{}{"test": "test"},
			},
		},
		{
			name: "missing all data",
			data: Data{
				Type:    "",
				Content: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.data.MarshalJSON()
			if err != nil && err.Error() != errorTypeEmptyWhenDataProvided {
				t.Errorf("test '%v' failed with error: %v", tc.name, err)
			}
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

func BenchmarkData_UnmarshalJSON(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	b.ReportAllocs()
	body := []byte(`{"status":"ok","code":200,"message":"","data":{"collection":{"language":"golang","tests":"ok"}}}`)
	for i := 0; i < b.N; i++ {
		var resp *Response
		json.Unmarshal(body, &resp)
	}
}

func TestData_Map(t *testing.T) {
	tt := []struct {
		name     string
		data     Data
		expected map[string]interface{}
	}{
		{
			name: "map valid data",
			data: Data{
				Type: "things",
				Content: map[string]interface{}{
					"thing_one": "a thing",
					"thing_two": "another thing",
				},
			},
			expected: map[string]interface{}{
				"things": map[string]interface{}{
					"thing_one": "a thing",
					"thing_two": "another thing",
				},
			},
		},
		{
			name: "map invalid data",
			data: Data{
				Content: map[string]interface{}{
					"thing_one": "a thing",
					"thing_two": "another thing",
				},
			},
			expected: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.expected, tc.data.Map()) {
				t.Errorf("want: %v, got: %v", tc.expected, tc.data.Map())
			}
		})
	}
}

func BenchmarkData_MarshalJSON(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	data := Data{
		Type: "test",
		Content: map[string]interface{}{
			"test1": "test1",
			"test2": "test2",
			"test3": "test3",
		},
	}

	for i := 0; i < b.N; i++ {
		data.MarshalJSON()
	}
}

func BenchmarkData_MarshalJSON_MissingType(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	data := Data{
		Content: map[string]interface{}{
			"test1": "test1",
			"test2": "test2",
			"test3": "test3",
		},
	}

	for i := 0; i < b.N; i++ {
		data.MarshalJSON()
	}
}

func BenchmarkData_Map(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	data := Data{
		Content: map[string]interface{}{
			"thing_one": "a thing",
			"thing_two": "another thing",
		},
	}

	for i := 0; i < b.N; i++ {
		data.Map()
	}
}

func BenchmarkResponse_ExtractData(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	resp := New(200, StatusOk, "", preparedData)
	for i := 0; i < b.N; i++ {
		var dst map[string]interface{}
		resp.ExtractData("tests", dst)
	}
}

// ExampleNew - Example usage for the New function.
func ExampleNew() {
	data := &Data{
		Type: "things",
		Content: map[string]interface{}{
			"thing_one": "a thing",
			"thing_two": "another thing",
		},
	}

	resp := New(200, StatusOk, "test message", data)
	fmt.Printf("%+v", resp)
}
