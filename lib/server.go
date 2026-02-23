package ghast

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// server represents an HTTP server that uses a Router to handle requests.
// It manages TCP listening, connection handling, request parsing, and routing across multiple routers.
// The server includes a root router for direct route registration and supports sub-routers with path prefixes.
type server struct {
	addr     string
	listener net.Listener // TODO: Add listener for graceful shutdown
	isDone   bool         // TODO: Add shutdown signal

	config *serverConfig // TODO: Add server configuration options (timeouts, max connections, etc.)

	requestHandler RequestHandler // Core request handling function that processes incoming requests and routes them

	// TODO: Add fields for future improvements:
	// - listener net.Listener (for graceful shutdown)
	// - done chan struct{} (shutdown signal)
	// - wg sync.WaitGroup (wait for goroutines)
	// - config ServerConfig (timeouts, max connections, etc.)
}

// serverConfig holds configuration options for the server.
// TODO: Implement and use this for:
// - ReadTimeout / WriteTimeout
// - MaxConnections / MaxRequestBodySize
// - TLS/HTTPS support
// - Custom error handlers
// - Access logging configuration
type serverConfig struct {
	// Placeholder for future configuration
	Address                 string      // Server listen address (e.g., ":8080")
	HidePort                bool        // Option to hide port in logs or responses
	GracefulShutdownTimeout int         // Timeout in seconds for graceful shutdown
	OnShutdownError         func(error) // Optional callback for shutdown errors
}

type RequestHandler interface {
	handleRequest(ResponseWriter, *Request)
}

// newServer creates a new server with a default root router and empty sub-router map.
func newServer(handler RequestHandler, config *serverConfig) *server {
	if config == nil {
		config = &serverConfig{
			Address:                 ":8080",
			HidePort:                false,
			GracefulShutdownTimeout: 30,
			OnShutdownError: func(err error) {
				log.Printf("Error during shutdown: %v", err)
			},
		}
	}
	return &server{
		config:         config,
		requestHandler: handler,
	}
}

// Listen starts the HTTP server on the given address (e.g., ":8080").
func (s *server) Listen(addr string) error {
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
func (s *server) Shutdown() error {
	// Placeholder for graceful shutdown implementation
	return nil
}

// handleConnection processes a single TCP connection and handles HTTP requests.
// It focuses purely on TCP connection I/O: reading request headers/body, parsing, and extracting metadata.
func (s *server) handleConnection(conn net.Conn) {
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
		req, err := parseRequest(strings.Join(headerLines, "\r\n"))
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

		// Extract client IP for logging or middleware use.
		// Very basic implementation - in production, handle proxies and X-Forwarded-For headers.
		// See echo's ip.go for reference: https://github.com/labstack/echo/blob/master/ip.go
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			req.ClientIP = conn.RemoteAddr().String() // Fallback to full address if splitting fails
		} else {
			req.ClientIP = host // Populate client IP for logging or middleware use
		}

		// Create response writer and serve the request through routing logic
		rw := newResponseWriter(conn)
		s.requestHandler.handleRequest(rw, req)

		// Check for connection keep-alive
		if shouldKeepAlive(req) {
			continue
		} else {
			return
		}

		// TODO: Add request timeout handling
		// TODO: Add support for HTTP/1.1 100 Continue
	}
}

// shouldKeepAlive checks the Connection header to determine if the connection should be kept alive.
func shouldKeepAlive(req *Request) bool {
	connHeader := req.Headers["Connection"]
	return strings.EqualFold(connHeader, "keep-alive")
}

// Note: Request parsing (headers, query params, etc.) is delegated to ParseRequest()
// in request.go. This keeps server concerns (connection handling) separate from
// protocol concerns (request parsing).
