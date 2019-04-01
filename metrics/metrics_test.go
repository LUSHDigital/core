package metrics_test

import (
	"log"

	"github.com/LUSHDigital/core/metrics"
)

func ExampleListenAndServe() {
	go func() { log.Fatal(metrics.ListenAndServe()) }()
}
