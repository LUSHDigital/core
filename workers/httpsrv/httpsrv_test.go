package httpsrv_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/LUSHDigital/core/rest"
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
	}).Run(ctx)
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
	if err := rest.UnmarshalJSONResponse(rr.Body.Bytes(), res); err != nil {
		t.Fatal(err)
	}

	if !strings.HasSuffix(res.Latency, " ms") {
		t.Errorf("handler returned unexpected latencey suffix: got %v want %v", res.Latency, "ms")
	}
}

func TestCORSHandler(t *testing.T) {
	var (
		keyOrigin  = "Access-Control-Allow-Origin"
		keyHeaders = "Access-Control-Allow-Headers"
		keyMethods = "Access-Control-Allow-Methods"
	)
	for _, tt := range [...]struct {
		method     string
		URL        string
		cors       httpsrv.CORS
		wantStatus int
	}{
		{http.MethodOptions, "/", httpsrv.DefaultCORS, http.StatusNoContent},
		{http.MethodOptions, "/foo", httpsrv.DefaultCORS, http.StatusNoContent},
		{http.MethodGet, "/", httpsrv.DefaultCORS, http.StatusOK},
		// Custom CORS headers
		{
			http.MethodOptions,
			"/",
			httpsrv.CORS{
				AllowOrigin: "https://foo.bar.org",
				AllowHeaders: []string{
					"Origin",
					"X-Requested-With",
				},
				AllowMethods: []string{
					http.MethodPost,
					http.MethodGet,
					http.MethodOptions,
					http.MethodDelete,
				},
			},
			http.StatusNoContent,
		},
		{
			http.MethodPost,
			"/",
			httpsrv.CORS{
				AllowOrigin: "https://foo.bar.org",
				AllowHeaders: []string{
					"Origin",
					"X-Requested-With",
				},
				AllowMethods: []string{
					http.MethodPost,
					http.MethodGet,
					http.MethodOptions,
					http.MethodDelete,
				},
			},
			http.StatusOK,
		},
	} {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(tt.method, tt.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		httpsrv.CORSHandler(tt.cors, handler)(w, req)
		res := w.Result()
		wantHeaders := map[string]string{
			keyOrigin:  tt.cors.AllowOrigin,
			keyHeaders: strings.Join(tt.cors.AllowHeaders, ", "),
			keyMethods: strings.Join(tt.cors.AllowMethods, ", "),
		}
		// Preflight should abort non-OPTIONS methods without updating headers.
		if tt.method != http.MethodOptions {
			wantHeaders[keyHeaders] = ""
			wantHeaders[keyMethods] = ""
		}
		for k, want := range wantHeaders {
			if got := res.Header.Get(k); got != want {
				t.Errorf("incorrect CORS header %v, %s request to %s: got %v want %v", tt.method, tt.URL, k, got, want)
			}
		}
		if res.StatusCode != tt.wantStatus {
			t.Errorf("incorrect status code, %s request to %s: got %v want %v", tt.method, tt.URL, res.StatusCode, tt.wantStatus)
		}
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
		go srv.Run(ctx)
	}
	for _, srv := range servers {
		test.NotEquals(t, ":0", srv.Addr().String())
	}

}
