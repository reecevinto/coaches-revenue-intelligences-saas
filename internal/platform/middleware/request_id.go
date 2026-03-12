package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID middleware injects a unique request ID into
// the request context and response headers for tracing.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// generate unique request ID
		requestID := uuid.New().String()

		// attach request ID to request context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// expose request ID in headers for clients
		w.Header().Set("X-Request-ID", requestID)

		// pass updated context forward
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
