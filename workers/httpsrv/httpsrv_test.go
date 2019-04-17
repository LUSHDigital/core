package httpsrv_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/LUSHDigital/core/response"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/httpsrv"
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
	go httpsrv.New(&http.Server{
		Handler:     handler,
		ReadTimeout: 1 * time.Second,
	}).Run(ctx, os.Stdout)
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := httpsrv.HealthHandler(nowf)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	res := &httpsrv.HealthResponse{}
	if err := response.UnmarshalJSONResponse(rr.Body.Bytes(), res); err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(res.Latency, " ms") {
		t.Errorf("handler returned unexpected latencey suffix: got %v want %v", res.Latency, "ms")
	}
}

func TestNotFoundHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	httpsrv.NotFoundHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestServer_Addr(t *testing.T) {
	cases := 10
	servers := make([]*httpsrv.Server, cases)
	for i := 0; i < cases; i++ {
		srv := httpsrv.New(&http.Server{
			Addr:    ":",
			Handler: handler,
		})
		servers[i] = srv
		go srv.Run(ctx, ioutil.Discard)
	}
	for _, srv := range servers {
		test.NotEquals(t, ":0", srv.Addr().String())
	}

}
