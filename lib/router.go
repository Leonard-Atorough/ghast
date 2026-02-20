package ghast

import (
	"regexp"
	"strings"
)

// Router interface defines the contract for HTTP routing and middleware management.
type Router interface {
	// Handle registers a handler for a specific HTTP method and path. The handler will be invoked when a request matches the method and path.
	Handle(method string, path string, handler Handler, middlewares ...Middleware)

	// Express-like convenience methods for common HTTP verbs
	Get(path string, handler Handler, middlewares ...Middleware) Router
	Post(path string, handler Handler, middlewares ...Middleware) Router
	Put(path string, handler Handler, middlewares ...Middleware) Router
	Delete(path string, handler Handler, middlewares ...Middleware) Router
	Patch(path string, handler Handler, middlewares ...Middleware) Router
	Head(path string, handler Handler, middlewares ...Middleware) Router
	Options(path string, handler Handler, middlewares ...Middleware) Router

	// ServeHTTP processes an incoming HTTP request by matching it to the appropriate handler based on the request's method and path.
	ServeHTTP(ResponseWriter, *Request)

	// Use adds a middleware function to the router. Middleware functions are applied to all handlers registered with the router, allowing you to add common
	// functionality (e.g., logging, authentication) across all routes without having to modify each handler individually.
	Use(middleware Middleware) Router

	// Listen starts the HTTP server on the given address (e.g., ":8080")
	Listen(addr string) error

	// Shutdown gracefully shuts down the server, allowing any in-flight requests to complete before closing the server.
	Shutdown() error
}

func (r *router) Shutdown() error {
	// Placeholder - actual implementation depends on embedding the server logic
	return nil
}

type router struct {
	routes      map[string]map[string]Handler // Nested map: first key is HTTP method (e.g., "GET", "POST"), second key is the path. Value is the Handler.
	middlewares []Middleware                  // Middleware applied to all routes.
	regexRoutes map[string]*pathRegex         // Regex patterns and params for routes with dynamic segments. Key is the path template.
}

// pathRegex stores compiled regex and parameter names for dynamic routes.
type pathRegex struct {
	regex  *regexp.Regexp // Compiled regex pattern for efficient matching.
	params []string       // Parameter names in order they appear in regex captures.
}

// NewRouter creates a new Router instance with empty routes and middleware.
func NewRouter() Router {
	return &router{
		routes:      make(map[string]map[string]Handler),
		regexRoutes: make(map[string]*pathRegex),
		middlewares: []Middleware{},
	}
}

// Handle registers a handler for a specific HTTP method and path. It also compiles regex patterns for dynamic routes and applies middleware.
func (r *router) Handle(method string, path string, handler Handler, middlewares ...Middleware) {
	// Extract route parameters and compile regex pattern for dynamic routes.
	params := extractRouteParams(path)
	pattern := pathToRegex(path)

	// Compile the regex pattern once during registration for efficient matching.
	compiledRegex := regexp.MustCompile(pattern)
	r.regexRoutes[path] = &pathRegex{
		regex:  compiledRegex,
		params: params,
	}

	// Collect middleware: global middleware + route-specific middleware.
	middlewareCollection := []Middleware{}
	middlewareCollection = append(middlewareCollection, r.middlewares...)
	middlewareCollection = append(middlewareCollection, middlewares...)

	// Apply middleware to the handler.
	handler = ChainMiddleware(handler, middlewareCollection)

	// Register the handler for the specified method and path.
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}
	r.routes[method][path] = handler
}

// Express-like convenience methods for HTTP verbs

// Routes HTTP GET requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Get(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("GET", path, handler, middlewares...)
	return r
}

// Routes HTTP POST requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Post(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("POST", path, handler, middlewares...)
	return r
}

// Routes HTTP PUT requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Put(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("PUT", path, handler, middlewares...)
	return r
}

// Routes HTTP DELETE requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Delete(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("DELETE", path, handler, middlewares...)
	return r
}

// Routes HTTP PATCH requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Patch(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("PATCH", path, handler, middlewares...)
	return r
}

// Routes HTTP HEAD requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Head(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("HEAD", path, handler, middlewares...)
	return r
}

// Routes HTTP OPTIONS requests to the specified path with the given handler. Returns the router for chaining.
func (r *router) Options(path string, handler Handler, middlewares ...Middleware) Router {
	r.Handle("OPTIONS", path, handler, middlewares...)
	return r
}

// ServeHTTP processes an incoming HTTP request by matching it to the appropriate handler.
func (r *router) ServeHTTP(w ResponseWriter, req *Request) {
	// First, try exact path match.
	if r.routes[req.Method] != nil {
		if handler, ok := r.routes[req.Method][req.Path]; ok {
			handler.ServeHTTP(w, req)
			return
		}
	}

	// Try matching against regex routes (paths with dynamic segments).
	for pathTemplate, route := range r.regexRoutes {
		matches := route.regex.FindStringSubmatch(req.Path)
		if len(matches) > 0 {
			// Extract captured parameters from regex matches.
			req.Params = make(map[string]string)
			for i, paramName := range route.params {
				if i+1 < len(matches) {
					req.Params[paramName] = matches[i+1]
				}
			}

			// Look up and invoke the handler for this route.
			if handler, ok := r.routes[req.Method][pathTemplate]; ok {
				handler.ServeHTTP(w, req)
				return
			}
		}
	}

	w.Status(404)
	w.Send([]byte("404 Not Found"))
}

// Use adds a middleware function to the router that applies to all routes.
func (r *router) Use(middleware Middleware) Router {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// Listen starts an HTTP server on the given address. This is a simplified implementation
// for demonstration purposes. See the main.go for the full TCP server setup.
// Deprecated: Server logic has been moved to server.go. This method is a placeholder and should not be used directly.
func (r *router) Listen(addr string) error {
	// Placeholder - actual implementation depends on embedding the server logic
	// For now, this is handled in main.go
	return nil
}

// extractRouteParams extracts parameter names from a path template.
// Example: "/users/:id/posts/:postId" returns ["id", "postId"].
func extractRouteParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params = append(params, paramName)
		}
	}
	return params
}

// pathToRegex converts a path template to a regex pattern.
// Example: "/users/:id/posts/:postId" returns "^/users/([^/]+)/posts/([^/]+)$".
func pathToRegex(path string) string {
	parts := strings.Split(path, "/")
	var regexParts []string
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			// Replace :paramName with a capture group for non-slash characters.
			regexParts = append(regexParts, "([^/]+)")
		} else {
			regexParts = append(regexParts, part)
		}
	}
	return "^" + strings.Join(regexParts, "/") + "$"
}
