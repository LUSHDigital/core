// Package httpsrv provides a default set of configuration for hosting a http server in a service.
package httpsrv

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/LUSHDigital/core/rest"
)

const (
	// Port is the default HTTP port.
	Port = 80
)

var (
	// NotFoundHandler responds with the default a 404 response.
	NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		res := &rest.Response{
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

	DefaultCORS = CORS{
		// Use broad defaults.
		// Origin: "*" is safe, see: https://fetch.spec.whatwg.org/#basic-safe-cors-protocol-setup
		AllowOrigin: "*",
		AllowHeaders: []string{
			"Authorization",
			"Origin",
			"Accept",
			"Content-Type",
			"X-Requested-With",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
	}
)

// WrapperHandler returns the wrapper handler for the http server.
func WrapperHandler(now func() time.Time, cors CORS, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.RequestURI == "/healthz":
			HealthHandler(now)(w, r)
		default:
			CORSHandler(cors, next)(w, r)
		}
	}
}

// HealthHandler responds with service health.
func HealthHandler(now func() time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := now()

		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)

		latency := time.Since(start).Nanoseconds() / (1 * 1000 * 1000) // Milliseconds
		res := &rest.Response{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
			Data: &rest.Data{
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
	}
}

// CORSHandler updates headers on CORS preflight requests and actual CORS requests.
func CORSHandler(c CORS, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			handlePreflight(c, w, r)
		} else {
			// This is an actual request. If CORS, it will require this header, if not, no harm done.
			w.Header().Set("Access-Control-Allow-Origin", c.AllowOrigin)
			next.ServeHTTP(w, r)
		}
	}
}

func handlePreflight(c CORS, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		log.Printf("CORS pre-flight aborted: %s!=OPTIONS", r.Method)
		return
	}
	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", c.AllowOrigin)
	headers.Set("Access-Control-Allow-Headers", strings.Join(c.AllowHeaders, ", "))
	headers.Set("Access-Control-Allow-Methods", strings.Join(c.AllowMethods, ", "))
	w.WriteHeader(http.StatusNoContent)
}

// NewDefault returns a http server
func NewDefault(handler http.Handler) *Server {
	server := &DefaultHTTPServer
	server.Handler = handler
	return New(server)
}

// New sets up a new HTTP server.
func New(server *http.Server) *Server {
	if server == nil {
		server = &DefaultHTTPServer
	}
	if server.WriteTimeout == 0 {
		server.WriteTimeout = DefaultHTTPServer.WriteTimeout
	}
	if server.ReadTimeout == 0 {
		server.ReadTimeout = DefaultHTTPServer.ReadTimeout
	}
	if server.IdleTimeout == 0 {
		server.IdleTimeout = DefaultHTTPServer.IdleTimeout
	}
	if server.ReadHeaderTimeout == 0 {
		server.ReadHeaderTimeout = DefaultHTTPServer.ReadHeaderTimeout
	}
	if server.Addr == "" {
		var addr string
		if addr = os.Getenv("HTTP_ADDR"); addr == "" {
			addr = net.JoinHostPort("0.0.0.0", strconv.Itoa(Port))
		}
		server.Addr = addr
	}
	return &Server{
		Server: server,
		Now:    time.Now,
		addrC:  make(chan *net.TCPAddr, 1),
		CORS:   DefaultCORS,
	}
}

type CORS struct {
	AllowOrigin  string
	AllowHeaders []string
	AllowMethods []string
}

// Server represents a collection of functions for starting and running an RPC server.
type Server struct {
	Server  *http.Server
	CORS    CORS
	Now     func() time.Time
	addrC   chan *net.TCPAddr
	tcpAddr *net.TCPAddr
}

// Run will start the gRPC server and listen for requests.
func (gs *Server) Run(_ context.Context) error {
	addr := gs.Server.Addr
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	gs.addrC <- lis.Addr().(*net.TCPAddr)

	if gs.Server.Handler == nil {
		return fmt.Errorf("http server needs a handler")
	}

	gs.Server.Handler = WrapperHandler(gs.Now, gs.CORS, gs.Server.Handler)
	log.Printf("serving http on http://%s", gs.Addr().String())
	return gs.Server.Serve(lis)
}

// Halt will attempt to gracefully shut down the server.
func (gs *Server) Halt(ctx context.Context) error {
	log.Printf("stopping serving http on http://%s...", gs.Addr().String())
	return gs.Server.Shutdown(ctx)
}

// Addr will block until you have received an address for your server.
func (gs *Server) Addr() *net.TCPAddr {
	if gs.tcpAddr != nil {
		return gs.tcpAddr
	}
	t := time.NewTimer(5 * time.Second)
	select {
	case addr := <-gs.addrC:
		gs.tcpAddr = addr
	case <-t.C:
		return &net.TCPAddr{}
	}
	return gs.tcpAddr
}

// HealthResponse contains information about the service health.
type HealthResponse struct {
	Latency       string `json:"latency"`
	StackInUse    string `json:"stack_in_use"`
	HeapInUse     string `json:"heap_in_use"`
	HeapAlloc     string `json:"heap_alloc"`
	NumGoRoutines int    `json:"num_go_routines"`
}
