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

// New sets up a new grpc server.
func New(options ...grpc.ServerOption) *Server {
	return &Server{grpc.NewServer(options...), net.JoinHostPort("", strconv.Itoa(Port)), time.Now}
}

// Server represents a collection of functions for starting and running an RPC server.
type Server struct {
	Connection *grpc.Server
	Addr       string
	Now        func() time.Time
}

// Run will start the gRPC server and listen for requests.
func (gs *Server) Run(ctx context.Context, out io.Writer) error {
	lis, err := net.Listen("tcp", gs.Addr)
	if err != nil {
		return err
	}

	hsrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gs.Connection, hsrv)

	fmt.Fprintf(out, "serving grpc on %s", gs.Addr)
	return gs.Connection.Serve(lis)
}
