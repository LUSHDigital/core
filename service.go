/*Package core consists of a set of packages which are used in writing micro-service
applications.

Each package defines conventional ways of handling common tasks,
as well as a suite of tests to verify their behaviour.
*/
package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	// GracePeriod represents the duration workers have to clean up before the process gets killed.
	GracePeriod time.Duration `json:"grace_period"`
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

// Worker represents the behaviour for running a service worker.
type Worker interface {
	// Run should be a blocking operation
	Run(context.Context) error
}

// StartWorkers will start the given service workers and block block indefinitely.
func (s *Service) StartWorkers(ctx context.Context, workers ...Worker) {
	if err := s.validate(); err != nil {
		log.Fatalln(err)
	}
	var (
		cancelled    <-chan int
		completed    <-chan int
		done, cancel func()
	)
	ctx, cancelled, cancel = ContextWithSignals(ctx)
	completed, done = WaitWithTimeout(len(workers), cancelled, s.grace())

	var work = func(ctx context.Context, worker Worker, done, cancel func()) {
		if err := worker.Run(ctx); err != nil {
			log.Println(err)
			cancel()
		}
		done()
	}

	log.Println(s.startmsg())
	for _, worker := range workers {
		go work(ctx, worker, done, cancel)
	}

	select {
	case code := <-completed:
		var message string
		switch code {
		case 0:
			message = "shutdown gracefully..."
		default:
			message = "failed to shutdown gracefully: killing!"
		}
		log.Println(message)
		os.Exit(code)
	}
}

func (s *Service) validate() error {
	if s.Name == "" || s.Type == "" {
		return fmt.Errorf("cannot start without a name or type")
	}
	return nil
}

func (s *Service) startmsg() string {
	msg := fmt.Sprintf("starting %s: %s", s.Type, s.Name)
	if s.Version != "" {
		msg = fmt.Sprintf("%s %s", msg, s.Version)
	}
	if s.Revision != "" {
		msg = fmt.Sprintf("%s (%s)", msg, s.Revision[0:6])
	}
	return msg
}

func (s *Service) grace() time.Duration {
	grace := s.GracePeriod
	if grace == 0 {
		grace = time.Second * 5
	}
	return grace
}

// WaitWithTimeout defines
func WaitWithTimeout(delta int, cancelled <-chan int, timeout time.Duration) (<-chan int, func()) {
	completed := make(chan int, 1)
	wg := &sync.WaitGroup{}
	wg.Add(delta)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		completed <- 0
	}(wg)
	go func() {
		select {
		case code := <-cancelled:
			time.Sleep(timeout)
			completed <- code
		}
	}()
	return completed, wg.Done
}

// ContextWithSignals creates a new instance of signal context.
func ContextWithSignals(ctx context.Context) (context.Context, <-chan int, context.CancelFunc) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	cancelled := make(chan int, 1)

	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	var cancelAndNotify = func() {
		cancel()
		cancelled <- 1
	}

	go func(cancel context.CancelFunc) {
		sig := <-sigs
		log.Printf("received signal: %s", sig)
		cancel()
	}(cancelAndNotify)

	return ctx, cancelled, cancelAndNotify
}
