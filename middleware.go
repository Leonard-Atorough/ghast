package ghast

// Middleware is a function type that wraps a Handler with additional functionality.
// It takes a Handler and returns a new Handler that wraps the original.
//
// Example:
//
//	loggingMiddleware := func(next ghast.Handler) ghast.Handler {
//	    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
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
func chainMiddleware(handler Handler, middlewares []Middleware) Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
