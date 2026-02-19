package gust

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Server represents an HTTP server that uses a Router to handle requests.
// It manages TCP listening, connection handling, and request parsing.
type Server struct {
	router   Router // TODO: Add support for multiple routers. We can do this by defining a RouterGroup struct that can hold multiple routers and route requests to the appropriate router based on path prefixes or other criteria.
	addr     string
	listener net.Listener // TODO: Add listener for graceful shutdown
	isDone   bool         // TODO: Add shutdown signal

	// TODO: Add fields for future improvements:
	// - listener net.Listener (for graceful shutdown)
	// - done chan struct{} (shutdown signal)
	// - wg sync.WaitGroup (wait for goroutines)
	// - config ServerConfig (timeouts, max connections, etc.)
	// - Server level middlewares []Middleware (for global middleware)
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

// NewServer creates a new Server with the given Router.
func NewServer(router Router) *Server {
	return &Server{
		router: router,
	}
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

	log.Printf("🌪️  Gust server listening on %s", addr)

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

		// Process the request through the router
		rw := NewResponseWriter(conn)
		s.router.ServeHTTP(rw, req)

		// Check for connection keep-alive
		if connHeader := req.Headers["Connection"]; !strings.EqualFold(connHeader, "keep-alive") {
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
