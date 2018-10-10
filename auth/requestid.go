package auth

import (
	"github.com/LUSHDigital/microservice-core-golang/response"
	"github.com/gofrs/uuid"
	"net/http"
)

// EnsureRequestID will create a Request ID header if one is not found.
// It will then place the request ID into the request's context.
func EnsureRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("x-request-id") == "" {
			requestID, err := uuid.NewV4()
			if err != nil {
				response.InternalError(err).WriteTo(w)
				return
			}
			r.Header.Add("x-request-id", requestID.String())
		}

		ctxWithReqID := NewContextWithRequestID(r.Context(), r)

		next.ServeHTTP(w, r.WithContext(ctxWithReqID))
	})
}
