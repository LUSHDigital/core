package auth_test

import (
	"log"
	"net"
	"net/http"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/keys"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

var (
	broker keys.RSAPublicKeyCopierRenewer
)

func ExampleStreamServerInterceptor() {
	srv := grpc.NewServer(
		grpc.StreamInterceptor(auth.StreamServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleUnaryServerInterceptor() {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor(broker)),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleHandlerValidateJWT() {
	r := mux.NewRouter()
	r.Handle("/users", auth.HandlerValidateJWT(broker, func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant("users.read") {
			http.Error(w, "access denied", http.StatusUnauthorized)
		}
	}))
}
