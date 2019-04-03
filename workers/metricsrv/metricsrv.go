// Package metricsrv provides a default set of configuration for hosting http prometheus metrics in a service.
package metricsrv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// DefaultInterface is the port that we listen to the prometheus path on by default.
	DefaultInterface = ":5117"

	// DefaultPath is the path where we expose prometheus by default.
	DefaultPath = "/metrics"
)

// New creates a new default metrics server.
func New() *Server {
	return &Server{DefaultInterface, DefaultPath}
}

// Server represents a prometheus metrics server.
type Server struct {
	Interface string
	Path      string
}

// Run will start the metrics server.
func (s *Server) Run(ctx context.Context, out io.Writer) error {
	addr := os.Getenv("PROMETHEUS_INTERFACE")
	if addr == "" {
		addr = DefaultInterface
	}
	path := os.Getenv("PROMETHEUS_PATH")
	if path == "" {
		path = DefaultPath
	}
	fmt.Fprintf(out, "running prometheus metrics server on %s%s", s.Interface, s.Path)
	http.Handle(path, promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
