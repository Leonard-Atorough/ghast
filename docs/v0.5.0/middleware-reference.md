---
title: "Middleware Reference"
description: "Comprehensive guide to built-in middleware in Ghast 0.5.0, including CORS, rate limiting, recovery, request ID, and response time tracking."
authors: ["Leonard"]
dateCreated: "2024-02-23"
dateUpdated: "2024-02-23"
---

# Middleware Reference

Ghast provides a set of built-in middleware modules in the `github.com/Leonard-Atorough/ghast/middleware` package. These middleware can be imported and applied at the application, router, or route level.

## Overview

All middleware in Ghast follow a consistent pattern:

```go
import "github.com/Leonard-Atorough/ghast/middleware"

app := ghast.New()

// Apply at app level
app.Use(middleware.CorsMiddleware(middleware.CorsOptions{...}))

// Apply at router level
router := ghast.NewRouter()
router.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{...}))

// Apply at route level
app.Get("/protected", handler, middleware.RecoveryMiddleware(middleware.Options{...}))
```

---

## CORS Middleware

Adds Cross-Origin Resource Sharing (CORS) headers to responses, enabling safe cross-origin requests from browsers.

**Import:**

```go
import "github.com/Leonard-Atorough/ghast/middleware"
```

**Function Signature:**

```go
func CorsMiddleware(options CorsOptions) ghast.Middleware
```

**Options:**

```go
type CorsOptions struct {
	AllowedOrigins    []string // Origins allowed to access the resource (default: "*")
	AllowedMethods    []string // HTTP methods allowed (default: "GET, POST, PUT, DELETE, OPTIONS")
	AllowedHeaders    []string // Request headers allowed (default: echoes requesting headers)
	PreflightMaxAge   int      // Seconds to cache preflight responses (default: not set)
	PreflightContinue bool     // Continue processing after preflight (default: false)
	Credentials       bool     // Allow credentials in cross-origin requests (default: false)
}
```

**Example: Allow all origins**

```go
app := ghast.New()
app.Use(middleware.CorsMiddleware(middleware.CorsOptions{}))
// Adds: Access-Control-Allow-Origin: *
```

**Example: Restrict origins and methods**

```go
app.Use(middleware.CorsMiddleware(middleware.CorsOptions{
	AllowedOrigins: []string{"https://example.com", "https://app.example.com"},
	AllowedMethods: []string{"GET", "POST"},
	AllowedHeaders: []string{"Content-Type", "Authorization"},
	PreflightMaxAge: 3600,     // Cache preflight for 1 hour
	Credentials: true,
}))
```

**Preflight Handling:**

The middleware automatically handles OPTIONS preflight requests from browsers. By default, it responds with HTTP 200 and terminates without calling the next handler. Set `PreflightContinue: true` to allow the request to proceed to your handler.

```go
app.Use(middleware.CorsMiddleware(middleware.CorsOptions{
	PreflightContinue: true, // OPTIONS preflight requests reach your handler
}))
```

**Headers Set:**

- `Access-Control-Allow-Origin`
- `Access-Control-Allow-Methods`
- `Access-Control-Allow-Headers`
- `Access-Control-Max-Age` (if `PreflightMaxAge > 0`)
- `Access-Control-Allow-Credentials` (if `Credentials: true`)

---

## Rate Limit Middleware

Implements per-IP rate limiting to prevent abuse and control traffic.

**Import:**

```go
import "github.com/Leonard-Atorough/ghast/middleware"
```

**Function Signature:**

```go
func RateLimitMiddleware(options RateLimitOptions) ghast.Middleware
```

**Options:**

```go
type RateLimitOptions struct {
	RequestsPerMinute int // Maximum requests allowed per minute per IP
}
```

**Example: 60 requests per minute per IP**

```go
app := ghast.New()
app.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
	RequestsPerMinute: 60,
}))
```

**Example: Strict rate limiting for an endpoint**

```go
api := ghast.NewRouter()

// All API routes are limited to 30 requests/minute
api.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
	RequestsPerMinute: 30,
}))

api.Get("/data", handlerFunc)
api.Post("/data", handlerFunc)

app.Route("/api", api)
```

**Behavior:**

- Tracks requests per client IP address
- Counts requests within a 60-second rolling window
- Returns HTTP 429 (Too Many Requests) when limit is exceeded
- Response body: `"Too Many Requests"`

**Note:** Rate limits are stored in-memory and reset when the process restarts. This is suitable for development and small deployments. For distributed systems, consider an external rate limiting service.

---

## Recovery Middleware

Catches panics in handlers and returns a graceful error response instead of crashing the server.

**Import:**

```go
import "github.com/Leonard-Atorough/ghast/middleware"
```

**Function Signature:**

```go
func RecoveryMiddleware(opts Options) ghast.Middleware
```

**Options:**

```go
type Options struct {
	Log    bool        // Log panic messages (default: true)
	Logger *log.Logger // Custom logger (default: standard logger)
}
```

**Example: Enable panic recovery with logging**

```go
app := ghast.New()
app.Use(middleware.RecoveryMiddleware(middleware.Options{
	Log: true,
}))
```

**Example: Use custom logger**

```go
customLogger := log.New(os.Stdout, "[APP] ", log.LstdFlags)

app.Use(middleware.RecoveryMiddleware(middleware.Options{
	Log:    true,
	Logger: customLogger,
}))
```

**Example: Disable logging**

```go
app.Use(middleware.RecoveryMiddleware(middleware.Options{
	Log: false,
}))
```

**Behavior:**

- Wraps handler execution in a defer-recover block
- Catches any panic thrown by the handler
- Logs panic message if enabled
- Responds with HTTP 500 and JSON error body: `{"error": "Internal Server Error"}`
- Continues running the server instead of crashing

**Best Practices:**

Apply recovery middleware early in your middleware chain to catch panics from all downstream handlers:

```go
app := ghast.New()

// Apply recovery first
app.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))

// Apply other middleware after
app.Use(corsMiddleware)
app.Use(loggingMiddleware)
```

---

## Request ID Middleware

Generates a unique ID for each request and includes it in the response header for request tracing.

**Import:**

```go
import "github.com/Leonard-Atorough/ghast/middleware"
```

**Function Signature:**

```go
func RequestIDMiddleware(opts RequestIDOptions) ghast.Middleware
```

**Options:**

```go
type RequestIDOptions struct {
	HeaderName string // Header name for storing request ID (default: "X-Request-ID")
}
```

**Example: Use default header**

```go
app := ghast.New()
app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))
// Adds header: X-Request-ID: <uuid>
```

**Example: Custom header name**

```go
app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{
	HeaderName: "Request-ID",
}))
// Adds header: Request-ID: <uuid>
```

**Behavior:**

- Generates a UUIDv4 for each request (falls back to timestamp if UUID generation fails)
- Sets the ID in the response header before calling the handler
- The ID is available to handlers via response headers
- Useful for correlating logs and request traces

**Example: Logging with request ID**

```go
app := ghast.New()
app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

app.Get("/data", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
	requestID := r.Headers["X-Request-ID"]
	log.Printf("[%s] Processing request...", requestID)
	w.JSON(200, map[string]string{"status": "ok"})
}))
```

---

## Response Time Middleware

Measures request processing time and includes it in the response header.

**Import:**

```go
import "github.com/Leonard-Atorough/ghast/middleware"
```

**Function Signature:**

```go
func ResponseTimeMiddleware(opts ResponseTimeOptions) ghast.Middleware
```

**Options:**

```go
type ResponseTimeOptions struct {
	HeaderName string // Header name for response time (default: "X-Response-Time")
	Suffix     string // Unit suffix: "ms", "s", "us", "ns" (default: "ms")
}
```

**Example: Default (milliseconds)**

```go
app := ghast.New()
app.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{}))
// Adds header: X-Response-Time: 42ms
```

**Example: Seconds**

```go
app.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
	Suffix: "s",
}))
// Adds header: X-Response-Time: 0s
```

**Example: Microseconds with custom header**

```go
app.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
	HeaderName: "Duration",
	Suffix:     "us",
}))
// Adds header: Duration: 42000us
```

**Supported Suffixes:**

- `"ms"` - Milliseconds (default)
- `"s"` - Seconds
- `"us"` - Microseconds
- `"ns"` - Nanoseconds

**Behavior:**

- Records time before calling the handler
- Calculates elapsed time after handler completes
- Converts to requested unit and appends suffix
- Sets response header before response is sent

---

## Middleware Composition

Combine multiple middleware to build a comprehensive middleware stack:

```go
import (
	"log"
	"github.com/Leonard-Atorough/ghast"
	"github.com/Leonard-Atorough/ghast/middleware"
)

app := ghast.New()

// Recovery must come first to catch panics from all handlers
app.Use(middleware.RecoveryMiddleware(middleware.Options{
	Log: true,
}))

// Add request tracing
app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

// Add performance monitoring
app.Use(middleware.ResponseTimeMiddleware(middleware.ResponseTimeOptions{
	Suffix: "ms",
}))

// Add rate limiting
app.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
	RequestsPerMinute: 100,
}))

// Add CORS for cross-origin requests
app.Use(middleware.CorsMiddleware(middleware.CorsOptions{
	AllowedOrigins: []string{"https://example.com"},
}))

app.Get("/api/users", getUsersHandler)
app.Listen(":8080")
```

**Response Headers with all middleware:**

```
HTTP/1.1 200 OK
X-Request-ID: a1b2c3d4-e5f6-4789-a0b1-c2d3e4f5g6h7
X-Response-Time: 12ms
Access-Control-Allow-Origin: https://example.com
```

---

## Router-Level Middleware

Apply middleware to specific route groups:

```go
app := ghast.New()

// Global middleware
app.Use(middleware.RecoveryMiddleware(middleware.Options{}))

// API routes with stricter rate limiting
apiRouter := ghast.NewRouter()
apiRouter.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
	RequestsPerMinute: 30,
}))
apiRouter.Get("/data", apiHandler)

// Public routes with relaxed limits
publicRouter := ghast.NewRouter()
publicRouter.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
	RequestsPerMinute: 100,
}))
publicRouter.Get("/health", healthHandler)

app.Route("/api", apiRouter)
app.Route("/public", publicRouter)
app.Listen(":8080")
```

---

## Route-Level Middleware

Apply middleware to individual routes:

```go
app := ghast.New()

// Recovery on all routes
app.Use(middleware.RecoveryMiddleware(middleware.Options{}))

// Protected route with request ID tracking
app.Get("/admin", adminHandler,
	middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

// Public endpoint with CORS
app.Get("/data", dataHandler,
	middleware.CorsMiddleware(middleware.CorsOptions{
		AllowedOrigins: []string{"*"},
	}))

app.Listen(":8080")
```
