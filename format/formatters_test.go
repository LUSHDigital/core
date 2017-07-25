package format

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/response"
)

var responseData = response.Data{
	Type: "tests",
	Content: map[string]interface{}{
		"tests":    "ok",
		"language": "golang",
	},
}

// The expected response body.
var expectedBody = []byte(`{"status":"ok","code":200,"message":"","data":{"tests":{"language":"golang","tests":"ok"}}}`)

// TestJSONResponseFormatter - Test a JSON response over HTTP.
func TestJSONResponseFormatter(t *testing.T) {
	// Start a HTTP server for testing purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response.
		resp := response.New(200, response.StatusOk, "", &responseData)
		// Format the response as JSON.
		JSONResponseFormatter(w, resp)
	}))
	defer ts.Close()

	// Check the server is working.
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Check the response code.
	if res.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("Expected %d, got %d", http.StatusOK, res.StatusCode))
	}

	// Get the response body.
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Check the response is what we expect.
	if !jsonEquals(body, expectedBody) {
		t.Error(fmt.Sprintf("Expected %v, got %v", string(body), string(expectedBody)))
	}
}

// jsonEquals - Compare the equality of two JSON byte arrays.
//
// Params:
//     jsonA []byte - The first JSON byte array.
//     jsonB []byte - The second JSON byte array.
//
// Return:
//     bool - Are the JSON byte arrays equal.
func jsonEquals(jsonA, jsonB []byte) bool {
	var interfaceA interface{}
	aErr := json.Unmarshal(jsonA, &interfaceA)
	if aErr != nil {
		log.Fatal(aErr)
	}

	var interfaceB interface{}
	bErr := json.Unmarshal(jsonB, &interfaceB)
	if bErr != nil {
		log.Fatal(bErr)
	}

	return reflect.DeepEqual(interfaceA, interfaceB)
}
