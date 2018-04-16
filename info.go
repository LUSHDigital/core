package microservicecore

import (
	"net/http"
	"github.com/LUSHDigital/microservice-core-golang/env"
)

// MicroserviceInfo - Represents information about this microservice.
type MicroserviceInfo struct {
	ServiceName    string  `json:"service_name"`
	ServiceType    string  `json:"service_type"`
	ServiceScope   string  `json:"service_scope"`
	ServiceVersion string  `json:"service_version"`
	Endpoints      []Route `json:"endpoints"`
}

// Route defines an HTTP route
type Route struct {
	Path    string                                   `json:"uri"`
	Method  string                                   `json:"method"`
	Handler func(http.ResponseWriter, *http.Request) `json:"-"`
}

// GetMicroserviceInfo - Get the information about this microservice.
//
// Return:
//     *MicroserviceInfo - Object representing this microservice.
func GetMicroserviceInfo() *MicroserviceInfo {
	return &MicroserviceInfo{
		ServiceName:    env.MustGet("SERVICE_NAME"),
		ServiceType:    env.MustGet("SERVICE_TYPE"),
		ServiceScope:   env.MustGet("SERVICE_SCOPE"),
		ServiceVersion: env.MustGet("SERVICE_VERSION"),
	}
}
