package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/LUSHDigital/microservice-core-golang/transport/models"
	transportErrors "github.com/LUSHDigital/microservice-core-golang/transport/errors"
	"github.com/LUSHDigital/microservice-core-golang/transport/config"
	"github.com/LUSHDigital/microservice-core-golang/transport/domain"
	"io/ioutil"
)

// AuthCredentials - Credentials needed to authenticate for a cloud service.
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CloudService - Responsible for communication with a cloud service.
type CloudService struct {
	Service                      // Inherit all properties of a normal service.
	Credentials *AuthCredentials // Authentication credentials for cloud service calls.
	Client      *http.Client
}

// NewCloudService - Prepare a new CloudService struct with the provided parameters.
func NewCloudService(client *http.Client, branch, env, namespace, name string, credentials *AuthCredentials) *CloudService {
	return &CloudService{
		Service: Service{
			Branch:      branch,
			Environment: env,
			Namespace:   namespace,
			Name:        name,
			Client:      client,
		},
		Client:      client,
		Credentials: credentials,
	}
}

// authenticate - Authenticate against the API gateway and return an auth token.
func (c *CloudService) authenticate(request *Request) (*models.Token, error) {
	//loginBody := new(bytes.Buffer)
	loginBody, err := json.Marshal(c.Credentials)
	if err != nil {
		return nil, fmt.Errorf("cannot encode json: %s", err)
	}

	loginReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.GetAPIGatewayURL(request), "login"), bytes.NewBuffer(loginBody))
	if err != nil {
		return nil, fmt.Errorf("cannot build login request: %s", err)
	}

	loginResp, err := c.Client.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("cannot perform login request: %s", err)
	}

	// Decode response.
	serviceResponse := response.Response{}

	content, err := ioutil.ReadAll(loginResp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}
	if err := json.Unmarshal(content, &serviceResponse); err != nil {
		return nil, fmt.Errorf("cannot decode login response: %v", err)
	}

	// Handle any error codes.
	switch loginResp.StatusCode {
	// Custom error for login failed.
	case http.StatusUnauthorized, http.StatusNotFound:
		return nil, transportErrors.LoginUnauthorisedError{}

		// 200 and 304 are all good.
	case http.StatusOK, http.StatusNotModified:
		break

		// Something somewhere broken!
	default:
		return nil, fmt.Errorf("api gateway login failed: %s", serviceResponse.Message)
	}

	// Extract the consumer from the response.
	var consumer *models.Consumer
	consumerErr := serviceResponse.ExtractData("consumer", &consumer)
	if consumerErr != nil {
		return nil, fmt.Errorf("could not extract consumer data: %s", consumerErr)
	}

	if len(consumer.Tokens) == 0 {
		return nil, transportErrors.ConsumerHasNoTokensError{}
	}

	return consumer.Tokens[0], nil
}

// GetAPIGatewayURL - Get the url of the API gateway.
func (c *CloudService) GetAPIGatewayURL(request *Request) string {
	// Check if a full URL is set in the environment.
	if config.GetGatewayURL() != "" {
		return config.GetGatewayURL()
	}

	// Fallback to constructing the URL ourselves.
	if c.Environment == "staging" {
		return fmt.Sprintf("%s://%s-%s.%s", request.getProtocol(), config.GetGatewayURI(), c.Environment, config.GetServiceDomain())
	}

	return fmt.Sprintf("%s://%s.%s", request.getProtocol(), config.GetGatewayURI(), config.GetServiceDomain())
}

// Call - Do the current service request.
func (c *CloudService) Call() (*http.Response, error) {
	return c.Client.Do(c.CurrentRequest)
}

// Dial - Create a request to a service resource.
func (c *CloudService) Dial(request *Request) error {
	if c.Credentials.Email == "" || c.Credentials.Password == "" {
		return errors.New("cannot authenticate for cloud service: missing credentials")
	}

	token, err := c.authenticate(request)
	if err != nil {
		return fmt.Errorf("cannot authenticate for cloud service: %s", err)
	}

	// Make any alterations based upon the namespace.
	switch c.Namespace {
	case "aggregators":
		c.Name = strings.Join([]string{config.AggregatorDomainPrefix, c.Name}, "-")
	}

	cloudServiceURL := domain.BuildCloudServiceURL(c.GetAPIGatewayURL(request), c.Namespace, c.Name)

	// Build the resource URL.
	resourceURL := fmt.Sprintf("%s/%s", cloudServiceURL, request.Resource)

	// Append the query string if we have any.
	if len(request.Query) > 0 {
		resourceURL = fmt.Sprintf("%s?%s", resourceURL, request.Query.Encode())
	}

	// Create the request.
	var reqErr error
	c.CurrentRequest, reqErr = http.NewRequest(request.Method, resourceURL, request.Body)
	if reqErr != nil {
		return reqErr
	}

	// Set the auth token header.
	c.CurrentRequest.Header.Set(config.AuthHeader, token.PrepareForHTTP())

	// Add the version header to the request if applicable.
	if c.Version != 0 {
		c.CurrentRequest.Header.Set(config.ServiceVersionHeader, strconv.Itoa(c.Version))
	}

	// Add the headers.
	for key, value := range request.Headers {
		c.CurrentRequest.Header.Set(key, value)
	}

	return nil
}

// GetName - Get the name of the service
func (c *CloudService) GetName() string {
	return c.Name
}
