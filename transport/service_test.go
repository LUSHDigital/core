package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-core-golang/transport/config"
)

func TestService_Dial(t *testing.T) {
	tt := []struct {
		name        string
		service     Service
		request     *Request
		postData    map[string]string
		expectedURL string
	}{
		{

			name: "name has agg- prefix",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "agg-myservice",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedURL: "http://agg-myservice-master-staging.agg-myservice/things",
		},
		{
			name: "GET HTTP",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedURL: "http://myservice-master-staging.myservice/things",
		},
		{
			name: "GET HTTPS",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedURL: "https://myservice-master-staging.myservice/things",
		},
		{
			name: "GET with query HTTP",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
				Query: url.Values{
					"baz": []string{"qux"},
					"foo": []string{"bar"},
				},
			},
			expectedURL: "http://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "GET with query HTTPS",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
				Query: url.Values{
					"foo": []string{"bar"},
					"baz": []string{"qux"},
				},
			},
			expectedURL: "https://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "POST HTTP",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedURL: "http://myservice-master-staging.myservice/things",
		},
		{
			name: "POST HTTPS",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedURL: "https://myservice-master-staging.myservice/things",
		},
		{
			name: "POST with query HTTP",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
				Query: url.Values{
					"baz": []string{"qux"},
					"foo": []string{"bar"},
				},
			},
			expectedURL: "http://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "POST with query HTTPS",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
				Query: url.Values{
					"foo": []string{"bar"},
					"baz": []string{"qux"},
				},
			},
			expectedURL: "https://myservice-master-staging.myservice/things?baz=qux&foo=bar",
		},
		{
			name: "POST with headers",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			postData: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
			request: &Request{
				Method:   http.MethodPost,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
				Headers: map[string]string{
					"Content-Type":    "application/json",
					"Accept-Language": "en-GB",
				},
			},
			expectedURL: "https://myservice-master-staging.myservice/things",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Add a body for POST requests.
			if tc.request.Method == http.MethodPost && len(tc.postData) > 0 {
				postBody := new(bytes.Buffer)
				json.NewEncoder(postBody).Encode(tc.postData)

				tc.request.Body = ioutil.NopCloser(postBody)
			}

			err := tc.service.Dial(tc.request)
			if err != nil {
				t.Errorf("TestService_Dial: %s: %s", tc.name, err)
			}

			if tc.service.CurrentRequest.Method != tc.request.Method {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.request.Method, tc.service.CurrentRequest.Method)
			}

			if tc.service.CurrentRequest.URL.String() != tc.expectedURL {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.expectedURL, tc.service.CurrentRequest.URL.String())
			}

			for key, value := range tc.request.Headers {
				if value != tc.service.CurrentRequest.Header.Get(key) {
					t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, value, tc.service.CurrentRequest.Header.Get(key))
				}
			}
		})
	}
}

func TestService_GetName(t *testing.T) {
	tt := []struct {
		name         string
		service      Service
		expectedName string
	}{
		{
			name: "Normal",
			service: Service{
				Resource: Resource{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
			},
			expectedName: "myservice",
		},
		{
			name: "Crazy",
			service: Service{
				Resource: Resource{
					Branch:      "massdsdfsdjf89uter",
					Environment: "sdfsdf34341",
					Namespace:   "l1j2312klj3k21j3",
					Name:        "-sf9s9f9ds0f9-",
				},
			},
			expectedName: "-sf9s9f9ds0f9-",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.service.GetName() != tc.expectedName {
				t.Errorf("TestService_GetName: %s: expected %v got %v", tc.name, tc.expectedName, tc.service.GetName())
			}
		})
	}
}

func ExampleService_Dial() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleService_Dial_post() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	postData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(postData)

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Body:     ioutil.NopCloser(postBody),
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleService_Dial_query() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Query: url.Values{
			"foo": []string{"bar"},
			"baz": []string{"qux"},
		},
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleService_Dial_headers() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Language": "en-GB",
		},
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleService_Dial_customClient() {
	// Instantiate the service.
	myService := NewService(
		&http.Client{
			Timeout: 5 * time.Microsecond,
		},
		"master",
		"staging",
		"services",
		"myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}
}

func ExampleService_Call() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}

	// Do the request.
	myServiceResp, err := myService.Call()
	if err != nil {
		fmt.Printf("call err: %s", err)
	}

	// Make sure we close the body once we're done.
	defer myServiceResp.Body.Close()

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(myServiceResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		fmt.Printf("decode err: %s", err)
	}

	// Handle any error codes.
	switch serviceResponse.Code {
	// Custom error for grants not found.
	case http.StatusNotFound:
		fmt.Println("response err: not found")

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		fmt.Println("response err: internal server error")
	}
}

func ExampleService_Call_post() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	postData := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	postBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(postData)

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Body:     ioutil.NopCloser(postBody),
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}

	// Do the request.
	myServiceResp, err := myService.Call()
	if err != nil {
		fmt.Printf("call err: %s", err)
	}

	// Make sure we close the body once we're done.
	defer myServiceResp.Body.Close()

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(myServiceResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		fmt.Printf("decode err: %s", err)
	}

	// Handle any error codes.
	switch serviceResponse.Code {
	// Custom error for grants not found.
	case http.StatusNotFound:
		fmt.Println("response err: not found")

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		fmt.Println("response err: internal server error")
	}
}

func ExampleService_Call_query() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Query: url.Values{
			"foo": []string{"bar"},
			"baz": []string{"qux"},
		},
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}

	// Do the request.
	myServiceResp, err := myService.Call()
	if err != nil {
		fmt.Printf("call err: %s", err)
	}

	// Make sure we close the body once we're done.
	defer myServiceResp.Body.Close()

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(myServiceResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		fmt.Printf("decode err: %s", err)
	}

	// Handle any error codes.
	switch serviceResponse.Code {
	// Custom error for grants not found.
	case http.StatusNotFound:
		fmt.Println("response err: not found")

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		fmt.Println("response err: internal server error")
	}
}

func ExampleService_Call_headers() {
	// Instantiate the service.
	myService := NewService(DefaultHTTPClient(), "master", "staging", "services", "myservice")

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodPost,
		Resource: "things",
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"Accept-Language": "en-GB",
		},
	}

	// Dial it up!
	err := myService.Dial(myServiceThingsRequest)
	if err != nil {
		fmt.Printf("dial err: %s", err)
	}

	// Do the request.
	myServiceResp, err := myService.Call()
	if err != nil {
		fmt.Printf("call err: %s", err)
	}

	// Make sure we close the body once we're done.
	defer myServiceResp.Body.Close()

	// Decode response.
	serviceResponse := response.Response{}
	jsonErr := json.NewDecoder(myServiceResp.Body).Decode(&serviceResponse)
	if jsonErr != nil {
		fmt.Printf("decode err: %s", err)
	}

	// Handle any error codes.
	switch serviceResponse.Code {
	// Custom error for grants not found.
	case http.StatusNotFound:
		fmt.Println("response err: not found")

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		fmt.Println("response err: internal server error")
	}
}
