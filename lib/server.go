package ghast

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
)

// Server represents an HTTP server that uses a Router to handle requests.
// It manages TCP listening, connection handling, request parsing, and routing across multiple routers.
// The Server includes a root router for direct route registration and supports sub-routers with path prefixes.
type Server struct {
	rootRouter  Router            // Default router for root-level routes
	routers     map[string]Router // Map of path prefixes to sub-routers (e.g., "/api" -> apiRouter)
	middlewares []Middleware      // Server-level middleware (applied to all routers)
	addr        string
	listener    net.Listener // TODO: Add listener for graceful shutdown
	isDone      bool         // TODO: Add shutdown signal

	// TODO: Add fields for future improvements:
	// - listener net.Listener (for graceful shutdown)
	// - done chan struct{} (shutdown signal)
	// - wg sync.WaitGroup (wait for goroutines)
	// - config ServerConfig (timeouts, max connections, etc.)
}

// ServerConfig holds configuration options for the server.
// TODO: Implement and use this for:
// - ReadTimeout / WriteTimeout
// - MaxConnections / MaxRequestBodySize
// - TLS/HTTPS support
// - Custom error handlers
// - Access logging configuration
type ServerConfig struct {
	// Placeholder for future configuration
}

type RouterPath struct {
	Path   string
	Router Router
}

// NewServer creates a new Server with a default root router and empty sub-router map.
func NewServer() *Server {
	return &Server{
		rootRouter: NewRouter(),
		routers:    make(map[string]Router),
	}
}

// NewRouter creates a new Router and registers it at the specified path prefix.
// Returns the new Router for chaining route registrations.
func (s *Server) NewRouter(prefix string) Router {
	if prefix == "" || prefix == "/" {
		log.Printf("Warning: Use the server's direct Http methods for root routes, not NewRouter")
		prefix = "/"
	}
	newRouter := NewRouter()
	s.registerRouter(RouterPath{Path: prefix, Router: newRouter})
	return newRouter
}

// registerRouter adds a Router to the server with an optional path prefix.
// If the path prefix is empty, it defaults to "/". Routers are matched based on longest prefix first.
// Returns the server instance for chaining.
func (s *Server) registerRouter(rp RouterPath) *Server {
	if _, exists := s.routers[rp.Path]; exists {
		log.Printf("Warning: Router for path %s already exists. Overwriting.", rp.Path)
	}
	if rp.Path == "" {
		rp.Path = "/" // Default to root if empty
	}
	s.routers[rp.Path] = rp.Router
	return s
}

// HTTP verb methods - delegate to root router for convenience

// Get registers a GET handler on the root router. Returns the server for chaining.
func (s *Server) Get(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Get(path, handler, middlewares...)
	return s
}

// Post registers a POST handler on the root router. Returns the server for chaining.
func (s *Server) Post(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Post(path, handler, middlewares...)
	return s
}

// Put registers a PUT handler on the root router. Returns the server for chaining.
func (s *Server) Put(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Put(path, handler, middlewares...)
	return s
}

// Delete registers a DELETE handler on the root router. Returns the server for chaining.
func (s *Server) Delete(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Delete(path, handler, middlewares...)
	return s
}

// Patch registers a PATCH handler on the root router. Returns the server for chaining.
func (s *Server) Patch(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Patch(path, handler, middlewares...)
	return s
}

// Head registers a HEAD handler on the root router. Returns the server for chaining.
func (s *Server) Head(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Head(path, handler, middlewares...)
	return s
}

// Options registers an OPTIONS handler on the root router. Returns the server for chaining.
func (s *Server) Options(path string, handler Handler, middlewares ...Middleware) *Server {
	s.rootRouter.Options(path, handler, middlewares...)
	return s
}

// Use adds middleware to the server. Server-level middleware is applied to all routes across all routers.
// Returns the server for chaining.
func (s *Server) Use(middleware Middleware) *Server {
	s.middlewares = append(s.middlewares, middleware)
	return s
}

// Listen starts the HTTP server on the given address (e.g., ":8080").
func (s *Server) Listen(addr string) error {
	s.addr = addr

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.listener = ln // Store listener for graceful shutdown support

	log.Printf("🌪️  Ghast server listening on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// TODO: Implement graceful shutdown check
			// if s.isDone() { return nil }
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// TODO: Add connection pooling / limiting
		// TODO: Add per-connection metrics and logging
		go s.handleConnection(conn)
	}
}

// Shutdown gracefully shuts down the server.
// TODO: Implement this to:
// - Signal all goroutines to stop accepting connections
// - Wait for existing requests to complete
// - Close the listener
// - Return after all connections are closed
func (s *Server) Shutdown() error {
	// Placeholder for graceful shutdown implementation
	return nil
}

// handleConnection processes a single TCP connection and its HTTP requests.
// It reads requests, parses them, and routes them to appropriate handlers.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		// Read HTTP request headers
		var headerLines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if line == "\r\n" {
				break
			}
			headerLines = append(headerLines, strings.TrimRight(line, "\r\n"))
		}

		if len(headerLines) == 0 {
			return
		}

		// Parse the request
		req, err := ParseRequest(strings.Join(headerLines, "\r\n"))
		if err != nil {
			// TODO: Send proper error response to client
			return
		}

		// Read request body if Content-Length is present
		if contentLength := req.Headers["Content-Length"]; contentLength != "" {
			var length int
			fmt.Sscanf(contentLength, "%d", &length)
			if length > 0 {
				// TODO: Add configurable max body size limit
				bodyBytes := make([]byte, length)
				reader.Read(bodyBytes)
				req.Body = string(bodyBytes)
			}
		}

		// Very basic client IP extraction for potential logging or middleware use. In a real implementation, this should be more robust and handle cases like proxies
		// and X-Forwarded-For headers.
		// See echo's ip.go for a more robust implementation: https://github.com/labstack/echo/blob/master/ip.go
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			req.ClientIP = conn.RemoteAddr().String() // Fallback to full address if splitting fails
		} else {
			req.ClientIP = host // Populate client IP for potential logging or middleware use
		}

		// Process the request through the router
		rw := NewResponseWriter(conn)

		// Determine the appropriate router to use based on request path.
		// Sort prefixes by length (longest first) to ensure more specific paths are matched first.
		var prefixes []string
		for prefix := range s.routers {
			prefixes = append(prefixes, prefix)
		}
		sort.Slice(prefixes, func(i, j int) bool {
			return len(prefixes[i]) > len(prefixes[j])
		})

		var matchedRouter Router = nil
		var matchedPrefix string
		for _, prefix := range prefixes {

			if strings.HasPrefix(req.Path, prefix) && (prefix == "/" || len(req.Path) == len(prefix) || req.Path[len(prefix)] == '/') {
				matchedRouter = s.routers[prefix]
				matchedPrefix = prefix
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

			routerWithMiddleware := ChainMiddleware(matchedRouter, s.middlewares)

			routerWithMiddleware.ServeHTTP(rw, req)

			req.Path = originalPath // Restore original path for logging or debugging
		} else {
			// Try the root router if no prefix matched
			routerWithMiddleware := ChainMiddleware(s.rootRouter, s.middlewares)
			routerWithMiddleware.ServeHTTP(rw, req)
			s.isDone = true // TODO: Implement proper shutdown signaling
			return
		}

		// TODO: Add request timeout handling
		// TODO: Add support for HTTP/1.1 100 Continue
	}
}

// Note: Request parsing (headers, query params, etc.) is delegated to ParseRequest()
// in request.go. This keeps server concerns (connection handling) separate from
// protocol concerns (request parsing).
