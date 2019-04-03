package auth_test

import (
	"log"
	"net"
	"net/http"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/middleware/authmw"

	"google.golang.org/grpc"
)

var (
	broker auth.RSAPublicKeyCopierRenewer
)

func ExampleStreamServerInterceptor() {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(authmw.StreamServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleUnaryServerInterceptor() {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authmw.UnaryServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleHandlerValidateJWT() {
	http.Handle("/users", authmw.HandlerValidateJWT(broker, func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant("users.read") {
			http.Error(w, "access denied", http.StatusUnauthorized)
		}
	}))
}
