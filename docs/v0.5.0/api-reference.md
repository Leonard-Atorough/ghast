---
title: "API Reference"
description: "Comprehensive reference for Ghast's public API, including top-level functions, Ghast, Router, Handler, Middleware, Request, and Response types and methods."
authors: ["Leonard"]
dateCreated: "2024-02-23"
dateUpdated: "2024-02-23"
---

# V 0.5.0 API Reference

> Ghast 0.5.0 requires Go 1.21 or later.

---

## Top-Level Functions

These are package-level functions exported directly from the `ghast` package. They serve as entry points and utilities for constructing the core building blocks of your application.

---

### `New() *Ghast`

Creates and returns a new `Ghast` instance. This is the primary entry point for the framework — use it to configure your server, register routes, apply middleware, and start listening.

**Returns:**

- `*Ghast`: A new application instance with an empty root router and no middleware.

**Example:**

```go
app := ghast.New()

app.Get("/", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.Plain(200, "Hello, World!")
})

app.Listen(":8080")
```

---

### `NewRouter() Router`

Creates and returns a new `Router` instance with empty routes and no middleware. Use this to build modular route groups that can be mounted onto a `Ghast` instance via `app.Route()`.

**Returns:**

- `Router`: A new router ready for route registration.

**Example:**

```go
userRouter := ghast.NewRouter()
userRouter.Get("/", listUsersHandler)
userRouter.Post("/", createUserHandler)

app.Route("/users", userRouter)
```

---

### `NewHandlerBuilder(handler Handler) *HandlerBuilder`

Creates a new `HandlerBuilder` wrapping the given base handler. Provides a fluent API for composing middleware onto a specific handler without applying it globally to a router or the whole app.

**Props:**

- `handler` (Handler): The base handler to wrap with middleware.

**Returns:**

- `*HandlerBuilder`: A builder instance for chaining `.Use()` calls.

**Example:**

```go
finalHandler := ghast.NewHandlerBuilder(myHandler).
    Use(authMiddleware).
    Use(loggingMiddleware).
    Build()

app.Get("/protected", finalHandler)
```

---

### `Error(rw ResponseWriter, statusCode int, message string) error`

Sends an HTTP error response with a JSON body. The response body has the shape `{ "status": <code>, "error": "<message>" }`.

**Props:**

- `rw` (ResponseWriter): The response writer for the current request.
- `statusCode` (int): HTTP status code to send (e.g. `404`, `500`).
- `message` (string): Human-readable error description.

**Returns:**

- `error`: Any error encountered while writing the response.

**Example:**

```go
app.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    user, ok := findUser(r.Param("id"))
    if !ok {
        ghast.Error(w, 404, "user not found")
        return
    }
    w.JSON(200, user)
})
```

---

### `ErrorString(rw ResponseWriter, statusCode int, message string) error`

Sends an HTTP error response with a plain text body. Prefer `Error` when the client expects JSON; use this when a plain text error is more appropriate.

**Props:**

- `rw` (ResponseWriter): The response writer for the current request.
- `statusCode` (int): HTTP status code to send.
- `message` (string): Error message to send as plain text.

**Returns:**

- `error`: Any error encountered while writing the response.

**Example:**

```go
app.Get("/ping", func(w ghast.ResponseWriter, r *ghast.Request) {
    if !serviceHealthy() {
        ghast.ErrorString(w, 503, "Service Unavailable")
        return
    }
    w.Plain(200, "pong")
})
```

---

## Application

The `Ghast` struct is the central application object returned by `New()`. It owns the root router, all mounted sub-routers, and the global middleware stack. All of its methods return `*Ghast` (except `Listen`) so calls can be chained.

---

### `app.Get(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP GET requests at the given path on the root router. Returns the app for chaining.

**Props:**

- `path` (string): URL path, supports dynamic segments prefixed with `:` (e.g. `"/users/:id"`).
- `handler` (Handler): Handler to invoke when the route matches.
- `middlewares` (...Middleware): Optional middleware applied only to this route.

**Example:**

```go
app.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"id": r.Param("id")})
})
```

---

### `app.Post(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP POST requests at the given path on the root router. Returns the app for chaining.

---

### `app.Put(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP PUT requests at the given path on the root router. Returns the app for chaining.

---

### `app.Delete(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP DELETE requests at the given path on the root router. Returns the app for chaining.

---

### `app.Patch(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP PATCH requests at the given path on the root router. Returns the app for chaining.

---

### `app.Head(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP HEAD requests at the given path on the root router. Returns the app for chaining.

---

### `app.Options(path string, handler Handler, middlewares ...Middleware) *Ghast`

Registers a handler for HTTP OPTIONS requests at the given path on the root router. Returns the app for chaining.

> **Note:** `Post`, `Put`, `Delete`, `Patch`, `Head`, and `Options` all share the same signature as `Get` above. See `Get` for a full description of the parameters.

---

### `app.Use(middleware Middleware) *Ghast`

Adds a middleware function to the global middleware stack. Global middleware runs on every request regardless of which router or route handles it. Returns the app for chaining.

**Props:**

- `middleware` (Middleware): Middleware function to apply globally.

**Example:**

```go
app.Use(func(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        log.Printf("%s %s", r.Method, r.Path)
        next.ServeHTTP(w, r)
    })
})
```

---

### `app.Route(prefix string, router Router, middlewares ...Middleware) *Ghast`

Mounts a `Router` under the given path prefix, with optional middleware scoped to that router. Incoming requests whose path begins with `prefix` are stripped of the prefix and forwarded to the mounted router. Returns the app for chaining.

**Props:**

- `prefix` (string): Path prefix to match (e.g. `"/api"`).
- `router` (Router): A router created with `ghast.NewRouter()` or returned by your own factory.
- `middlewares` (...Middleware): Optional middleware applied only to routes within this router.

**Example:**

```go
userRouter := ghast.Router()
userRouter.Get("/", listUsersHandler)
userRouter.Get("/:id", getUserHandler)

app.Route("/users", userRouter, authMiddleware)
// GET /users     → listUsersHandler (through authMiddleware)
// GET /users/42  → getUserHandler   (through authMiddleware)
```

---

### `app.Router() Router`

Returns the root `Router` instance that backs the app directly. In most cases you won't need this — `app.Get()`, `app.Post()`, etc. register routes on the root router for you. Use `Router()` when you need direct access to the `Router` interface, for example to call `router.Handle()` with a custom method or to pass the root router to a helper that accepts a `Router`.

**Returns:**

- `Router`: The app's root router.

**Example:**

```go
// Register a custom HTTP method directly on the root router
app.Router().Handle("PURGE", "/cache", purgeHandler)
```

---

### `app.Listen(addr string) error`

Starts the HTTP server on the given address and blocks until a fatal error occurs. This should be the last call in your application setup.

**Props:**

- `addr` (string): TCP address to listen on (e.g. `":8080"`, `"0.0.0.0:3000"`).

**Returns:**

- `error`: Any error encountered while starting or running the server.

**Example:**

```go
if err := app.Listen(":8080"); err != nil {
    log.Fatal(err)
}
```

---

## Router

`Router` is an interface that defines how routes are registered and how incoming requests are dispatched. The concrete implementation is returned by `ghast.NewRouter()`.

### Interface

```go
type Router interface {
    Handle(method string, path string, handler Handler, middlewares ...Middleware)
    Get(path string, handler Handler, middlewares ...Middleware) Router
    Post(path string, handler Handler, middlewares ...Middleware) Router
    Put(path string, handler Handler, middlewares ...Middleware) Router
    Delete(path string, handler Handler, middlewares ...Middleware) Router
    Patch(path string, handler Handler, middlewares ...Middleware) Router
    Head(path string, handler Handler, middlewares ...Middleware) Router
    Options(path string, handler Handler, middlewares ...Middleware) Router
    ServeHTTP(ResponseWriter, *Request)
    Use(middleware Middleware) Router
}
```

---

### `router.Handle(method, path string, handler Handler, middlewares ...Middleware)`

Registers a handler for an explicit HTTP method and path. All HTTP verb convenience methods (`Get`, `Post`, etc.) delegate to this. Use it directly when you need a method not covered by the convenience set (e.g. a custom verb).

When called, the path is compiled into a regex pattern stored internally for dynamic segment matching on each request.

**Props:**

- `method` (string): HTTP method string (e.g. `"GET"`, `"PURGE"`).
- `path` (string): URL path, supports dynamic segments prefixed with `:` (e.g. `"/users/:id"`).
- `handler` (Handler): Handler to invoke on a match.
- `middlewares` (...Middleware): Optional middleware applied only to this route, layered after any router-level middleware.

**Example:**

```go
router.Handle("PURGE", "/cache/:key", func(w ghast.ResponseWriter, r *ghast.Request) {
    purge(r.Param("key"))
    w.Plain(200, "purged")
})
```

---

### `router.Get(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP GET requests at the given path. Returns the router for chaining.

**Props:**

- `path` (string): URL path, supports dynamic segments (e.g. `"/posts/:id"`).
- `handler` (Handler): Handler to invoke when the route matches.
- `middlewares` (...Middleware): Optional per-route middleware.

**Example:**

```go
router.Get("/posts/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"id": r.Param("id")})
})
```

---

### `router.Post(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP POST requests. Returns the router for chaining.

---

### `router.Put(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP PUT requests. Returns the router for chaining.

---

### `router.Delete(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP DELETE requests. Returns the router for chaining.

---

### `router.Patch(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP PATCH requests. Returns the router for chaining.

---

### `router.Head(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP HEAD requests. Returns the router for chaining.

---

### `router.Options(path string, handler Handler, middlewares ...Middleware) Router`

Registers a handler for HTTP OPTIONS requests. Returns the router for chaining.

> **Note:** `Post`, `Put`, `Delete`, `Patch`, `Head`, and `Options` share the same signature as `Get` above.

---

### `router.Use(middleware Middleware) Router`

Adds middleware that applies to all routes registered on this router. Middleware added via `Use` is prepended before any per-route middleware. Returns the router for chaining.

**Props:**

- `middleware` (Middleware): Middleware function to apply to all routes.

**Example:**

```go
router.Use(loggingMiddleware)
router.Use(authMiddleware)

router.Get("/users", listUsersHandler)  // both middleware functions apply here
```

---

### `router.ServeHTTP(w ResponseWriter, req *Request)`

Dispatches an incoming request to the correct registered handler. Exact path matches are tried first; if none is found, dynamic (regex-compiled) routes are tried in turn. If no route matches, a `404 Not Found` response is written.

This method satisfies the `Handler` interface, which means a `Router` can itself be used anywhere a `Handler` is accepted.

---

## Handler

### Interface

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any type that implements `ServeHTTP` can be used as a handler. The framework provides `HandlerFunc` as a convenience adapter so ordinary functions can satisfy the interface without defining a named type.

---

### `HandlerFunc`

Function type that implements `Handler`. Cast any compatible function literal to `HandlerFunc` to use it directly as a handler.

```go
type HandlerFunc func(ResponseWriter, *Request)
```

**Example:**

```go
var h ghast.Handler = ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
    w.Plain(200, "ok")
})
```

---

## Middleware

### Type

```go
type Middleware func(Handler) Handler
```

A `Middleware` is a function that receives the next `Handler` in the chain and returns a new `Handler` that wraps it. This lets you execute logic before and/or after the inner handler.

**Example:**

```go
var logger ghast.Middleware = func(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        log.Printf("--> %s %s", r.Method, r.Path)
        next.ServeHTTP(w, r)
        log.Printf("<-- %s %s done", r.Method, r.Path)
    })
}
```

---

### `HandlerBuilder`

Fluent builder for applying an ordered set of middleware to a single handler. Use `NewHandlerBuilder()` to create one, chain `.Use()` calls, and call `.Build()` to retrieve the wrapped handler.

```go
type HandlerBuilder struct { ... }
```

#### `hb.Use(middleware Middleware) *HandlerBuilder`

Wraps the current handler with the given middleware. Returns the builder for chaining. Middleware is applied in call order — the first `Use` call wraps outermost.

**Props:**

- `middleware` (Middleware): Middleware function to layer onto the handler.

#### `hb.Build() Handler`

Returns the final handler with all middleware applied.

**Example:**

```go
handler := ghast.NewHandlerBuilder(myHandler).
    Use(rateLimitMiddleware).
    Use(authMiddleware).
    Build()

app.Get("/api/data", handler)
```

---

## Request

Represents a parsed incoming HTTP request. Populated by the framework before your handler is called.

### Type

```go
type Request struct {
    Method   string            // HTTP method (e.g. "GET", "POST")
    Path     string            // URL path without query string (e.g. "/users/42")
    Headers  map[string]string // Request headers
    Body     string            // Raw request body as a string
    Version  string            // HTTP version (e.g. "HTTP/1.1")
    Params   map[string]string // Dynamic path parameters extracted by the router
    Queries  map[string]string // Query string parameters
    ClientIP string            // Remote client IP address
}
```

---

### `req.Param(key string) string`

Returns the value of a dynamic path parameter by name. Returns an empty string if the parameter does not exist.

**Props:**

- `key` (string): Parameter name as declared in the route path (without the `:` prefix).

**Example:**

```go
// Route: /users/:id
app.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")  // "42" for a request to /users/42
    w.Plain(200, id)
})
```

---

### `req.Query(key string) string`

Returns the value of a query string parameter by name. Returns an empty string if the key is not present.

**Props:**

- `key` (string): Query parameter name.

**Example:**

```go
// Request: GET /search?q=ghast&page=2
app.Get("/search", func(w ghast.ResponseWriter, r *ghast.Request) {
    q    := r.Query("q")     // "ghast"
    page := r.Query("page")  // "2"
    _ = q; _ = page
})
```

---

### `req.GetHeader(key string) string`

Returns the value of a request header. Matching is case-insensitive. Returns an empty string if the header is absent.

**Props:**

- `key` (string): Header name (e.g. `"Content-Type"`, `"authorization"`).

**Example:**

```go
token := r.GetHeader("Authorization")
```

---

### `req.ContentType() string`

Convenience shorthand for `req.GetHeader("Content-Type")`.

**Example:**

```go
if r.ContentType() == "application/json" {
    // decode JSON body
}
```

---

### `req.JSON(v any) error`

Unmarshals the request body as JSON into `v`. `v` must be a pointer.

**Props:**

- `v` (any): Pointer to a value to unmarshal into.

**Returns:**

- `error`: Any JSON decoding error.

**Example:**

```go
type CreateUserBody struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

app.Post("/users", func(w ghast.ResponseWriter, r *ghast.Request) {
    var body CreateUserBody
    if err := r.JSON(&body); err != nil {
        ghast.Error(w, 400, "invalid JSON")
        return
    }
    w.JSON(201, body)
})
```

---

## ResponseWriter

Interface for constructing and sending HTTP responses. Every handler receives a fully initialised `ResponseWriter`; methods that set state (`Status`, `SetHeader`) return the writer itself so calls can be chained.

### Interface

```go
type ResponseWriter interface {
    Header() map[string]string
    Status(statusCode int) ResponseWriter
    SetHeader(key, value string) ResponseWriter
    Send([]byte) (int, error)
    SendString(string) (int, error)
    JSON(statusCode int, data interface{}) error
    JSONPretty(statusCode int, data interface{}) error
    HTML(statusCode int, html string) error
    Plain(statusCode int, text string) error
}
```

---

### `rw.Header() map[string]string`

Returns the mutable response headers map. Set headers directly on the returned map before the first write, or use `SetHeader` for a chainable alternative.

**Example:**

```go
rw.Header()["X-Request-Id"] = "abc123"
```

---

### `rw.Status(statusCode int) ResponseWriter`

Sets the HTTP status code for the response. Has no effect after the first write. Returns the writer for chaining.

**Props:**

- `statusCode` (int): Standard HTTP status code (e.g. `200`, `404`, `500`).

**Example:**

```go
rw.Status(201).SetHeader("Location", "/users/42")
```

---

### `rw.SetHeader(key, value string) ResponseWriter`

Sets a single response header. Returns the writer for chaining.

**Props:**

- `key` (string): Header name.
- `value` (string): Header value.

**Example:**

```go
rw.SetHeader("Cache-Control", "no-store")
```

---

### `rw.Send(data []byte) (int, error)`

Writes a raw byte slice as the response body. The status line and headers are flushed on the first call.

**Props:**

- `data` ([]byte): Bytes to write.

**Returns:**

- `int`: Number of bytes written.
- `error`: Any write error.

---

### `rw.SendString(s string) (int, error)`

Writes a plain string as the response body. Convenience wrapper around `Send`.

**Props:**

- `s` (string): String to write.

**Returns:**

- `int`: Number of bytes written.
- `error`: Any write error.

**Example:**

```go
rw.Status(200).SendString("pong")
```

---

### `rw.JSON(statusCode int, data interface{}) error`

Marshals `data` to JSON, sets `Content-Type: application/json`, and writes the response.

**Props:**

- `statusCode` (int): HTTP status code.
- `data` (interface{}): Value to marshal as JSON.

**Returns:**

- `error`: Any marshalling or write error.

**Example:**

```go
rw.JSON(200, map[string]string{"status": "ok"})
```

---

### `rw.JSONPretty(statusCode int, data interface{}) error`

Same as `JSON` but the output is indented with two spaces per level — useful during development or for APIs where human-readable output is expected.

**Props:**

- `statusCode` (int): HTTP status code.
- `data` (interface{}): Value to marshal.

**Returns:**

- `error`: Any marshalling or write error.

---

### `rw.HTML(statusCode int, html string) error`

Sends an HTML response, setting `Content-Type: text/html`.

**Props:**

- `statusCode` (int): HTTP status code.
- `html` (string): HTML string to send.

**Returns:**

- `error`: Any write error.

**Example:**

```go
rw.HTML(200, "<h1>Hello</h1>")
```

---

### `rw.Plain(statusCode int, text string) error`

Sends a plain text response, setting `Content-Type: text/plain`.

**Props:**

- `statusCode` (int): HTTP status code.
- `text` (string): Text to send.

**Returns:**

- `error`: Any write error.

**Example:**

```go
rw.Plain(200, "pong")
```

---

## Error

### Types

#### `HTTPError`

The JSON shape produced by `Error()`. Also available directly if you need to construct or inspect error values.

```go
type HTTPError struct {
    StatusCode int    `json:"status"`
    Message    string `json:"error"`
}
```

---

## HTTP Method Constants

String constants for the standard HTTP methods, exported for use in `Handle` calls and request comparisons.

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
