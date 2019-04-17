package grpcsrv_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/LUSHDigital/core/middleware/paginationmw"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/grpcsrv"

	"google.golang.org/grpc"
)

var (
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		return
	})
	ctx context.Context
	now time.Time
)

func nowf() time.Time { return now }

func TestMain(m *testing.M) {
	ctx = context.Background()
	now = time.Now()
	os.Exit(m.Run())
}

func TestHealthCheck(t *testing.T) {
	server := grpcsrv.New(&grpcsrv.Config{
		Addr: "",
	})
	go server.Run(ctx, ioutil.Discard)
	addr := server.Addr()
	host := fmt.Sprintf("127.0.0.1:%d", addr.Port)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error(err)
	}
	client := grpc_health_v1.NewHealthClient(conn)
	res, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "",
	})
	if err != nil {
		t.Error(err)
	}
	test.Equals(t, "SERVING", res.Status.String())
}

func Example() {
	srv := grpcsrv.New(&grpcsrv.Config{
		Addr: ":8080",
	},
		grpc.StreamInterceptor(paginationmw.StreamServerInterceptor),
		grpc.UnaryInterceptor(paginationmw.UnaryServerInterceptor),
	)
	srv.Run(ctx, ioutil.Discard)
}

func TestServer_Addr(t *testing.T) {
	cases := 10
	servers := make([]*grpcsrv.Server, cases)
	for i := 0; i < cases; i++ {
		srv := grpcsrv.New(&grpcsrv.Config{
			Addr: ":",
		})
		servers[i] = srv
		go srv.Run(ctx, ioutil.Discard)
	}
	for _, srv := range servers {
		test.NotEquals(t, ":0", srv.Addr().String())
	}

}
