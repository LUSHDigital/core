package transport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LUSHDigital/microservice-core-golang/transport/config"
	"github.com/LUSHDigital/microservice-core-golang/transport/domain"
)

// Resource - Useful for defining a remote resource without an attached transport.
type Resource struct {
	Branch         string        // VCS branch the service is built from.
	CurrentRequest *http.Request // Current HTTP request being actioned.
	Environment    string        // CI environment the service operates in.
	Namespace      string        // Namespace of the service.
	Name           string        // Name of the service.
	Version        int           // Major API version of the service.
}

// DNSPath returns internal dns path for the resource
func (r *Resource) DNSPath() string {
	// Make any alterations based upon the namespace.
	switch r.Namespace {
	case "aggregators":
		if !strings.HasPrefix(r.Name, config.AggregatorDomainPrefix) {
			r.Name = strings.Join([]string{config.AggregatorDomainPrefix, r.Name}, "-")
		}
	}
	// Determine the service namespace to use based on the service version.
	serviceNamespace := r.Name
	if r.Version != 0 {
		serviceNamespace = fmt.Sprintf("%r-%d", serviceNamespace, r.Version)
	}
	return serviceNamespace
}

// Service - Responsible for communication with a service.
type Service struct {
	Resource
	Client *http.Client // http client implementation
}

// NewService - prepares a new service with the provided parameters and client.
func NewService(client *http.Client, branch, env, namespace, name string) *Service {
	return &Service{
		Resource: Resource{
			Branch:      branch,
			Name:        name,
			Environment: env,
			Namespace:   namespace,
		},
		Client: client,
	}
}

// Call - Do the current service request.
func (s *Service) Call() (*http.Response, error) {
	return s.Client.Do(s.CurrentRequest)
}

// Dial - Create a request to a service resource.
func (s *Service) Dial(request *Request) error {
	var err error

	serviceNamespace := s.DNSPath()
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
