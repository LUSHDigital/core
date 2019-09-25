/*Package core consists of a set of packages which are used in writing micro-service
applications.

Each package defines conventional ways of handling common tasks,
as well as a suite of tests to verify their behaviour.
*/
package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	// tag and ref can be used by build stage with build flags to set things like git commit hash and tag.
	// --ldflags "-X github.com/LUSHDigital/core.tag=${GIT_TAG}"
	// --ldflags "-X github.com/LUSHDigital/core.ref=${GIT_COMMIT_HASH}"
	tag string
	ref string
)

// Service represents the minimal information required to define a working service.
type Service struct {
	// Name represents the name of the service, typically the same as the github repository.
	Name string `json:"name"`
	// Type represents the type of the service, eg. service or aggregator.
	Type string `json:"type"`
	// Version represents the latest version or SVC tag of the service.
	Version string `json:"version"`
	// Revision represents the SVC revision or commit hash of the service.
	Revision string `json:"revision"`
}

// ServiceOption represents behaviour for applying options to a new service.
type ServiceOption interface {
	Apply(*Service)
}

// NewService creates a new service based on
func NewService(name, kind string, opts ...ServiceOption) *Service {
	if v := os.Getenv("SERVICE_VERSION"); v != "" {
		tag = v
	}
	if r := os.Getenv("SERVICE_REVISION"); r != "" {
		ref = r
	}
	return &Service{
		Name:     name,
		Type:     kind,
		Version:  tag,
		Revision: ref,
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
	if s.Name == "" || s.Type == "" {
		log.Fatalln("cannot start without a name or type")
	}
	var out = writer(func(b []byte) (int, error) {
		log.Print(string(b))
		return len(b), nil
	})
	var work = func(worker ServiceWorker) {
		if err := worker.Run(ctx, out); err != nil {
			log.Fatalln(err)
		}
	}
	msg := fmt.Sprintf("starting %s: %s", s.Type, s.Name)
	if s.Version != "" {
		msg = fmt.Sprintf("%s %s", msg, s.Version)
	}
	if s.Revision != "" {
		msg = fmt.Sprintf("%s (%s)", msg, s.Revision[0:6])
	}
	log.Println(msg)
	for _, worker := range workers {
		go work(worker)
	}
	select {} // blocks the thread indefinitely
}
