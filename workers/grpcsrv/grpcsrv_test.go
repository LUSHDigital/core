package grpcsrv_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/core/middleware/paginationmw"
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

func Example() {
	srv := grpcsrv.New(
		grpc.StreamInterceptor(paginationmw.StreamServerInterceptor),
		grpc.UnaryInterceptor(paginationmw.UnaryServerInterceptor),
	)
	srv.Port = 8080
	srv.Run(ctx, ioutil.Discard)
}
