// Package metrics is used to record and expose metrics for an application.
package metrics

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultInterface is the port that we listen to the prometheus path on by default.
const DefaultInterface = ":5117"

// DefaultPath is the path where we expose prometheus by default.
const DefaultPath = "/metrics"

// ListenAndServe will start the metrics server.
func ListenAndServe() error {
	addr := os.Getenv("PROMETHEUS_INTERFACE")
	if addr == "" {
		addr = DefaultInterface
	}
	path := os.Getenv("PROMETHEUS_PATH")
	if path == "" {
		path = DefaultPath
	}
	http.Handle(path, promhttp.Handler())
	return http.ListenAndServe(addr, nil)
}
