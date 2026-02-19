package gust

import (
	"log"
	"time"
)

// Middleware is a function type that wraps a Handler with additional functionality.
// It takes a Handler and returns a new Handler that wraps the original.
//
// Example:
//
//	loggingMiddleware := func(next gust.Handler) gust.Handler {
//	    return gust.HandlerFunc(func(w gust.ResponseWriter, r *gust.Request) {
//	        log.Printf("Request: %s %s", r.Method, r.Path)
//	        next.ServeHTTP(w, r)
//	    })
//	}
type Middleware func(Handler) Handler

// HandlerBuilder helps build handlers with middleware using a fluent API.
type HandlerBuilder struct {
	handler Handler
}

// NewHandlerBuilder creates a new HandlerBuilder with the given Handler.
func NewHandlerBuilder(handler Handler) *HandlerBuilder {
	return &HandlerBuilder{handler: handler}
}

// Use adds a middleware layer to the handler being built (chainable).
func (hb *HandlerBuilder) Use(middleware Middleware) *HandlerBuilder {
	hb.handler = middleware(hb.handler)
	return hb
}

// Build returns the final Handler with all middleware applied.
func (hb *HandlerBuilder) Build() Handler {
	return hb.handler
}

// ChainMiddleware applies a slice of middleware to a handler in order.
func ChainMiddleware(handler Handler, middlewares []Middleware) Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

// LoggingMiddleware logs HTTP request details and response timing.
func LoggingMiddleware(next Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		start := time.Now()
		log.Printf("[%s] %s %s", req.Method, req.Path, req.Version)
		defer func() {
			log.Printf("Completed in %v", time.Since(start))
		}()
		next.ServeHTTP(rw, req)
	})
}

// RecoveryMiddleware recovers from panics and returns a 500 error instead of crashing.
func RecoveryMiddleware(next Handler) Handler {
	return HandlerFunc(func(rw ResponseWriter, req *Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				rw.Status(500).SendString("Internal Server Error")
			}
		}()
		next.ServeHTTP(rw, req)
	})
}
