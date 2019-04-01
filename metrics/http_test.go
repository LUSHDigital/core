package metrics_test

import (
	"net/http"

	"github.com/LUSHDigital/core/metrics"
	"github.com/gorilla/mux"
)

func ExampleMeasureRequestsMiddleware() {
	r := mux.NewRouter()
	r.Use(metrics.MeasureRequestsMiddleware)
}
func ExampleMeasureRequests() {
	http.Handle("/check", metrics.MeasureRequests(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
}
