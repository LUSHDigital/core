// Package httpsrv provides a default set of configuration for hosting a http server in a service.
package httpsrv

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/LUSHDigital/core/response"
	"github.com/LUSHDigital/core/workers/internal/portfmt"

	"github.com/dustin/go-humanize"
)

const (
	// Port is the default HTTP port.
	Port = 80
)

var (
	// NotFoundHandler responds with the default a 404 response.
	NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		res := &response.Response{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		res.WriteTo(w)
	})

	// DefaultHTTPServer represents the default configuration for the http server
	DefaultHTTPServer = http.Server{
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}
)

// WrapperHandler returns the wrapper handler for the http server.
func WrapperHandler(now func() time.Time, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/healthz":
			HealthHandler(now)(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	})

}

// HealthHandler responds with service health.
func HealthHandler(now func() time.Time) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := now()

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		latency := time.Since(start).Nanoseconds() / (1 * 1000 * 1000) // Milliseconds
		res := &response.Response{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
			Data: &response.Data{
				Type: "health",
				Content: HealthResponse{
					Latency:       fmt.Sprintf("%d ms", latency),
					HeapInUse:     humanize.Bytes(mem.HeapInuse),
					HeapAlloc:     humanize.Bytes(mem.HeapAlloc),
					StackInUse:    humanize.Bytes(mem.StackInuse),
					NumGoRoutines: runtime.NumGoroutine(),
				},
			},
		}
		res.WriteTo(w)
	})

}

// New sets up a new HTTP server.
func New(handler http.Handler, servers ...*http.Server) *Server {
	var server *http.Server
	if len(servers) > 1 {
		server = servers[0]
	} else {
		server = &DefaultHTTPServer
	}
	return &Server{
		Server:  server,
		Handler: handler,
		Port:    Port,
		Now:     time.Now,
	}
}

// Server represents a collection of functions for starting and running an RPC server.
type Server struct {
	Server  *http.Server
	Handler http.Handler
	Port    int
	Now     func() time.Time
}

// Run will start the gRPC server and listen for requests.
func (gs *Server) Run(ctx context.Context, out io.Writer) error {
	port := portfmt.Port(gs.Port)
	lis, err := net.Listen("tcp", port.String())
	if err != nil {
		return err
	}

	if gs.Handler == nil {
		if gs.Server.Handler == nil {
			return fmt.Errorf("http server needs a handler")
		}
		gs.Handler = gs.Server.Handler
	}

	gs.Server.Handler = WrapperHandler(gs.Now, gs.Handler)

	gs.Port = lis.Addr().(*net.TCPAddr).Port
	fmt.Fprintf(out, "serving http on 0.0.0.0:%d", gs.Port)
	return gs.Server.Serve(lis)
}

// HealthResponse contains information about the service health.
type HealthResponse struct {
	Latency       string `json:"latency"`
	StackInUse    string `json:"stack_in_use"`
	HeapInUse     string `json:"heap_in_use"`
	HeapAlloc     string `json:"heap_alloc"`
	NumGoRoutines int    `json:"num_go_routines"`
}
