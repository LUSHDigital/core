package metricsrv_test

import (
	"context"
	"io/ioutil"

	"github.com/LUSHDigital/core/workers/metricsrv"
)

var (
	ctx context.Context
)

func ExampleServer_Run() {
	srv := metricsrv.New()
	srv.Run(ctx, ioutil.Discard)
}
