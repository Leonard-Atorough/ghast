package ghast

import "strings"

// Router interface defines the contract for HTTP routing and middleware management.
type Router interface {
	// Handle registers a handler for a specific HTTP method and path. The handler will be invoked when a request matches the method and path.
	Handle(method string, path string, handler Handler)

	// Express-like convenience methods for common HTTP verbs
	Get(path string, handler Handler) Router
	Post(path string, handler Handler) Router
	Put(path string, handler Handler) Router
	Delete(path string, handler Handler) Router
	Patch(path string, handler Handler) Router
	Head(path string, handler Handler) Router
	Options(path string, handler Handler) Router

	// ServeHTTP processes an incoming HTTP request by matching it to the appropriate handler based on the request's method and path.
	ServeHTTP(ResponseWriter, *Request)

	// Use adds a middleware function to the router. Middleware functions are applied to all handlers registered with the router, allowing you to add common
	// functionality (e.g., logging, authentication) across all routes without having to modify each handler individually.
	Use(middleware Middleware) Router

	// UsePath adds a middleware function to a specific path. Middleware functions added with this method will only be applied to handlers registered for the specified path, allowing you to add functionality that is specific to certain routes without affecting others.
	UsePath(path string, middleware Middleware) Router

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
	routes          map[string]map[string]Handler // This is a nested map where the first key is the HTTP method (e.g., "GET", "POST") and the second key is the path (e.g., "/hello"). The value is the Handler that should be invoked for that method and path.
	routeParams     map[string][]routeParam       // This map stores the route parameters for each registered path. The key is the path, and the value is a slice of routeParam structs that contain information about the parameters defined in that path (e.g., ":id" in "/users/:id").
	middlewares     []Middleware                  // This slice can be used to store middleware functions that should be applied to all handlers registered with the router. When a request is processed, the router can apply these middleware functions in order before invoking the final handler for the matched route.
	pathMiddlewares map[string][]Middleware       // This map stores middleware functions that should be applied to specific paths. The key is the path, and the value is a slice of middleware functions to be applied to that path.
}

type routeParam struct {
	name     string
	value    string
	position int
}

// NewRouter creates a new Router instance with empty routes and middleware.
func NewRouter() Router {
	return &router{
		routes:          make(map[string]map[string]Handler),
		routeParams:     make(map[string][]routeParam), // Initialize the routeParams map to an empty map to avoid nil pointer issues when adding route parameters later.
		middlewares:     []Middleware{},                // Initialize the middlewares slice to an empty slice to avoid nil pointer issues when adding middleware later.
		pathMiddlewares: make(map[string][]Middleware), // Initialize the pathMiddlewares map to an empty map to avoid nil pointer issues when adding path-specific middleware later.
	}
}

// Handle registers a handler for a specific HTTP method and path.
func (r *router) Handle(method string, path string, handler Handler) {

	// Extract route parameters from the path and store them in the routeParams map for later use when matching incoming requests to registered routes.
	// For example, if the path is "/users/:id", we would extract "id" as a route parameter and store it in the routeParams map with the key "/users/:id".
	// This allows us to later match incoming requests to this route and extract the value of the "id" parameter from the request path.
	params := extractRouteParams(path)
	r.routeParams[path] = params

	middlwareCollection := []Middleware{}

	middlwareCollection = append(middlwareCollection, r.middlewares...)
	if pathMiddlewares, ok := r.pathMiddlewares[path]; ok {
		for _, middleware := range pathMiddlewares {
			middlwareCollection = append(middlwareCollection, middleware)
		}
	}
	// Build the final handler with all middleware applied before registering it in the routes map.
	handler = ChainMiddleware(handler, middlwareCollection)

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}

	r.routes[method][path] = handler
}

// Express-like convenience methods for HTTP verbs
func (r *router) Get(path string, handler Handler) Router {
	r.Handle("GET", path, handler)
	return r
}

func (r *router) Post(path string, handler Handler) Router {
	r.Handle("POST", path, handler)
	return r
}

func (r *router) Put(path string, handler Handler) Router {
	r.Handle("PUT", path, handler)
	return r
}

func (r *router) Delete(path string, handler Handler) Router {
	r.Handle("DELETE", path, handler)
	return r
}

func (r *router) Patch(path string, handler Handler) Router {
	r.Handle("PATCH", path, handler)
	return r
}

func (r *router) Head(path string, handler Handler) Router {
	r.Handle("HEAD", path, handler)
	return r
}

func (r *router) Options(path string, handler Handler) Router {
	r.Handle("OPTIONS", path, handler)
	return r
}

// ServeHTTP processes an incoming HTTP request by matching it to the appropriate handler.
func (r *router) ServeHTTP(w ResponseWriter, req *Request) {

	// TODO: Include logic to handle route parameters when matching incoming requests to registered routes.
	// This will likely involve iterating through the registered routes for the request's method and checking
	// if the request path matches any of the registered paths, including those with route parameters (e.g., "/users/:id").
	// If a match is found, we would need to extract the values of the route parameters from the request path and populate the Params
	// field of the Request struct before invoking the handler.

	if r.routes[req.Method] != nil {
		if handler, ok := r.routes[req.Method][req.Path]; ok {
			handler.ServeHTTP(w, req)
			return
		}
	}
	// If no handler is found, respond with 404 Not Found
	w.Status(404)
	w.Send([]byte("404 Not Found"))
}

// Use adds a middleware function to the router that applies to all routes.
func (r *router) Use(middleware Middleware) Router {
	r.middlewares = append(r.middlewares, middleware)
	return r
}

// UsePath adds a middleware function to a specific path.
func (r *router) UsePath(path string, middleware Middleware) Router {
	if r.pathMiddlewares == nil {
		r.pathMiddlewares = make(map[string][]Middleware)
	}
	r.pathMiddlewares[path] = append(r.pathMiddlewares[path], middleware)
	return r
}

// Listen starts an HTTP server on the given address. This is a simplified implementation
// for demonstration purposes. See the main.go for the full TCP server setup.
func (r *router) Listen(addr string) error {
	// Placeholder - actual implementation depends on embedding the server logic
	// For now, this is handled in main.go
	return nil
}

// extractRouteParams is a helper function to extract route parameters from a path.
// For example, if the path is "/users/:id", it would extract "id" as a parameter.
func extractRouteParams(path string) []routeParam {
	// returns a routeParam struct array

	params := []routeParam{}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params = append(params, routeParam{name: paramName, position: i})

		}
	}
	return params
}
