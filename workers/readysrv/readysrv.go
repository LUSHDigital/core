// Package readysrv is used to provide readiness checks for a service.
package readysrv

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const (
	// DefaultInterface is the port that we listen to the prometheus path on by default.
	DefaultInterface = "0.0.0.0:3674"

	// DefaultPath is the path where we expose prometheus by default.
	DefaultPath = "/ready"
)

// New creates a new default metrics server.
func New(checks Checks) *Server {
	return &Server{
		Interface: DefaultInterface,
		Path:      DefaultPath,
		Checks:    checks,
	}
}

// Server defines a readiness server.
type Server struct {
	Interface string
	Path      string
	Checks    Checks
}

// Run will start the metrics server.
func (s *Server) Run(ctx context.Context) error {
	addr := os.Getenv("READINESS_INTERFACE")
	if addr == "" {
		addr = DefaultInterface
	}
	path := os.Getenv("READINESS_PATH")
	if path == "" {
		path = DefaultPath
	}
	log.Printf("serving readiness checks server over http on http://%s%s", s.Interface, s.Path)
	http.Handle(path, CheckHandler(s.Checks))
	return http.ListenAndServe(addr, nil)

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
