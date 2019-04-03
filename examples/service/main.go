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

var (
	service = &core.Service{
		Name:    "example",
		Type:    "service",
		Version: "1.0.0",
	}
)

func main() {
	ctx := context.Background()

	workers := []core.ServiceWorker{}

	httphandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})
	httpserver := &http.Server{
		ReadTimeout: 10 * time.Second,
	}

	workers = append(workers, httpsrv.New(httphandler, httpserver))
	workers = append(workers, metricsrv.New())
	workers = append(workers, keybroker.NewRSA(nil))

	service.StartWorkers(ctx, workers...)
}
