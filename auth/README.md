# Auth
The `core/auth` package provides functions for services to issue and sign api consumer tokens.
It contains several middlewares for HTTP and GRPC to aid streamlining the authentication process.

## Examples

### HTTP middleware

```go
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
```

### gRPC middleware

```go
package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/keys"

	"google.golang.org/grpc"
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
```