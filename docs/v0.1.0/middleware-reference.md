---
title: "Middleware Reference"
description: "Comprehensive reference for Ghast's built-in middleware components. Learn how to use recovery, request ID tracking, response timing, and other middleware for common tasks like error handling, tracing, and monitoring."
authors: ["Leonardo"]
dateCreated: "2024-02-20"
dateUpdated: "2024-02-20"
---

# V 0.1.0 Middleware Reference

Middleware functions wrap around handlers to perform additional processing before or after the handler executes. Ghast provides several built-in middlewares for common tasks like error recovery, request tracing, and performance monitoring. Middleware can be applied globally (server or router level) or to specific paths.

---

## Recovery Middleware

Catches panics that occur during handler execution and prevents the server from crashing. Returns a 500 Internal Server Error response to the client.

### Type

#### `Options`

Configuration for the Recovery Middleware.

```go
type Options struct {
    Log    bool         // Whether to log panic errors (default: true)
    Logger *log.Logger  // Optional custom logger (default: standard logger)
}
```

### Functions

#### `RecoveryMiddleware(opts Options) ghast.Middleware`

Creates and returns a Recovery Middleware with the specified options.

**Props:**

- `opts` (Options): Configuration options
  - `Log` (bool): Enable logging of recovered panics. Set to false to suppress logs
  - `Logger` (\*log.Logger): Optional custom logger instance. If nil, uses the standard logger

**Returns:** A Middleware function that wraps handlers with panic recovery

**Behavior:**

- Defers execution and catches any panic in the handler
- Logs the panic error if `Log` is true
- Sends a 500 Internal Server Error JSON response to the client
- Does not propagate the panic further up the stack

**Example:**

```go
// Default recovery with standard logger
router.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))

// Custom recovery with disabled logging
router.Use(middleware.RecoveryMiddleware(middleware.Options{Log: false}))

// With custom logger
customLogger := log.New(logFile, "ghast: ", log.LstdFlags)
router.Use(middleware.RecoveryMiddleware(middleware.Options{
    Log:    true,
    Logger: customLogger,
}))
```

---

## RequestID Middleware

Generates a unique request ID for each incoming request using UUIDv4 and sets it in the response header. Useful for request tracing, debugging, and log aggregation.

### Type

#### `RequestIDOptions`

Configuration for the RequestID Middleware.

```go
type RequestIDOptions struct {
    HeaderName string // The header name for the request ID (default: "X-Request-ID")
}
```

### Functions

#### `RequestIDMiddleware(opts RequestIDOptions) ghast.Middleware`

Creates and returns a RequestID Middleware with the specified options.

**Props:**

- `opts` (RequestIDOptions): Configuration options
  - `HeaderName` (string): HTTP header name for storing the request ID. Defaults to "X-Request-ID" if empty

**Returns:** A Middleware function that injects request IDs into responses

**Behavior:**

- Generates a unique UUIDv4 for each request
- Sets the UUID in the specified response header
- Falls back to nanosecond timestamp if UUID generation fails
- Executes before the handler via middleware chain

**Example:**

```go
// Default request ID header (X-Request-ID)
router.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

// Custom header name
router.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{
    HeaderName: "X-Trace-ID",
}))

// Router-level middleware
appRouter := ghast.NewRouter()
appRouter.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

// Server-level middleware (applies to all routers)
server := ghast.NewServer()
server.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))
```

**Response Header Example:**

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

---

## Response Time Middleware

Measures the time taken to process a request and sets it in the response header. Useful for performance monitoring and identifying slow endpoints.

### Type

#### `ResponseTimeOptions`

Configuration for the ResponseTime Middleware.

```go
type ResponseTimeOptions struct {
    HeaderName string // The header name for response time (default: "X-Response-Time")
    Suffix     string // Time unit suffix (default: "ms" for milliseconds)
}
```

### Functions

#### `ResponseTimeMiddleware(opts ResponseTimeOptions) ghast.Middleware`

Creates and returns a ResponseTime Middleware with the specified options.

**Props:**

- `opts` (ResponseTimeOptions): Configuration options
  - `HeaderName` (string): HTTP header name for storing the response time. Defaults to "X-Response-Time" if empty
  - `Suffix` (string): Time unit suffix appended to the value. Supports "ms" (milliseconds), "us" (microseconds), "ns" (nanoseconds). Defaults to "ms" if empty

**Returns:** A Middleware function that measures and injects response times into responses

**Behavior:**

- Records the start time before handler execution
- Measures elapsed time after handler completes
- Converts elapsed time to the specified unit (milliseconds by default)
- Sets the header with format: `{duration}{suffix}` (e.g., "45ms")
- Executes before and after the handler via middleware chain

**Example:**

```go
// Default response time in milliseconds (X-Response-Time)
router.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{}))

// Custom header name
router.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
    HeaderName: "X-Duration",
    Suffix:     "ms",
}))

// Microsecond precision
router.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
    HeaderName: "X-Processing-Time",
    Suffix:     "us",
}))

// No time unit suffix
router.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
    HeaderName: "X-Time",
    Suffix:     "",
}))
```

**Response Header Example:**

```
X-Response-Time: 45ms
```

---

## Middleware Composition

Middleware can be composed at multiple levels:

### Global Server Middleware

Applied to all routers and routes:

```go
server := ghast.NewServer()
server.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))
server.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{}))
server.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))
```

### Router-Level Middleware

Applied to all routes on a specific router:

```go
router := ghast.NewRouter()
router.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))
router.Get("/api/users", handler)
```

### Path-Specific Middleware

Applied only to a specific path:

```go
router := ghast.NewRouter()
router.UsePath("/admin", middleware.RecoveryMiddleware(middleware.Options{Log: true}))
```

### Handler-Specific Middleware

Pass middleware directly when registering a handler:

```go
router.Get("/users/:id", handler,
    middleware.RequestIDMiddleware(middleware.RequestIDOptions{}),
    middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{}),
)
```

### Order Matters

Middleware executes in the order it was added. The first middleware to be added will be the outermost wrapper:

```go
router.Use(firstMiddleware)  // Executes first
router.Use(secondMiddleware) // Executes second
// Handler executes third
// secondMiddleware completes
// firstMiddleware completes
```
