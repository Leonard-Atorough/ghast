package ghast

import (
	"sort"
	"strings"
)

const Version = "0.5.0"

type Ghast struct {
	config     *serverConfig
	rootRouter Router
	routers    []routeGroup
	server     *server

	middlewares []Middleware
}

// New creates and returns a new Server instance, ready for route registration and listening.
// This is the primary entry point for the Ghast framework.
// Example usage:
//
//	app := ghast.New()
//	app.Get("/hello", func(w ghast.ResponseWriter, r *ghast.Request) {
//	    w.SendString("Hello, World!")
//	})
//	app.Listen(":8080")
func New() *Ghast {
	return &Ghast{
		config:      &serverConfig{},
		rootRouter:  NewRouter(),
		routers:     []routeGroup{},
		middlewares: []Middleware{},
	}
}

// Router returns the root Router instance for direct route registration. This allows you to register routes directly on the main router without needing to create sub-routers or groups.
// In most cases, you can use the convenience methods on the Ghast struct (e.g., Get, Post) which internally register routes on the root router. However, if you need direct access to the Router interface for more advanced use cases (e.g., registering custom middleware, accessing route parameters), you can use this method to get the root router instance.
func (g *Ghast) Router() Router {
	return g.rootRouter
}

// Route allows you to mount an existing Router instance under a specific path prefix with optional middleware. This is useful for modularizing your application and reusing routers across different parts of your app or even across different projects.
//
// For example, you could create a router in a separate package and then mount it in your main application:
//
//	// In a separate package
//	func NewUserRouter() Router {
//	    r := ghast.NewRouter()
//	    r.Get("/users", handler)
//	    return r
//	}
//
//	// In your main application
//	userRouter := NewUserRouter()
//	app.Route("/api", userRouter, loggingMiddleware)
//
// The Route method takes a path prefix, a Router instance, and an optional list of middleware functions that will be applied to all routes within the mounted router. The mounted router's routes will be accessible under the specified path prefix.
func (g *Ghast) Route(prefix string, router Router, middlewares ...Middleware) *Ghast {
	rg := &routeGroup{
		prefix:      prefix,
		middlewares: middlewares,
		router:      router,
	}
	g.routers = append(g.routers, *rg)
	return g
}

func (g *Ghast) Use(middleware Middleware) *Ghast {
	g.middlewares = append(g.middlewares, middleware)
	return g
}

// Get registers a GET handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Get(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Get(path, handler, middlewares...)
	return g
}

// Post registers a POST handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Post(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Post(path, handler, middlewares...)
	return g
}

// Put registers a PUT handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Put(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Put(path, handler, middlewares...)
	return g
}

// Delete registers a DELETE handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Delete(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Delete(path, handler, middlewares...)
	return g
}

// Patch registers a PATCH handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Patch(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Patch(path, handler, middlewares...)
	return g
}

// Head registers a HEAD handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Head(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Head(path, handler, middlewares...)
	return g
}

// Options registers an OPTIONS handler on the root router at the entry point. Returns the server for chaining.
func (g *Ghast) Options(path string, handler Handler, middlewares ...Middleware) *Ghast {
	g.rootRouter.Options(path, handler, middlewares...)
	return g
}

func (g *Ghast) Listen(addr string) error {
	if g.server == nil {
		g.server = newServer(g, g.config)
	}
	return g.server.Listen(addr)
}

func (g *Ghast) handleRequest(rw ResponseWriter, req *Request) {
	var prefixes []string
	for _, rg := range g.routers {
		prefixes = append(prefixes, rg.prefix)
	}
	sort.Slice(prefixes, func(i, j int) bool {
		return len(prefixes[i]) > len(prefixes[j])
	})

	var matchedRouter Router = nil
	var matchedPrefix string
	for _, prefix := range prefixes {
		if strings.HasPrefix(req.Path, prefix) && (prefix == "/" || len(req.Path) == len(prefix) || req.Path[len(prefix)] == '/') {
			for _, rg := range g.routers {
				if rg.prefix == prefix {
					matchedRouter = rg.router
					matchedPrefix = prefix
					break
				}
			}
			break
		}
	}

	if matchedRouter != nil {
		// Strip the prefix from the path before passing to the router
		originalPath := req.Path
		if matchedPrefix != "/" {
			req.Path = strings.TrimPrefix(req.Path, matchedPrefix)
			if req.Path == "" {
				req.Path = "/"
			}
		}

		routerWithMiddleware := chainMiddleware(matchedRouter, g.middlewares)
		routerWithMiddleware.ServeHTTP(rw, req)

		req.Path = originalPath // Restore original path for logging or debugging
	}

	// Fall back to root router if no prefix matched
	routerWithMiddleware := chainMiddleware(g.rootRouter, g.middlewares)
	routerWithMiddleware.ServeHTTP(rw, req)
}
