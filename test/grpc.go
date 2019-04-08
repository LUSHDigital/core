package test

import (
	"log"

	"google.golang.org/grpc"
)

// DialGRPC will connect to a grpc server on a specific port.
func DialGRPC(addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Panicf("did not connect: %v\n", err)
	}
	return conn
}
