// Package grpcsrv provides a default set of configuration for hosting a grpc server in a service.
package grpcsrv

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// Port is the default gRPC port used in examples.
	Port = 50051
)

// Config represents configuration for the GRPC server.
type Config struct {
	Addr string
}

// New sets up a new grpc server.
func New(config *Config, options ...grpc.ServerOption) *Server {
	if config == nil {
		config = &Config{
			Addr: net.JoinHostPort("", strconv.Itoa(Port)),
		}
	}
	return &Server{
		Connection: grpc.NewServer(options...),
		Now:        time.Now,
		addr:       config.Addr,
		addrC:      make(chan *net.TCPAddr, 1),
	}
}

// Server represents a collection of functions for starting and running an RPC server.
type Server struct {
	Connection *grpc.Server
	Now        func() time.Time
	addr       string
	addrC      chan *net.TCPAddr
	tcpAddr    *net.TCPAddr
}

// Run will start the gRPC server and listen for requests.
func (gs *Server) Run(ctx context.Context, out io.Writer) error {
	defer close(gs.addrC)
	lis, err := net.Listen("tcp", gs.addr)
	if err != nil {
		return err
	}
	gs.addrC <- lis.Addr().(*net.TCPAddr)

	hsrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gs.Connection, hsrv)

	fmt.Fprintf(out, "serving grpc on %s", gs.addr)
	return gs.Connection.Serve(lis)
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
		gs.tcpAddr = &net.TCPAddr{}
	}
	return gs.tcpAddr
}
