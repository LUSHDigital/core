package main

import (
	"context"
	"net/http"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/keys"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	broker, cancel := keys.BrokerRSAPublicKey(context.Background(), keys.JWTPublicKeySources, 5*time.Second)
	defer cancel()

	r.Handle("/", auth.HandlerValidateJWT(broker, func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant("users.read") {
			http.Error(w, "access denied", http.StatusUnauthorized)
		}
	}))
}
