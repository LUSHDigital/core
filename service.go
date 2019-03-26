package microservicecore

import (
	"github.com/LUSHDigital/core/env"
)

// Service represents the minimal information required to define a working micro-service.
type Service struct {
	Name    string `json:"service_name"`
	Type    string `json:"service_type"`
	Version string `json:"service_version"`
}

// GetService derives a service from the environment.
func GetService() *Service {
	return &Service{
		Name:    env.MustGet("SERVICE_NAME"),
		Type:    env.MustGet("SERVICE_TYPE"),
		Version: env.MustGet("SERVICE_VERSION"),
	}
}
