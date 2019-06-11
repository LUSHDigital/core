package tracingmw

import (
	"net/http"

	"github.com/LUSHDigital/uuid"

	"github.com/LUSHDigital/core/rest"
)

const (
	httpHeaderRequestIDKey = "X-Request-Id"
)

// EnsureRequestID will create a Request ID header if one is not found
func EnsureRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(httpHeaderRequestIDKey) == "" {
			requestID, err := uuid.NewV4()
			if err != nil {
				rest.InternalError(err).WriteTo(w)
				return
			}
			r.Header.Add(httpHeaderRequestIDKey, requestID.String())
		}
		ctxWithReqID := ContextWithRequestID(r.Context(), r.Header.Get(httpHeaderRequestIDKey))
		next.ServeHTTP(w, r.WithContext(ctxWithReqID))
	})
}
