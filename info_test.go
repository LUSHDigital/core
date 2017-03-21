package microservicecore

import (
	"fmt"
	"os"
	"testing"
)

// Example environment variables for testing.
var exampleEnvVars = map[string]string{
	"SERVICE_NAME":    "example-service",
	"SERVICE_TYPE":    "examples",
	"SERVICE_SCOPE":   "testing",
	"SERVICE_VERSION": "0.0.1",
}

// TestServiceInfo - Test the GetMicroserviceInfo function is working.
func TestServiceInfo(t *testing.T) {
	// Set our expected environment variables.
	for key, value := range exampleEnvVars {
		os.Setenv(key, value)
	}

	// Get the service info.
	info := GetMicroserviceInfo()

	// Check each expected env var.
	if info.ServiceName != os.Getenv("SERVICE_NAME") {
		t.Error(fmt.Sprintf("Expected %v, got %v", os.Getenv("SERVICE_NAME"), info.ServiceName))
	}

	if info.ServiceType != os.Getenv("SERVICE_TYPE") {
		t.Error(fmt.Sprintf("Expected %v, got %v", os.Getenv("SERVICE_TYPE"), info.ServiceType))
	}

	if info.ServiceScope != os.Getenv("SERVICE_SCOPE") {
		t.Error(fmt.Sprintf("Expected %v, got %v", os.Getenv("SERVICE_SCOPE"), info.ServiceScope))
	}

	if info.ServiceVersion != os.Getenv("SERVICE_VERSION") {
		t.Error(fmt.Sprintf("Expected %v, got %v", os.Getenv("SERVICE_VERSION"), info.ServiceVersion))
	}
}

// ExampleGetMicroserviceInfo - Example for GetMicroserviceInfo.
func ExampleGetMicroserviceInfo() {
	envVars := map[string]string{
		"SERVICE_NAME":    "example-service",
		"SERVICE_TYPE":    "examples",
		"SERVICE_SCOPE":   "testing",
		"SERVICE_VERSION": "0.0.1",
	}

	// Set our expected environment variables.
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Get the service info.
	info := GetMicroserviceInfo()

	// Print the details.
	fmt.Println(info.ServiceName)
	fmt.Println(info.ServiceType)
	fmt.Println(info.ServiceScope)
	fmt.Println(info.ServiceVersion)

	// Output:
	// example-service
	// examples
	// testing
	// 0.0.1
}
