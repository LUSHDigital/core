package pagination_test

import (
	"log"
	"net"

	"github.com/LUSHDigital/core/pagination"

	"google.golang.org/grpc"
)

func ExampleStreamServerInterceptor() {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(pagination.StreamServerInterceptor),
	)
	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleUnaryServerInterceptor() {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(pagination.UnaryServerInterceptor),
	)
	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}
