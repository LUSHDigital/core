package middleware_test

import (
	"github.com/LUSHDigital/core/middleware"
	"github.com/LUSHDigital/core/middleware/metricsmw"
	"github.com/LUSHDigital/core/middleware/paginationmw"
	"github.com/LUSHDigital/core/middleware/tracingmw"
	"github.com/LUSHDigital/core/workers/grpcsrv"
)

var (
	server *grpcsrv.Server
)

func Example() {
	server = grpcsrv.New(nil, middleware.WithUnaryServerChain(
		metricsmw.UnaryServerInterceptor,
		tracingmw.UnaryServerInterceptor,
		paginationmw.UnaryServerInterceptor,
	))
}
