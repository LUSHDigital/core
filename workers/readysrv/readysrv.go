// Package readysrv is used to provide readiness checks for a service.
package readysrv

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// DefaultAddr is the port that we listen to the prometheus path on by default.
	DefaultAddr = "0.0.0.0:3674"

	// DefaultPath is the path where we expose prometheus by default.
	DefaultPath = "/ready"
)

// Config represents the configuration for the metrics server.
type Config struct {
	Path   string
	Server *http.Server
}

// New creates a new default metrics server.
func New(config *Config, checks Checks) *Server {
	if config == nil {
		config = &Config{}
	}
	if path := os.Getenv("READINESS_PATH"); path != "" && config.Path == "" {
		config.Path = path
	}
	if config.Path == "" {
		config.Path = DefaultPath
	}
	if config.Server == nil {
		config.Server = &http.Server{}
	}
	if addr := os.Getenv("READINESS_ADDR"); addr != "" && config.Server.Addr == "" {
		config.Server.Addr = addr
	}
	if config.Server.Addr == "" {
		config.Server.Addr = DefaultAddr
	}
	config.Server.Handler = CheckHandler(checks)
	return &Server{
		Checks: checks,
		Server: config.Server,
		Path:   path.Join("/", config.Path),
		addrC:  make(chan *net.TCPAddr, 1),
	}
}

// Server defines a readiness server.
type Server struct {
	Path    string
	Checks  Checks
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
		return &net.TCPAddr{}
	}
	return s.tcpAddr
}

// Run will start the metrics server.
func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return err
	}
	s.addrC <- lis.Addr().(*net.TCPAddr)
	mux := http.NewServeMux()
	mux.Handle(s.Path, promhttp.Handler())
	s.Server.Handler = mux
	log.Printf("serving readiness checks server over http on http://%s%s", s.Addr(), s.Path)
	return s.Server.Serve(lis)
}

// Halt will attempt to gracefully shut down the server.
func (s *Server) Halt(ctx context.Context) error {
	log.Printf("stopping readiness checks server over http on http://%s...", s.Addr().String())
	return s.Server.Shutdown(ctx)
}

// CheckHandler provides a function for providing health checks over http.
func CheckHandler(checks Checks) http.HandlerFunc {
	type health struct {
		OK       bool     `json:"ok"`
		Messages []string `json:"messages"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var ready = true
		res := make(map[string]health)
		for name, check := range checks {
			messages, ok := check.Check()
			if !ok {
				ready = false
			}
			res[name] = health{
				OK:       ok,
				Messages: messages,
			}
		}

		w.Header().Add("Content-Type", "application/json")
		bts, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		code := http.StatusOK
		if !ready {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		w.Write(bts)
	}
}
