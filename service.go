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
	"strings"
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

// Runner represents the behaviour for running a service worker.
type Runner interface {
	// Run should run start processing the worker and be a blocking operation.
	Run(context.Context) error
}

// Halter represents the behaviour for stopping a service worker.
type Halter interface {
	// Halt should tell the worker to stop doing work.
	Halt(context.Context) error
}

// Worker represents the behaviour for a service worker.
type Worker interface {
	Runner
	Halter
}

// StartWorkers will start the given service workers and block block indefinitely, until interupted.
// The process with an appropriate status code.
// DEPRECATED: Use MustRun in favour of StartWorkers.
func (s *Service) StartWorkers(ctx context.Context, workers ...Worker) {
	log.Println("DEPRECATED: Use core.MustRun in favour of core.StartWorkers")
	s.MustRun(ctx, workers...)
}

// MustRun will start the given service workers and block block indefinitely, until interupted.
// The process with an appropriate status code.
func (s *Service) MustRun(ctx context.Context, workers ...Worker) {
	os.Exit(s.Run(ctx, workers...))
}

// Run will start the given service workers and block block indefinitely, until interupted.
func (s *Service) Run(ctx context.Context, workers ...Worker) int {
	const fail int = 1
	nWorkers := len(workers)
	if nWorkers < 1 {
		log.Println("need at least 1 service worker")
		return fail
	}
	if err := s.validate(); err != nil {
		log.Println(err)
		return fail
	}
	var (
		cancelled    <-chan int
		completed    <-chan int
		done, cancel func()
	)
	ctx, cancelled, cancel = ContextWithSignals(ctx)
	completed, cancelled, done = WaitWithTimeout(nWorkers, cancelled, s.grace())

	var run = func(ctx context.Context, worker Worker, done, cancel func()) {
		if err := worker.Run(ctx); err != nil {
			log.Println("service errored:", err)
			go cancel()
		}
		done()
	}
	var halt = func(ctx context.Context, worker Worker) {
		if err := worker.Halt(ctx); err != nil {
			log.Println("service halted:", err)
		}
	}

	log.Printf("starting %s: %s", s.Type, s.name())

	for _, worker := range workers {
		go run(ctx, worker, done, cancel)
	}
	for {
		select {
		case <-cancelled:
			for _, worker := range workers {
				go halt(ctx, worker)
			}
		case code := <-completed:
			message := "shutdown gracefully..."
			if code > 0 {
				message = "failed to shutdown gracefully: killing!"
			}
			log.Println(message)
			return code
		}
	}
}

func (s *Service) validate() error {
	if s.Name == "" || s.Type == "" {
		return fmt.Errorf("cannot start without a name or type")
	}
	return nil
}

func (s *Service) name() string {
	var w strings.Builder

	w.WriteString(s.Name)

	if len(s.Revision) > 5 {
		w.WriteString(" (" + s.Revision[0:6] + ")")
	}

	if s.Version != "" {
		w.WriteString(" " + s.Version)
	}
	return w.String()
}

func (s *Service) grace() time.Duration {
	grace := s.GracePeriod
	if grace == 0 {
		grace = time.Second * 5
	}
	return grace
}

// WaitWithTimeout will wait for a number of pieces of work has finished and send a message on the completed channel.
func WaitWithTimeout(delta int, cancelled <-chan int, timeout time.Duration) (<-chan int, <-chan int, func()) {
	completedC := make(chan int, 1)
	cancelledC := make(chan int, 1)
	wg := &sync.WaitGroup{}
	wg.Add(delta)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		completedC <- 0
	}(wg)
	go func() {
		select {
		case code := <-cancelled:
			cancelledC <- code
			time.Sleep(timeout)
			completedC <- code
		}
	}()
	return completedC, cancelledC, wg.Done
}

// ContextWithSignals creates a new instance of signal context.
func ContextWithSignals(ctx context.Context) (context.Context, <-chan int, context.CancelFunc) {
	var cancelCtx context.CancelFunc
	ctx, cancelCtx = context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	cancelled := make(chan int, 1)

	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	var cancel = func() {
		cancelCtx()
		cancelled <- 1
	}

	go func() {
		sig := <-sigs
		log.Printf("received signal: %s", sig)
		cancel()
	}()

	return ctx, cancelled, cancel
}
