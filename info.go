package microservicecore

import "github.com/LUSHDigital/microservice-core-golang/routing"

// MicroserviceInfo - Represents information about this microservice.
type MicroserviceInfo struct {
	ServiceName    string          `json:"service_name"`
	ServiceType    string          `json:"service_type"`
	ServiceScope   string          `json:"service_scope"`
	ServiceVersion string          `json:"service_version"`
	Endpoints      []routing.Route `json:"endpoints"`
}

// GetMicroserviceInfo - Get the information about this microservice.
//
// Return:
//     *MicroserviceInfo - Object representing this microservice.
func GetMicroserviceInfo() *MicroserviceInfo {
	return &MicroserviceInfo{
		ServiceName:    GetEnvOrFail("SERVICE_NAME"),
		ServiceType:    GetEnvOrFail("SERVICE_TYPE"),
		ServiceScope:   GetEnvOrFail("SERVICE_SCOPE"),
		ServiceVersion: GetEnvOrFail("SERVICE_VERSION"),
	}
}
