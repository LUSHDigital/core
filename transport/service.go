package transport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LUSHDigital/core/transport/config"
)

// NewResource returns a new service
func NewResource(branch, env, namespace, name string) *Resource {
	// Make any alterations based upon the namespace
	switch namespace {
	case "aggregators":
		if !strings.HasPrefix(name, config.AggregatorDomainPrefix) {
			name = strings.Join([]string{config.AggregatorDomainPrefix, name}, "-")
		}
	}
	return &Resource{
		Branch:      branch,
		Name:        name,
		Environment: env,
		Namespace:   namespace,
	}
}

// Resource defines a remote service
type Resource struct {
	Branch      string // VCS branch the service is built from.
	Environment string // CI environment the service operates in.
	Namespace   string // Namespace of the service.
	Name        string // Name of the service.
	Version     int    // Major API version of the service.
}

// DomainName returns the resource domain name for the internal DNS
func (s *Service) DomainName() string {
	name := s.Name
	// Determine the service namespace to use based on the service version
	if s.Version != 0 {
		name = fmt.Sprintf("%s-%d", name, s.Version)
	}
	return fmt.Sprintf("%s-%s-%s.%s", s.Name, s.Branch, s.Environment, name)
}

// Service - Responsible for communication with a service.
type Service struct {
	Resource
	CurrentRequest *http.Request // Current HTTP request being actioned.
	Client         *http.Client  // http client implementation
}

// NewService - prepares a new service with the provided parameters and client.
func NewService(client *http.Client, branch, env, namespace, name string) *Service {
	return &Service{
		Resource: *NewResource(branch, env, namespace, name),
		Client:   client,
	}
}

// Call - Do the current service request.
func (s *Service) Call() (*http.Response, error) {
	return s.Client.Do(s.CurrentRequest)
}

// Dial - Create a request to a service resource.
func (s *Service) Dial(request *Request) error {
	var err error

	// Build the resource URL.
	resourceURL := fmt.Sprintf("%s://%s/%s", request.getProtocol(), s.DomainName(), request.Resource)

	// Append the query string if we have any.
	if len(request.Query) > 0 {
		resourceURL = fmt.Sprintf("%s?%s", resourceURL, request.Query.Encode())
	}

	// Create the request.
	s.CurrentRequest, err = http.NewRequest(request.Method, resourceURL, request.Body)

	// Add the headers.
	for key, value := range request.Headers {
		s.CurrentRequest.Header.Set(key, value)
	}

	return err
}

// GetName - Get the name of the service
func (s *Service) GetName() string {
	return s.Name
}
