package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/keys"
)

func main() {
	broker, cancel := keys.BrokerRSAPublicKey(context.Background(), keys.JWTPublicKeySources, 5*time.Second)
	defer cancel()

	srv := grpc.NewServer(
		grpc.StreamInterceptor(auth.StreamServerInterceptor(broker)),
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	srv.Serve(l)
}
