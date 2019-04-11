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

func main() {
	core.SetupLogs()

	metrics := metricsrv.New()
	broker := keybroker.NewRSA(nil)
	readiness := readysrv.New(readysrv.Checks{
		"public_key": broker,
	})

	server := httpsrv.New(&http.Server{
		ReadTimeout: 10 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("hello world"))
		}),
	})

	ctx := context.Background()
	svc := &core.Service{
		Name:    "example",
		Type:    "service",
		Version: "1.0.0",
	}
	svc.StartWorkers(ctx,
		server,
		metrics,
		broker,
		readiness,
	)
}
