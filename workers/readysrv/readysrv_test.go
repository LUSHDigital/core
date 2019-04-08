package readysrv_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/readysrv"
)

var (
	ctx context.Context
)

func Example() {
	srv := readysrv.New(readysrv.Checks{
		"google": readysrv.CheckerFunc(func() ([]string, bool) {
			if _, err := http.Get("https://google.com"); err != nil {
				return []string{err.Error()}, false
			}
			return []string{"google can be accessed"}, true
		}),
	})
	srv.Run(ctx, ioutil.Discard)
}

func TestCheckerFunc(t *testing.T) {
	yes := readysrv.CheckerFunc(func() ([]string, bool) { return []string{}, true })
	no := readysrv.CheckerFunc(func() ([]string, bool) { return []string{}, false })

	cases := []struct {
		name     string
		checks   readysrv.Checks
		expected int
	}{
		{
			name: "all ok",
			checks: readysrv.Checks{
				"a": yes,
				"b": yes,
			},
			expected: 200,
		},
		{
			name: "all ok",
			checks: readysrv.Checks{
				"a": yes,
				"b": no,
			},
			expected: 500,
		},
		{
			name: "all ok",
			checks: readysrv.Checks{
				"a": no,
				"b": no,
			},
			expected: 500,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			handler := readysrv.CheckHandler(c.checks)
			handler.ServeHTTP(rr, req)

			test.Equals(t, c.expected, rr.Code)
		})
	}
}
