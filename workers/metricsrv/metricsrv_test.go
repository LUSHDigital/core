package metricsrv_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/LUSHDigital/core/test"

	"github.com/LUSHDigital/core/workers/metricsrv"
)

var (
	ctx context.Context
)

func ExampleServer_Run() {
	srv := metricsrv.New(&metricsrv.Config{
		Server: &http.Server{
			Addr: "0.0.0.0:5117",
		},
		Path: "/metrics",
	})
	srv.Run(ctx, ioutil.Discard)
}

func TestNew(t *testing.T) {
	os.Setenv("PROMETHEUS_ADDR", "0.0.0.0:1111")
	os.Setenv("PROMETHEUS_PATH", "/testmetrics")
	srv := metricsrv.New(nil)
	test.Equals(t, "0.0.0.0:1111", srv.Server.Addr)
	test.Equals(t, "/testmetrics", srv.Path)
	srv = metricsrv.New(&metricsrv.Config{})
	test.Equals(t, "0.0.0.0:1111", srv.Server.Addr)
	test.Equals(t, "/testmetrics", srv.Path)
	srv = metricsrv.New(&metricsrv.Config{
		Path: "zero",
		Server: &http.Server{
			Addr: "0.0.0.0:2222",
		},
	})
	test.Equals(t, "0.0.0.0:2222", srv.Server.Addr)
	test.Equals(t, "/zero", srv.Path)
}

func TestServer_Addr(t *testing.T) {
	cases := 100
	servers := make([]*metricsrv.Server, cases)
	for i := 0; i < cases; i++ {
		srv := metricsrv.New(&metricsrv.Config{
			Server: &http.Server{
				Addr: ":",
			},
		})
		servers[i] = srv
		go srv.Run(ctx, ioutil.Discard)
	}
	for _, srv := range servers {
		test.NotEquals(t, ":0", srv.Addr().String())
	}

}
