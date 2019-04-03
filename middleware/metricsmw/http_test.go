package metricsmw_test

import (
	"net/http"

	"github.com/LUSHDigital/core/middleware/metricsmw"
)

func ExampleMeasureRequests() {
	http.Handle("/check", metricsmw.MeasureRequests(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
}
