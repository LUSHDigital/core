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

	"strconv"

	"net/url"

	"github.com/LUSHDigital/microservice-core-golang/pagination"
	"github.com/LUSHDigital/microservice-core-golang/response"
)

var (
	rawData = map[string]interface{}{
		"tests":    "ok",
		"language": "golang",
	}

	responseData = response.Data{
		Type:    "tests",
		Content: rawData,
	}
)

// The expected response body.
var expectedBody = []byte(`{"status":"ok","code":200,"message":"","data":{"tests":{"language":"golang","tests":"ok"}}}`)

// TestJSONResponseFormatter - Test a JSON response over HTTP.
func TestJSONResponseFormatter(t *testing.T) {
	// Start a HTTP server for testing purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response.
		resp := response.New(200, "", &responseData)
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
		t.Error(fmt.Sprintf("Expected %v, got %v", string(expectedBody), string(body)))
	}
}

// TestJSONResponseFormatter - Test a JSON response over HTTP. This time it's paginated!
func TestJSONResponseFormatter2(t *testing.T) {
	// Start a HTTP server for testing purposes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		// Get the pagination values.
		var page = 1
		var pageErr error
		if len(params.Get("page")) > 0 {
			page, pageErr = strconv.Atoi(params.Get("page"))
			if pageErr != nil {
				log.Printf("Could not parse 'page' parameter: %s", pageErr)
			}
		}

		var perPage = 5
		var perPageErr error
		if len(params.Get("per_page")) > 0 {
			perPage, perPageErr = strconv.Atoi(params.Get("per_page"))
			if perPageErr != nil {
				log.Printf("Could not parse 'per_page' parameter: %s", perPageErr)
			}
		}

		// Create a paginator.
		paginator, err := pagination.NewPaginator(perPage, page, len(rawData))
		if err != nil {
			log.Printf("Could not create paginator: %s", perPageErr)
		}

		// Hacky way to slice the raw data up.
		returnData := make(map[string]interface{}, perPage)
		var i = 0
		for key, value := range rawData {
			if i < perPage {
				returnData[key] = value
			}

			i++
		}
		responseData.Content = returnData

		// Create a response.
		resp := response.NewPaginated(paginator, 200, response.StatusOk, "", &responseData)

		// Format the response as JSON.
		JSONResponseFormatter(w, resp)
	}))
	defer ts.Close()

	tt := []struct {
		name           string
		page           int
		perPage        int
		expectedOutput []byte
	}{
		{
			name:           "2 items. 1 per page. page 1",
			page:           1,
			perPage:        1,
			expectedOutput: []byte(`{"status":"ok","code":200,"message":"","data":{"tests":{"tests":"ok"}},"pagination":{"total":2,"per_page":1,"current_page":1,"last_page":2,"next_page":2,"prev_page":null}}`),
		},
		{
			name:           "2 items. 1 per page. page 2",
			page:           2,
			perPage:        1,
			expectedOutput: []byte(`{"status":"ok","code":200,"message":"","data":{"tests":{"tests":"ok"}},"pagination":{"total":2,"per_page":1,"current_page":2,"last_page":2,"next_page":null,"prev_page":1}}`),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			queryParams := url.Values{}
			queryParams.Add("page", strconv.Itoa(tc.page))
			queryParams.Add("per_page", strconv.Itoa(tc.perPage))

			//Build the request.
			req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", ts.URL, queryParams.Encode()), nil)
			if err != nil {
				log.Fatal(err)
			}

			// Check the server is working.
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			// Check the response code.
			if res.StatusCode != http.StatusOK {
				t.Error(fmt.Sprintf("TestJSONResponseFormatter2: %s, expected %d, got %d", tc.name, http.StatusOK, res.StatusCode))
			}

			// Get the response body.
			body, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			// Check the response is what we expect.
			if !jsonEquals(body, tc.expectedOutput) {
				t.Error(fmt.Sprintf("TestJSONResponseFormatter2: %s, expected %v, got %v", tc.name, string(tc.expectedOutput), string(body)))
			}
		})
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
