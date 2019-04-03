/*Package core consists of a set of packages which are used in writing micro-service
applications.

Each package defines conventional ways of handling common tasks,
as well as a suite of tests to verify their behaviour.
*/
package core

import (
	"context"
	"io"
	"log"

	"github.com/LUSHDigital/core/env"
)

// Service represents the minimal information required to define a working service.
type Service struct {
	// Name represents the name of the service, typically the same as the github repository.
	Name string `json:"service_name"`
	// Type represents the type of the service, eg. service or aggregator.
	Type string `json:"service_type"`
	// Version represents the SVC tagged version of the service.
	Version string `json:"service_version"`
}

// NewService reads a service definition from the environment.
func NewService() *Service {
	return &Service{
		Name:    env.MustGet("SERVICE_NAME"),
		Type:    env.MustGet("SERVICE_TYPE"),
		Version: env.MustGet("SERVICE_VERSION"),
	}
}

// ServiceWorker represents the behaviour for running a service worker.
type ServiceWorker interface {
	// Run should be a blocking operation
	Run(context.Context, io.Writer) error
}

type writer func(b []byte) (int, error)

func (f writer) Write(b []byte) (int, error) {
	return f(b)
}

// StartWorkers will start the given service workers and block block indefinitely.
func (s *Service) StartWorkers(ctx context.Context, workers ...ServiceWorker) {
	var out = writer(func(b []byte) (int, error) {
		log.Println(b)
		return len(b), nil
	})
	var work = func(worker ServiceWorker) {
		log.Fatalln(worker.Run(ctx, out))
	}
	log.Printf("starting %s: %s (%s)\n", s.Type, s.Name, s.Version)
	for _, worker := range workers {
		go work(worker)
	}
	select {} // blocks the thread indefinitely
}
