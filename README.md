# Ghast

An Express.js-inspired HTTP framework for Go. Built from first principles to learn networking, concurrency, and elegant API design.

## Why Ghast?

**Ghast** bridges two worlds:

- **For Node.js developers:** Familiar Express-like API for building HTTP services in Go
- **For learning:** Understand HTTP, TCP, and concurrency by reading clean, documented code

This is an educational framework that prioritizes clarity and approachability. While not production-ready (use [Echo](https://echo.labstack.com/), [Gin](https://gin-gonic.com/), or stdlib for production), Ghast is perfect for:

- Learning Go networking and concurrency
- Building small projects or APIs
- Understanding HTTP server design patterns
- **Portfolio projects** showing thoughtful architecture

### What does ghast mean?

A "ghast" is a ghostly creature from folklore, often depicted as a haunting presence. In the context of this project, it symbolizes the "ghost" of Express.js in the Go ecosystem—a familiar API that allows developers to build web servers with ease while learning the underlying mechanics. Just as a ghast can be both eerie and fascinating, this framework aims to be an intriguing blend of simplicity and depth for those looking to explore Go's capabilities in web development.

`It's also a portmanteau of "Go", "HTTP", and "Fast" - emphasizing the goal of creating a fast, efficient HTTP framework in Go.`

## Quick Start

### Installation

```bash
go get github.com/Leonard-Atorough/ghast
```

### Hello World

```go
package main

import (
	"log"
	"github.com/Leonard-Atorough/ghast"
)

func main() {
    app := ghast.New()

    app.Get("/hello", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, map[string]string{
            "message": "Hello, World!",
        })
    }))

    log.Println("Server running on http://localhost:8080")
    app.Listen(":8080")
}
```

Run it:

```bash
go run main.go
curl http://localhost:8080/hello
```

Output:

```json
{ "message": "Hello, World!" }
```

## Features ✨

### Express-Like API

```go
app.Get("/users/:id", getUserHandler)
app.Post("/users", createUserHandler)
app.Put("/users/:id", updateUserHandler)
app.Delete("/users/:id", deleteUserHandler)
app.Patch("/users/:id", patchUserHandler)
```

### Chainable Response Methods

```go
w.Status(201).
  SetHeader("X-Custom", "value").
  JSON(map[string]interface{}{
    "id": 123,
    "name": "Alice",
  })
```

### Request Helpers

```go
// Query parameters
name := r.Query("name")      // ?name=alice
// Or alias:
name := r.Param("name")

// JSON parsing
var user User
r.JSON(&user)

// Headers
token := r.GetHeader("Authorization")
contentType := r.ContentType()
```

### Middleware Support

Ghast supports middleware composition at multiple levels:

```go
app := ghast.New()

// App-level middleware - applies to all routes
app.Use(recoveryMiddleware)
app.Use(requestIDMiddleware)

// Router-level middleware - applies to all routes on that router
userRouter := ghast.NewRouter()
userRouter.Use(authMiddleware)
app.Route("/users", userRouter)

// Route-specific middleware - applies only to that route
app.Get("/admin", adminHandler, authMiddleware)

// Custom middleware
customMiddleware := func(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        // Before handler
        log.Println("Request:", r.Path)

        // Call handler
        next.ServeHTTP(w, r)

        // After handler
        log.Println("Response sent")
    })
}
app.Use(customMiddleware)
```

**Built-in Middleware** (see [Middleware Reference](docs/v0.5.0/middleware-reference.md)):

- **Recovery:** Catches panics and returns graceful error responses
- **Request ID:** Generates unique IDs for request tracing
- **Response Time:** Measures and reports request processing time
- **CORS:** Enables cross-origin resource sharing
- **Rate Limit:** Implements per-IP rate limiting

## Project Structure

```
ghast/
├── ghast.go               # Main app struct and public API
├── router.go              # HTTP routing engine
├── handler.go             # Handler interface and types
├── request.go             # Request parsing & helpers
├── response.go            # Response writing & helpers
├── middleware.go          # Middleware system
├── error.go               # Error handling
├── group.go               # Route grouping
├── server.go              # Server and listener
├── http_status.go         # HTTP status constants
│
├── middleware/            # Built-in middleware packages
│   ├── cors.go           # Cross-Origin Resource Sharing
│   ├── rate_limit.go     # Per-IP rate limiting
│   ├── recovery.go       # Panic recovery
│   ├── request_id.go     # Request ID generation
│   └── response_time.go  # Response time tracking
│
├── examples/              # Example applications
│   ├── 1_basic_example.go
│   ├── 2_grouped_routes_example.go
│   └── 3_custom_middleware_example.go
│
├── docs/                  # Documentation
│   └── v0.5.0/
│       ├── guide.md                    # Routing and middleware guide
│       ├── api-reference.md            # Complete API documentation
│       └── middleware-reference.md     # Middleware reference
│
└── README.md
```

### Key Design Decisions

1. **Interfaces over Implementations** - `Handler` and `ResponseWriter` are interfaces, promoting composability
2. **Chainable Builders** - Response methods return `ResponseWriter` for method chaining (Express-style)
3. **Middleware as Functions** - Simple functional middleware that wrap handlers
4. **Minimal Dependencies** - Only standard library (no external packages needed)

## Usage Examples

### Basic Routing & JSON APIs

```go
import "github.com/Leonard-Atorough/ghast"

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

app := ghast.New()

// GET /users/:id
app.Get("/users/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
    w.JSON(200, user)
}))

// POST /users
app.Post("/users", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
    var user User
    if err := r.JSON(&user); err != nil {
        ghast.Error(w, 400, "Invalid request body")
        return
    }
    user.ID = 2
    w.Status(201).JSON(user)
}))

// GET /health (query parameter example)
app.Get("/health", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
    verbose := r.Query("verbose") == "true"
    status := "healthy"
    w.JSON(200, map[string]interface{}{
        "status": status,
        "verbose": verbose,
    })
}))

app.Listen(":8080")
```

### Organized Routing with Sub-Routers

```go
app := ghast.New()

// Create API v1 router
apiv1 := ghast.NewRouter()
apiv1.Get("/status", statusHandler)
apiv1.Get("/users", listUsersHandler)
apiv1.Post("/users", createUserHandler)

// Create API v2 router with different handlers
apiv2 := ghast.NewRouter()
apiv2.Get("/status", statusHandlerV2)
apiv2.Get("/users", listUsersHandlerV2)

// Mount routers at different prefixes
app.Route("/api/v1", apiv1)
app.Route("/api/v2", apiv2)

app.Listen(":8080")
```

### Middleware

```go
import (
    "log"
    "github.com/Leonard-Atorough/ghast"
    "github.com/Leonard-Atorough/ghast/middleware"
)

app := ghast.New()

// Apply global middleware
app.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))
app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{}))

// Router with rate limiting
apiRouter := ghast.NewRouter()
apiRouter.Use(middleware.RateLimitMiddleware(middleware.RateLimitOptions{
    RequestsPerMinute: 100,
}))

apiRouter.Get("/data", dataHandler)
app.Route("/api", apiRouter)

// CORS for public endpoints
publicRouter := ghast.NewRouter()
publicRouter.Use(middleware.CorsMiddleware(middleware.CorsOptions{
    AllowedOrigins: []string{"*"},
}))

publicRouter.Get("/status", statusHandler)
app.Route("/public", publicRouter)

app.Listen(":8080")
```

## Testing

Unit tests are included in the repository:

```bash
go test ./...
```

Tests cover:

- ✅ Router matching (GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH)
- ✅ Route parameters and wildcards
- ✅ 404 and error handling
- ✅ Middleware chaining and order
- ✅ Request parsing (query, JSON, headers)
- ✅ Response writing and chainable methods
- ✅ Mixed middleware application
- ✅ Sub-router mounting and prefix stripping

## Documentation

Complete documentation is available in the `docs/v0.5.0/` directory:

- **[Routing & Middleware Guide](docs/v0.5.0/guide.md)** - Learn routing, path parameters, middleware composition, and sub-routers
- **[API Reference](docs/v0.5.0/api-reference.md)** - Comprehensive reference for all types, interfaces, and functions
- **[Middleware Reference](docs/v0.5.0/middleware-reference.md)** - Detailed guide to built-in middleware (Recovery, Request ID, Response Time, CORS, Rate Limiting)

## Examples

Learning examples are available in the `examples/` folder:

1. **basic_example.go** - Simple hello-world and JSON API example
2. **grouped_routes_example.go** - Sub-routers and route organization
3. **custom_middleware_example.go** - Writing custom middleware

## Development Roadmap

### Completed (v0.5.0) ✅

- [x] Core router with HTTP methods
- [x] Express-like convenience methods
- [x] Response/request helpers
- [x] Middleware system with composition
- [x] Path parameters (`:id`, `:slug`)
- [x] Sub-routers with prefix mounting
- [x] Built-in middleware (Recovery, Request ID, Response Time, CORS, Rate Limit)
- [x] Comprehensive documentation
- [x] Unit tests

### Future Enhancements (Optional)

- [ ] Static file serving
- [ ] Cookie handling
- [ ] Form data parsing
- [ ] WebSocket support
- [ ] Template rendering
- [ ] Request validation middleware
- [ ] Compressed response support

## Code Quality

- **Readable:** Clear variable names, comprehensive comments
- **Testable:** Comprehensive unit tests demonstrating all features
- **Documented:** Complete API documentation and guides
- **Idiomatic:** Follows Go conventions and best practices
- **Educational:** Clean code designed to be understandable and learnable

## Contributing

This is primarily an educational project. Feel free to:

- Fork and modify for learning purposes
- Read the source code to understand HTTP server design
- Add more examples or tests
- Improve documentation
- Experiment with additional features

## License

MIT - Feel free to learn from, modify, and adapt this code.
See [LICENSE.md](LICENSE.md) for details.

---

**Built with ❤️ for learning Go, networking, and elegant API design**

For questions or to dive deeper:

- Check the [examples/](examples/) folder for working examples
- Read the [documentation](docs/v0.5.0/) for comprehensive guides
- Review the well-commented source code
