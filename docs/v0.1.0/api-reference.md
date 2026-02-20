---
title: "API Reference"
description: "Comprehensive reference for Ghast's public API, including Server, Router, Handler, Middleware, Request, and Response types and methods. This is the go-to resource for understanding how to use Ghast's core functionality to build web applications."
authors: ["Leonardo"]
dateCreated: "2024-02-20"
dateUpdated: "2024-02-20"
---

# V 0.1.0 API Reference

> note: Ghast 0.1.0 is an early release focused on core routing and middleware features. The API may change in future versions as we add dynamic routing, error handling, and other enhancements.

> Ghast 0.1.0 needs Go 1.20 or later.

## Server

Represents an HTTP server that manages multiple routers and handles incoming TCP connections.

### Types

#### `Server`

```go
type Server struct {
    routers     map[string]Router  // Map of path prefixes to routers
    middlewares []Middleware       // Server-level middleware
    addr        string             // Server address
    listener    net.Listener       // TCP listener
    isDone      bool               // Shutdown signal
}
```

### Methods

#### `NewServer() *Server`

Creates a new Server instance with empty routers and middleware maps.

#### `AddRouter(rp RouterPath) *Server`

Registers a router at a specific path prefix. Returns the server for chaining.

**Props:**

- `rp.Path` (string): Path prefix (e.g., "/api"). Defaults to "/" if empty
- `rp.Router` (Router): Router instance to register

#### `Use(middleware Middleware) *Server`

Adds server-level middleware that applies to all routers. Returns the server for chaining.

**Props:**

- `middleware` (Middleware): Middleware function to apply globally

#### `Listen(addr string) error`

Starts the HTTP server on the given address. Blocks until an error occurs.

**Props:**

- `addr` (string): Listen address (e.g., ":8080")

#### `Shutdown() error`

Gracefully shuts down the server.

---

## Router

Interface for HTTP routing and middleware management. Handles path matching (exact and dynamic parameters) and middleware composition.

### Interface

```go
type Router interface {
    Handle(method string, path string, handler Handler)
    Get(path string, handler Handler) Router
    Post(path string, handler Handler) Router
    Put(path string, handler Handler) Router
    Delete(path string, handler Handler) Router
    Patch(path string, handler Handler) Router
    Head(path string, handler Handler) Router
    Options(path string, handler Handler) Router
    ServeHTTP(ResponseWriter, *Request)
    Use(middleware Middleware) Router
    UsePath(path string, middleware Middleware) Router
    Listen(addr string) error
    Shutdown() error
}
```

### Methods

#### `NewRouter() Router`

Creates a new Router instance with empty routes and middleware.

#### `Handle(method, path string, handler Handler)`

Registers a handler for a specific HTTP method and path. Compiles regex patterns for dynamic routes.

**Props:**

- `method` (string): HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- `path` (string): URL path, can include dynamic parameters prefixed with `:` (e.g., "/users/:id")
- `handler` (Handler): Handler to execute for matching requests

#### `Get(path string, handler Handler) Router`

Routes HTTP GET requests to the specified path.

#### `Post(path string, handler Handler) Router`

Routes HTTP POST requests to the specified path.

#### `Put(path string, handler Handler) Router`

Routes HTTP PUT requests to the specified path.

#### `Delete(path string, handler Handler) Router`

Routes HTTP DELETE requests to the specified path.

#### `Patch(path string, handler Handler) Router`

Routes HTTP PATCH requests to the specified path.

#### `Head(path string, handler Handler) Router`

Routes HTTP HEAD requests to the specified path.

#### `Options(path string, handler Handler) Router`

Routes HTTP OPTIONS requests to the specified path.

#### `ServeHTTP(w ResponseWriter, req *Request)`

Processes incoming HTTP requests. Matches requests to registered handlers using exact path matching first, then regex matching for dynamic routes.

#### `Use(middleware Middleware) Router`

Adds middleware that applies to all routes on this router.

**Props:**

- `middleware` (Middleware): Middleware function

#### `UsePath(path string, middleware Middleware) Router`

Adds middleware to a specific path only.

**Props:**

- `path` (string): The path to apply middleware to
- `middleware` (Middleware): Middleware function

#### `Listen(addr string) error` (Deprecated)

**Note:** Router.Listen is deprecated. Use Server.Listen instead.

#### `Shutdown() error`

Gracefully shuts down the router.

---

## Handler

Interface for processing HTTP requests and constructing responses.

### Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

### Types

#### `HandlerFunc`

Function type that implements the Handler interface. Allows ordinary functions to be used as handlers.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

#### `ServeHTTP(rw ResponseWriter, req *Request)`

Implements the Handler interface for HandlerFunc by calling the underlying function.

---

## Middleware

Functions that wrap handlers with additional functionality. Middleware can be applied globally (server or router level) or to specific paths.

### Types

#### `Middleware`

Function type that wraps a Handler and returns a new Handler.

```go
type Middleware func(Handler) Handler
```

### Utilities

#### `HandlerBuilder`

Fluent API for building handlers with middleware layers.

```go
type HandlerBuilder struct {
    handler Handler
}
```

##### `NewHandlerBuilder(handler Handler) *HandlerBuilder`

Creates a new HandlerBuilder with the given handler.

##### `Use(middleware Middleware) *HandlerBuilder`

Adds a middleware layer to the handler. Returns the builder for chaining.

##### `Build() Handler`

Returns the final handler with all applied middleware.

#### `ChainMiddleware(handler Handler, middlewares []Middleware) Handler`

Applies multiple middleware functions to a handler in order.

**Props:**

- `handler` (Handler): Base handler to wrap
- `middlewares` ([]Middleware): Slice of middleware functions to apply sequentially

### Built-in Middleware

#### `LoggingMiddleware(next Handler) Handler`

Logs HTTP request details (method and path) and response timing.

#### `RecoveryMiddleware(next Handler) Handler`

Recovers from panics and returns a 500 error instead of crashing.

---

## Request

Represents an HTTP request with parsed components.

### Type

#### `Request`

```go
type Request struct {
    Method  string            // HTTP method (GET, POST, etc.)
    Path    string            // URL path (without query string)
    Headers map[string]string // HTTP headers
    Body    string            // Request body as string
    Version string            // HTTP version (e.g., "HTTP/1.1")
    Params  map[string]string // Route parameters from dynamic paths
    Queries map[string]string // Query parameters
}
```

### Methods

#### `Query(key string) string`

Retrieves a query parameter value. Returns empty string if not found.

**Props:**

- `key` (string): Query parameter name

#### `Param(key string) string`

Retrieves a route parameter value from dynamic path segments (e.g., `:id`). Returns empty string if not found.

**Props:**

- `key` (string): Parameter name

#### `JSON(v any) error`

Unmarshals the request body as JSON into the provided value.

**Props:**

- `v` (any): Pointer to a value to unmarshal into

#### `GetHeader(key string) string`

Retrieves a header value with case-insensitive matching. Returns empty string if not found.

**Props:**

- `key` (string): Header name

#### `ContentType() string`

Returns the Content-Type header value.

#### `ParseRequest(rawRequest string) (*Request, error)`

Parses a raw HTTP request string into a Request struct. Extracts method, path, version, headers, and query parameters.

**Props:**

- `rawRequest` (string): Raw HTTP request text

---

## Response

Interface for constructing and sending HTTP responses.

### Interface

#### `ResponseWriter`

```go
type ResponseWriter interface {
    Header() map[string]string
    Status(statusCode int) ResponseWriter
    SetHeader(key, value string) ResponseWriter
    Send([]byte) (int, error)
    SendString(string) (int, error)
    JSON(statusCode int, data interface{}) error
    JSONPretty(statusCode int, data interface{}) error
    Plain(statusCode int, data string) error
    HTML(statusCode int, data string) error
}
```

### Methods

#### `NewResponseWriter(conn net.Conn) ResponseWriter`

Creates a new ResponseWriter for the given connection.

#### `Header() map[string]string`

Returns the response headers map for setting headers before writing the body.

#### `Status(statusCode int) ResponseWriter`

Sets the HTTP status code. Returns the writer for chaining.

**Props:**

- `statusCode` (int): HTTP status code (200, 404, 500, etc.)

#### `SetHeader(key, value string) ResponseWriter`

Sets a response header. Returns the writer for chaining.

**Props:**

- `key` (string): Header name
- `value` (string): Header value

#### `Send(data []byte) (int, error)`

Writes data to the response body.

**Props:**

- `data` ([]byte): Bytes to write

**Returns:**

- (int): Number of bytes written
- (error): Error if write failed

#### `SendString(s string) (int, error)`

Writes a string response.

**Props:**

- `s` (string): String to write

**Returns:**

- (int): Number of bytes written
- (error): Error if write failed

#### `JSON(statusCode int, data interface{}) error`

Marshals data as JSON and sends with application/json content-type.

**Props:**

- `statusCode` (int): HTTP status code
- `data` (interface{}): Data to marshal as JSON

#### `JSONPretty(statusCode int, data interface{}) error`

Marshals data as pretty-printed JSON and sends with application/json content-type.

**Props:**

- `statusCode` (int): HTTP status code
- `data` (interface{}): Data to marshal as JSON

#### `Plain(statusCode int, data string) error`

Sends a plain text response.

**Props:**

- `statusCode` (int): HTTP status code
- `data` (string): Text to send

#### `HTML(statusCode int, data string) error`

Sends an HTML response.

**Props:**

- `statusCode` (int): HTTP status code
- `data` (string): HTML string to send

---

## Error

Utilities for sending HTTP error responses.

### Types

#### `HTTPError`

```go
type HTTPError struct {
    StatusCode int    `json:"status"`
    Message    string `json:"error"`
}
```

### Functions

#### `Error(rw ResponseWriter, statusCode int, message string) error`

Sends an error response as JSON.

**Props:**

- `rw` (ResponseWriter): Response writer
- `statusCode` (int): HTTP status code
- `message` (string): Error message

#### `ErrorString(rw ResponseWriter, statusCode int, message string) error`

Sends an error response with a plain text body.

**Props:**

- `rw` (ResponseWriter): Response writer
- `statusCode` (int): HTTP status code
- `message` (string): Error message

---

## HTTP Methods

Constants for HTTP methods:

```go
const (
    GET     = "GET"
    POST    = "POST"
    PUT     = "PUT"
    DELETE  = "DELETE"
    HEAD    = "HEAD"
    OPTIONS = "OPTIONS"
    PATCH   = "PATCH"
)
```
