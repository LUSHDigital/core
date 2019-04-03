package metricsmw

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// RequestCounter counts the number of requests for a combination of method, code and path.
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "A running count for HTTP requests",
		},
		[]string{"method", "code", "path"},
	)

	// RequestDurationHistogram measures the duration in nanoseconds for requests.
	RequestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_nanoseconds",
			Help: "Duration in nanoseconds of each request",
		},
		[]string{"method", "code", "path"},
	)

	// ResponseSizeHistogram measures the size in bytes for responses.
	ResponseSizeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_byte_size",
			Help: "Size in bytes of each response",
		},
		[]string{"method", "code", "path"},
	)

	// All represents a combination of all HTTP metric collectors.
	All = []prometheus.Collector{
		RequestCounter,
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
func Register() {
	prometheus.MustRegister(All...)
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
			rdh.Observe(float64(duration.Nanoseconds()))
		}

		if rc, err := RequestCounter.GetMetricWith(labels); err != nil {
			log.Printf("metrics: cannot get http request counter: %v\n", err)
		} else {
			rc.Inc()
		}
	}
}
