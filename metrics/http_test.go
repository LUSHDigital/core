package metrics_test

import (
	"net/http"

	"github.com/LUSHDigital/core/metrics"
)

func ExampleMeasureRequests() {
	http.Handle("/check", metrics.MeasureRequests(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
}
