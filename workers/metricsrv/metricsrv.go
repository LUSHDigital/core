// Package metricsrv provides a default set of configuration for hosting http prometheus metrics in a service.
package metricsrv

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// DefaultAddr is the port that we listen to the prometheus path on by default.
	DefaultAddr = "0.0.0.0:5117"

	// DefaultPath is the path where we expose prometheus by default.
	DefaultPath = "/metrics"
)

// Config represents the configuration for the metrics server.
type Config struct {
	Path   string
	Server *http.Server
}

// New creates a new default metrics server.
func New(config *Config) *Server {
	if config == nil {
		config = &Config{}
	}
	if path := os.Getenv("PROMETHEUS_PATH"); path != "" && config.Path == "" {
		config.Path = path
	}
	if config.Path == "" {
		config.Path = DefaultPath
	}
	if config.Server == nil {
		config.Server = &http.Server{}
	}
	if addr := os.Getenv("PROMETHEUS_ADDR"); addr != "" && config.Server.Addr == "" {
		config.Server.Addr = addr
	}
	if config.Server.Addr == "" {
		config.Server.Addr = DefaultAddr
	}
	return &Server{
		Path:   path.Join("/", config.Path),
		Server: config.Server,
		addrC:  make(chan *net.TCPAddr, 1),
	}
}

// Server represents a prometheus metrics server.
type Server struct {
	Path    string
	Server  *http.Server
	addrC   chan *net.TCPAddr
	tcpAddr *net.TCPAddr
}

// Addr will block until you have received an address for your server.
func (s *Server) Addr() *net.TCPAddr {
	if s.tcpAddr != nil {
		return s.tcpAddr
	}
	t := time.NewTimer(5 * time.Second)
	select {
	case addr := <-s.addrC:
		s.tcpAddr = addr
	case <-t.C:
		s.tcpAddr = &net.TCPAddr{}
	}
	return s.tcpAddr
}

// Run will start the metrics server.
func (s *Server) Run(ctx context.Context, out io.Writer) error {
	lis, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return err
	}
	s.addrC <- lis.Addr().(*net.TCPAddr)

	mux := http.NewServeMux()
	mux.Handle(s.Path, promhttp.Handler())
	s.Server.Handler = mux

	fmt.Fprintf(out, "serving prometheus metrics over http on %s%s", s.Addr().String(), s.Path)
	return s.Server.Serve(lis)
}
