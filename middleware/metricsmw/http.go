package metricsmw

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestDurationHistogram measures the duration in seconds for requests.
	RequestDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration in seconds of each request",
		},
		[]string{"method", "code", "path"},
	)

	// ResponseSizeHistogram measures the size in bytes for responses.
	ResponseSizeHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_byte_size",
			Help: "Size in bytes of each response",
		},
		[]string{"method", "code", "path"},
	)

	// All represents a combination of all HTTP metric collectors.
	// TODO: Remove once we move to v1.0 since we no longer need to register the collectors manually.
	All = []prometheus.Collector{
		RequestDurationHistogram,
		ResponseSizeHistogram,
	}
)

type recorder struct {
	http.ResponseWriter
	status int
	length int
}

func (w *recorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *recorder) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// Register registers all the metric collectors with prometheus.
// DEPRECATED: metricsmw.Register() does not need to be called since registering of metrics now happens automatically.
// TODO: Remove once we move to v1.0 since we no longer need to register the collectors manually.
func Register() {
	log.Println("DEPRECATED: metricsmw.Register() does not need to be called since registering of metrics now happens automatically")
}

// MiddlewareFunc represents a middleware func for use with gorilla mux.
type MiddlewareFunc func(http.Handler) http.Handler

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}

// MeasureRequestsMiddleware wraps the measure requests handler in a gorilla mux middleware.
var MeasureRequestsMiddleware = MiddlewareFunc(MeasureRequestsHandler)

// MeasureRequestsHandler wraps the measure requests handler in a http handler.
func MeasureRequestsHandler(next http.Handler) http.Handler {
	switch nxt := next.(type) {
	case http.HandlerFunc:
		return MeasureRequests(nxt)
	default:
		log.Printf("could not create measure requests handler: invalid handler function\n")
		return next
	}
}

// MeasureRequests returns a middleware for collecting metrics on http requests.
func MeasureRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			now = time.Now()
			rec = &recorder{}
		)

		// Pass the request through to the handler.
		next.ServeHTTP(rec, r)

		duration := time.Now().Sub(now)
		labels := prometheus.Labels{
			"method": r.Method,
			"code":   strconv.Itoa(rec.status),
			"path":   r.RequestURI,
		}

		if rsh, err := ResponseSizeHistogram.GetMetricWith(labels); err != nil {
			log.Printf("metrics: cannot get http response size histogram: %v\n", err)
		} else {
			rsh.Observe(float64(rec.length))
		}

		if rdh, err := RequestDurationHistogram.GetMetricWith(labels); err != nil {
			log.Printf("metrics: cannot get http request duration histogram: %v\n", err)
		} else {
			rdh.Observe(duration.Seconds())
		}
	}
}
