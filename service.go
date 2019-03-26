/*Package core consists of a set of packages which are used in writing micro-service
applications.

Each package defines conventional ways of handling common tasks,
as well as a suite of tests to verify their behaviour.
*/
package core

import (
	"github.com/LUSHDigital/core/env"
)

// Service represents the minimal information required to define a working
// micro-service.
//
// This information is purely for the reader's benefit, and
// should ideally be used to report on a given service's characteristics when
// interrogating it's info or health endpoints or any similar pattern you may
// use.
type Service struct {
	Name    string `json:"service_name"`
	Type    string `json:"service_type"`
	Version string `json:"service_version"`
}

// GetService reads a service definition from the environment.
func GetService() *Service {
	return &Service{
		Name:    env.MustGet("SERVICE_NAME"),
		Type:    env.MustGet("SERVICE_TYPE"),
		Version: env.MustGet("SERVICE_VERSION"),
	}
}
