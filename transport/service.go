package microservicetransport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-transport-golang/config"
	"github.com/LUSHDigital/microservice-transport-golang/domain"
)

// Service - Responsible for communication with a service.
type Service struct {
	Branch         string        // VCS branch the service is built from.
	CurrentRequest *http.Request // Current HTTP request being actioned.
	Environment    string        // CI environment the service operates in.
	Namespace      string        // Namespace of the service.
	Name           string        // Name of the service.
	Version        int           // Major API version of the service.
	Client         *http.Client  // http client implementation
}

// NewService - prepares a new service with the provided parameters and client.
func NewService(client *http.Client, branch, env, namespace, name string) *Service {
	return &Service{
		Branch:      branch,
		Name:        name,
		Environment: env,
		Namespace:   namespace,
		Client:      client,
	}
}

// Call - Do the current service request.
func (s *Service) Call() (*http.Response, error) {
	return s.Client.Do(s.CurrentRequest)
}

// Dial - Create a request to a service resource.
func (s *Service) Dial(request *Request) error {
	var err error

	// Make any alterations based upon the namespace.
	switch s.Namespace {
	case "aggregators":
		if !strings.HasPrefix(s.Name, config.AggregatorDomainPrefix) {
			s.Name = strings.Join([]string{config.AggregatorDomainPrefix, s.Name}, "-")
		}
	}

	// Determine the service namespace to use based on the service version.
	serviceNamespace := s.Name
	if s.Version != 0 {
		serviceNamespace = fmt.Sprintf("%s-%d", serviceNamespace, s.Version)
	}

	// Get the name of the service.
	dnsName := domain.BuildServiceDNSName(s.Name, s.Branch, s.Environment, serviceNamespace)

	// Build the resource URL.
	resourceURL := fmt.Sprintf("%s://%s/%s", request.getProtocol(), dnsName, request.Resource)

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
