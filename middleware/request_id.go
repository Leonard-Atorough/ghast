package middleware

import (
	"fmt"
	"time"
	uuid "github.com/google/uuid"
	ghast "ghast/lib"
)

const defaultRequestIDHeader = "X-Request-ID"

type RequestIDOptions struct {
	HeaderName string // The name of the header to set the request ID in (default: "X-Request-ID")
}

// RequestIDMiddleware is a middleware that generates a unique request ID for each incoming request and sets it in the response header.
func RequestIDMiddleware(opts RequestIDOptions) ghast.Middleware {
	headerName := defaultRequestIDHeader
	if opts.HeaderName != "" {
		headerName = opts.HeaderName
	}
	return func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			requestID := generateRequestID()
			w.SetHeader(headerName, requestID)
			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID generates a unique request ID using UUIDv4.
func generateRequestID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		// Fallback to a simple random string if UUID generation fails (should be rare)
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return id.String()
}