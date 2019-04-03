package main

import (
	"context"
	"net/http"
	"time"

	"github.com/LUSHDigital/core"
	"github.com/LUSHDigital/core/workers/httpsrv"
	"github.com/LUSHDigital/core/workers/keybroker"
	"github.com/LUSHDigital/core/workers/metricsrv"
)

func main() {
	core.SetupLogs()

	service := &core.Service{
		Name:    "example",
		Type:    "service",
		Version: "1.0.0",
	}

	ctx := context.Background()

	httphandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})
	httpserver := &http.Server{
		ReadTimeout: 10 * time.Second,
	}

	service.StartWorkers(ctx,
		httpsrv.New(httphandler, httpserver),
		metricsrv.New(),
		keybroker.NewRSA(nil),
	)
}
