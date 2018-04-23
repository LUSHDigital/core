package config

import "os"

const (
	// AuthHeader - Name of the HTTP header to use for authentication and
	// authorization.
	AuthHeader = "Authorization"

	// AuthHeaderPrefix - Prefix expected for the HTTP auth header value.
	AuthHeaderPrefix = "Bearer"

	// ProtocolHTTP - Protocol string for non-ssl requests.
	ProtocolHTTP = "http"

	// ProtocolHTTPS - Protocol string for ssl requests.
	ProtocolHTTPS = "https"

	// RequestKey - Name of the context key used to pass a service request
	// between layers of middleware.
	RequestKey = "request"

	// ServiceVersionHeader - Name of the HTTP header to use for service version.
	ServiceVersionHeader = "x-service-version"

	// AggregatorDomainPrefix - The prefix value used for aggregator domains.
	AggregatorDomainPrefix = "agg"
)

// GetServiceDomain - Get the top level domain of the service environment.
func GetServiceDomain() string {
	return os.Getenv("SOA_DOMAIN")
}

// GetGatewayURI - Get the URI of the API gateway.
func GetGatewayURI() string {
	return os.Getenv("SOA_GATEWAY_URI")
}

// GetGatewayURL - Get the URL of the API gateway.
func GetGatewayURL() string {
	return os.Getenv("SOA_GATEWAY_URL")
}
