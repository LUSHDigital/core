package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-core-golang/transport/models"
	"github.com/LUSHDigital/microservice-core-golang/transport/config"
)

func TestCloudService_Dial(t *testing.T) {
	// Start a HTTP server to act as a fake API gateway.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := response.New(http.StatusOK, "", &response.Data{
			Type: "consumer",
			Content: models.Consumer{
				Tokens: []*models.Token{
					{
						Type:  "JWT",
						Value: "xxxx.xxxx.xxxx",
					},
				},
			},
		})
		resp.WriteTo(w)
	}))
	defer ts.Close()

	// Set the URL of the fake gateway as an environment variable.
	os.Setenv("SOA_GATEWAY_URL", ts.URL)

	tt := []struct {
		name        string
		service     CloudService
		request     *Request
		postData    map[string]string
		expectedURI string
	}{
		{
			name: "GET HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTP,
			},
			expectedURI: "things",
		},
		{
			name: "GET HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			request: &Request{
				Method:   http.MethodGet,
				Resource: "things",
				Protocol: config.ProtocolHTTPS,
			},
			expectedURI: "things",
		},
		{
			name: "GET with query HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things?baz=qux&foo=bar",
		},
		{
			name: "GET with query HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things?baz=qux&foo=bar",
		},
		{
			name: "POST HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things",
		},
		{
			name: "POST HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things",
		},
		{
			name: "POST with query HTTP",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things?baz=qux&foo=bar",
		},
		{
			name: "POST with query HTTPS",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
					Client:      DefaultHTTPClient(),
				},
				Client: DefaultHTTPClient(),
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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
			expectedURI: "things?baz=qux&foo=bar",
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
				t.Fatalf("TestCloudService_Dial: %s: %s", tc.name, err)
			}

			if tc.service.CurrentRequest.Method != tc.request.Method {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, tc.request.Method, tc.service.CurrentRequest.Method)
			}

			expectedURL := fmt.Sprintf("%s/%s/%s/%s", ts.URL, tc.service.Namespace, tc.service.Name, tc.expectedURI)
			if tc.service.CurrentRequest.URL.String() != expectedURL {
				t.Errorf("TestService_Dial: %s: expected %v got %v", tc.name, expectedURL, tc.service.CurrentRequest.URL.String())
			}
		})
	}
}

func TestCloudService_GetName(t *testing.T) {
	tt := []struct {
		name         string
		service      CloudService
		expectedName string
	}{
		{
			name: "Normal",
			service: CloudService{
				Service: Service{
					Branch:      "master",
					Environment: "staging",
					Namespace:   "services",
					Name:        "myservice",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
				},
			},
			expectedName: "myservice",
		},
		{
			name: "Crazy",
			service: CloudService{
				Service: Service{
					Branch:      "massdsdfsdjf89uter",
					Environment: "sdfsdf34341",
					Namespace:   "l1j2312klj3k21j3",
					Name:        "-sf9s9f9ds0f9-",
				},
				Credentials: &AuthCredentials{
					Email:    "test@test.com",
					Password: "1234",
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

func TestCloudService_GetApiGatewayUrl(t *testing.T) {
	tt := []struct {
		name               string
		gatewayURL         string
		gatewayURI         string
		serviceDomain      string
		expectedGatewayURL string
	}{
		{
			name:               "Just URL",
			gatewayURL:         "https://api-gateway.test.com",
			expectedGatewayURL: "https://api-gateway.test.com",
		},
		{
			name:               "Just URI + domain",
			gatewayURI:         "api-gateway",
			serviceDomain:      "test.com",
			expectedGatewayURL: "http://api-gateway-staging.test.com",
		},
		{
			name:               "URL + URI + domain",
			gatewayURL:         "https://api-gateway.wibble.com",
			gatewayURI:         "api-gateway",
			serviceDomain:      "test.com",
			expectedGatewayURL: "https://api-gateway.wibble.com",
		},
	}

	// Instantiate the service.
	myService := &CloudService{
		Service: Service{
			Branch:      "master",
			Environment: "staging",
			Namespace:   "services",
			Name:        "myservice",
		},
		Credentials: &AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	}

	// Prepare the request.
	myServiceThingsRequest := &Request{
		Method:   http.MethodGet,
		Resource: "things",
	}

	for _, tc := range tt {
		os.Setenv("SOA_DOMAIN", tc.serviceDomain)
		os.Setenv("SOA_GATEWAY_URI", tc.gatewayURI)
		os.Setenv("SOA_GATEWAY_URL", tc.gatewayURL)

		t.Run(tc.name, func(t *testing.T) {
			actualGatewayURL := myService.GetAPIGatewayURL(myServiceThingsRequest)
			if actualGatewayURL != tc.expectedGatewayURL {
				t.Errorf("TestCloudService_GetApiGatewayUrl: %s: expected %v got %v", tc.name, tc.expectedGatewayURL, actualGatewayURL)
			}
		})
	}
}

func ExampleCloudService_Dial() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Dial_post() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Dial_query() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Dial_headers() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Dial_customClient() {
	// Instantiate the service.
	myService := NewCloudService(
		&http.Client{
			Timeout: 500 * time.Microsecond,
		},
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Call() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Call_post() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Call_query() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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

func ExampleCloudService_Call_headers() {
	// Instantiate the service.
	myService := NewCloudService(
		DefaultHTTPClient(),
		"master",
		"staging",
		"services",
		"myservice",
		&AuthCredentials{
			Email:    "test@test.com",
			Password: "1234",
		},
	)

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
