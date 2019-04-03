package core_test

import (
	"context"
	"net/http"

	"github.com/LUSHDigital/core"
	"github.com/LUSHDigital/core/workers/grpcsrv"
	"github.com/LUSHDigital/core/workers/httpsrv"
	"github.com/LUSHDigital/core/workers/keybroker"
	"github.com/LUSHDigital/core/workers/metricsrv"
)

var (
	ctx     context.Context
	handler http.Handler
)

func ExampleNewService() {
	core.NewService()
}

func ExampleService_StartWorkers() {
	svc := core.NewService()
	svc.StartWorkers(ctx,
		grpcsrv.New(),
		httpsrv.New(handler, nil),
		metricsrv.New(),
		keybroker.NewRSA(),
	)
}
