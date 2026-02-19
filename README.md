# 🌪️ Gust

An Express.js-inspired HTTP framework for Go. Built from first principles to learn networking, concurrency, and elegant API design.

## Why Gust?

**Gust** bridges two worlds:

- **For Node.js developers:** Familiar Express-like API for building HTTP services in Go
- **For learning:** Understand HTTP, TCP, and concurrency by reading clean, documented code

This is an educational framework that prioritizes clarity and approachability. While not production-ready (use [Echo](https://echo.labstack.com/), [Gin](https://gin-gonic.com/), or stdlib for production), Gust is perfect for:

- Learning Go networking and concurrency
- Building small projects or APIs
- Understanding HTTP server design patterns
- **Portfolio projects** showing thoughtful architecture

## Quick Start

### Installation

```bash
go get github.com/YourUsername/gust
```

### Hello World

```go
package main

import (
	"log"
	"gust/lib"
)

func main() {
	router := gust.NewRouter()

	// Simple GET route
	router.Get("/hello", gust.HandlerFunc(func(w gust.ResponseWriter, r *gust.Request) {
		w.JSON(200, map[string]string{"message": "Hello, World!"})
	}))

    server := gust.NewServer(router)
	// Start server (see main.go for full implementation)
	if err := server.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
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

```go
// Global middleware - applies to all routes
app.Use(LoggingMiddleware)
app.Use(RecoveryMiddleware)

// Path-specific middleware
app.UsePath("/api", authMiddleware)

// Custom middleware
customMiddleware := func(next Handler) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
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

## Architecture

```
gust/
├── lib/                    # Core framework
│   ├── router.go          # HTTP routing
│   ├── handler.go         # Handler interface
│   ├── request.go         # Request parsing & helpers
│   ├── response.go        # Response writing & helpers
│   ├── middleware.go      # Middleware system
│   └── error.go           # Error handling
│
├── examples/              # Example handlers
│   └── handlers.go        # Sample implementations
│
├── main.go               # TCP server & entry point
├── go.mod
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
import "gust/lib"

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

router := lib.NewRouter()

// GET /users/:id
router.Get("/users/:id", lib.HandlerFunc(func(w lib.ResponseWriter, r *lib.Request) {
    id := r.Param("id")
    user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
    w.JSON(200, user)
}))

// POST /users
router.Post("/users", lib.HandlerFunc(func(w lib.ResponseWriter, r *lib.Request) {
    var user User
    if err := r.JSON(&user); err != nil {
        lib.Error(w, 400, "Invalid request body")
        return
    }
    user.ID = 2
    w.Status(201).JSON(user)
}))

// GET /health (query parameter example)
router.Get("/health", lib.HandlerFunc(func(w lib.ResponseWriter, r *lib.Request) {
    verbose := r.Query("verbose") == "true"
    status := "healthy"
    w.JSON(200, map[string]interface{}{
        "status": status,
        "verbose": verbose,
    })
}))
```

### Middleware

```go
// Logging middleware
loggingMiddleware := func(next lib.Handler) lib.Handler {
    return lib.HandlerFunc(func(w lib.ResponseWriter, r *lib.Request) {
        log.Printf("[%s] %s %s", r.Method, r.Path, r.Version)
        next.ServeHTTP(w, r)
    })
}

// Authentication middleware
authMiddleware := func(next lib.Handler) lib.Handler {
    return lib.HandlerFunc(func(w lib.ResponseWriter, r *lib.Request) {
        token := r.GetHeader("Authorization")
        if token == "" {
            lib.Error(w, 401, "Missing authorization token")
            return
        }
        next.ServeHTTP(w, r)
    })
}

router := lib.NewRouter()
router.Use(loggingMiddleware)
router.UsePath("/api", authMiddleware)
```

## Testing

Basic tests are included in `lib/gust_test.go`:

```bash
go test ./lib/...
```

Tests cover:

- ✅ Router matching (GET, POST, PUT, DELETE, HEAD, OPTIONS)
- ✅ 404 handling
- ✅ Middleware chaining order
- ✅ Request query parameters
- ✅ JSON marshaling/unmarshaling
- ✅ Response header chaining

## Development Roadmap

### Phase 1: Foundation ✅

- [x] Core router with HTTP methods
- [x] Express-like convenience methods
- [x] Response/request helpers
- [x] Middleware system
- [x] Basic tests
- [x] Documentation

### Phase 2: Dynamic Routing (Planned)

- [ ] Path parameters (`:id`, `:slug`)
- [ ] Regex route matching
- [ ] Route groups with shared middleware
- [ ] 404 and error handling improvements
- [ ] Graceful shutdown support
- [ ] Server configuration options (timeouts, max connections)
- [ ] Server-level middleware (global middleware that applies to all routes)

### Phase 3: Advanced Features (Optional)

- [ ] JSON schema validation
- [ ] Rate limiting middleware
- [ ] CORS middleware
- [ ] Static file serving
- [ ] Cookie handling

## Learning Resources

If you want to understand how this works:

1. Start with [lib/handler.go](lib/handler.go) - Learn the Handler interface
2. Read [lib/router.go](lib/router.go) - See how routing works
3. Check [lib/middleware.go](lib/middleware.go) - Understand middleware chaining
4. Look at [main.go](main.go) - See the TCP server implementation
5. Study [lib/gust_test.go](lib/gust_test.go) - Test examples showing all features

## Helpful Go Packages

- `net` - TCP/IP networking
- `bufio` - Buffered I/O for reading HTTP headers
- `encoding/json` - JSON marshaling/unmarshaling
- `strings` - String utilities for parsing

## Code Style & Quality

- **Readable:** Clear variable names, comprehensive comments
- **Testable:** Tests demonstrate intended usage
- **Documented:** Godoc comments on all public types and functions
- **Idiomatic:** Follows Go conventions and philosophy

## Contributing

This is primarily an educational project. Feel free to:

- Fork and modify for learning
- Add more test examples
- Improve documentation
- Experiment with features

## License

MIT - Feel free to learn from, modify, and adapt this code.

---

**Built with ❤️ as a learning project**

Questions? Check the example handlers in `examples/` or read through the well-commented code in `lib/`.
