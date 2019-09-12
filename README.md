[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/LUSHDigital/core/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/LUSHDigital/core)](https://goreportcard.com/report/github.com/LUSHDigital/core) [![Build Status](https://travis-ci.org/LUSHDigital/core.svg?branch=master)](https://travis-ci.org/LUSHDigital/core)
[![GoDoc](https://godoc.org/github.com/LUSHDigital/core?status.svg)](https://godoc.org/github.com/LUSHDigital/core)
# Core (Go)
A collection of packages for building a Go microservice on the LUSH platform.

## Quick start
Below there's an example for how to get running quickly with a service using the LUSHDigital core package.

```go
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/LUSHDigital/core"
	"github.com/LUSHDigital/core/workers/httpsrv"
	"github.com/LUSHDigital/core/workers/keybroker"
	"github.com/LUSHDigital/core/workers/metricsrv"
	"github.com/LUSHDigital/core/workers/readysrv"
)

var service = &core.Service{
	Name:    "example",
	Type:    "service",
	Version: "1.0.0",
}

func main() {
	core.SetupLogs()

	metrics := metricsrv.New(nil)
	broker := keybroker.NewPublicRSA(nil)
	readiness := readysrv.New(readysrv.Checks{
		"rsa_key": broker,
	})

	server := httpsrv.New(&http.Server{
		ReadTimeout: 10 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("hello world"))
		}),
	})

	service.StartWorkers(context.Background(),
		server,
		metrics,
		broker,
		readiness,
	)
}
```

## Documentation
Documentation and examples are provided in README files in each pacakge.

### Core Concepts
These packages contain functionality for the core concepts of our services.

- [core/auth](https://github.com/LUSHDigital/core/tree/master/auth#auth)
- [core/env](https://github.com/LUSHDigital/core/tree/master/env#env)
- [core/i18n](https://github.com/LUSHDigital/core/tree/master/i18n#internationalisation)
- [core/middleware](https://github.com/LUSHDigital/core/tree/master/middleware#middleware)
- [core/pagination](https://github.com/LUSHDigital/core/tree/master/pagination#pagination)
- [core/rest](https://github.com/LUSHDigital/core/tree/master/rest#rest)
- [core/test](https://github.com/LUSHDigital/core/tree/master/test#test)

### Middlewares
These packages contain convenient middlewares for transport protocols like HTTP REST and gRPC.

- [core/middleware/authmw](https://github.com/LUSHDigital/core/tree/master/middleware/authmw#auth-middleware)
- [core/middleware/i18nmw](https://github.com/LUSHDigital/core/tree/master/middleware/i18nmw#internationalisation-middleware)
- [core/middleware/metricsmw](https://github.com/LUSHDigital/core/tree/master/middleware/metricsmw#metrics-middleware)
- [core/middleware/paginationmw](https://github.com/LUSHDigital/core/tree/master/middleware/paginationmw#pagination-middleware)
- [core/middleware/tracingmw](https://github.com/LUSHDigital/core/tree/master/middleware/tracingmw#tracing-middleware)

### Service Workers
These packages contain convenient service workers things like network servers, background workers and message brokers.

- [core/workers/grpcsrv](https://github.com/LUSHDigital/core/tree/master/workers/grpcsrv#grpc-server)
- [core/workers/httpsrv](https://github.com/LUSHDigital/core/tree/master/workers/httpsrv#http-server)
- [core/workers/keybroker](https://github.com/LUSHDigital/core/tree/master/workers/keybroker#key-broker)
- [core/workers/metricsrv](https://github.com/LUSHDigital/core/tree/master/workers/metricsrv#metric-server)
- [core/workers/readysrv](https://github.com/LUSHDigital/core/tree/master/workers/readysrv#ready-server)

## Plugins
There are a few libraries that can be used in conjunction with the core library containing their own service workers, ready checks and/or middlewares.

- [core-redis](https://github.com/LUSHDigital/core-redis#core-redis): Libraries for connecting to, and working with a Redis store.
- [core-sql](https://github.com/LUSHDigital/core-sql#core-sql): Libraries for connecting to, and working with an SQL database.

## Tools
There are a few tools that can be used with projects that use the core libary.

- [LUSHDigital/jwtl](https://github.com/LUSHDigital/jwtl#jwtl-json-web-token-command-line-tool): A command line tool to help generating JWTs during development.
- [LUSHDigital/core-mage](https://github.com/LUSHDigital/core-mage): A library for the [mage build tool](https://magefile.org/) including convenient build targets used in conjunction with a projcet using this core library.

## Recommended Libraries
Some libraries have been designed to work together with the core library and some are even dependencies.
Consider using these if you need extended functionality for certain things.

- [LUSHDigital/scan](https://github.com/LUSHDigital/scan): Scan database/sql rows directly to structs, slices, and primitive types, originally forked from github.com/blockloop/scan
- [LUSHDigital/uuid](https://github.com/LUSHDigital/uuid): A UUID package originally forked from github.com/gofrs/uuid & github.com/satori/go.uuid
