package main

import (
	"context"
	"net/http"
	"time"

	"github.com/LUSHDigital/core"
	"github.com/LUSHDigital/core/workers/httpsrv"
	"github.com/LUSHDigital/core/workers/keybroker"
	"github.com/LUSHDigital/core/workers/metricsrv"
	"github.com/LUSHDigital/core/workers/readysrv"
)

var service = &core.Service{
	Name:    "example",
	Type:    "service",
	Version: "1.0.0",
}

func main() {
	metrics := metricsrv.New(nil)
	broker := keybroker.NewPublicRSA(nil)
	readiness := readysrv.New(nil, readysrv.Checks{
		"public_rsa_key": broker,
	})

	server := httpsrv.New(&http.Server{
		ReadTimeout: 10 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("hello world"))
		}),
	})

	service.MustRun(context.Background(),
		server,
		metrics,
		broker,
		readiness,
	)
}
