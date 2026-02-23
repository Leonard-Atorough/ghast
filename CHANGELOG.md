# Changelog

All notable changes to the Ghast HTTP framework are documented in this file. This project follows [Semantic Versioning](https://semver.org/).

---

## [0.5.0] - 2026-02-23

### Major Changes

#### API Refactoring

- **Renamed `Server` to `Ghast`**: The primary application type is now called `Ghast` for consistency with the framework name
- **Introduced `Router` mounting via `app.Route()`**: Cleaner separation of concerns with mounted sub-routers
- **Replaced `NewRouter(prefix)` with module-level `NewRouter()`**: Routers are now created independently and mounted to the app, promoting modularity

#### Handler & Middleware Improvements

- **New `HandlerBuilder` utility**: Fluent API for composing middleware onto individual handlers without global registration
- **`NewHandlerBuilder(handler Handler)` function**: Factory for creating handler builders with chainable `.Use()` method
- **Enhanced middleware composition**: Fine-grained control over middleware layering at route level

#### Application API Enhancements

- **New `app.Router()` method**: Direct access to the root router for advanced use cases
- **New `app.Route(prefix string, router Router, middlewares ...Middleware)` method**: Mount standalone routers with optional scoped middleware
- **Improved method chaining**: All `Ghast` methods (except `Listen`) return `*Ghast` for fluent configuration

#### Error Handling

- **Package-level `Error()` function**: Sends JSON error responses with consistent structure `{ "status": <code>, "error": "<message>" }`
- **`ErrorString()` function**: Alternative for plain text error responses when JSON isn't appropriate

#### Request & Response Enhancements

- **Implicit content-type handling**: Methods like `Plain()`, `HTML()` automatically set proper content-type headers
- **Consistent response builder pattern**: Chainable response methods for intuitive API usage

### Features

**Complete HTTP/1.1 server implementation** with all core features:

- Full routing with dynamic path parameters (`:param` syntax)
- Global and route-specific middleware chains
- Automatic path parameter extraction
- Comprehensive request/response APIs
- Multiple response format shortcuts

**Built-in middleware library**:

- CORS support (`middleware/cors.go`)
- Rate limiting (`middleware/rate_limit.go`)
- Panic recovery (`middleware/recovery.go`)
- Request ID injection (`middleware/request_id.go`)
- Response timing (`middleware/response_time.go`)

**Full API documentation** with examples for every public type and method

### Breaking Changes

- `Server` type renamed to `Ghast` — update all `ghast.New()` usage to work with `Ghast` type
- `NewRouter(prefix string)` removed — use `NewRouter()` + `app.Route(prefix, router)` pattern
- `app.NewRouter(prefix)` method removed — migrate to standalone `Router` creation
- Router is no longer tied to a specific prefix at creation time

### Migration Guide

**Old code (v0.4.0)**:

```go
app := ghast.New()
apiRouter := app.NewRouter("/api")
apiRouter.Get("/users", handler)
```

**New code (v0.5.0)**:

```go
app := ghast.New()
apiRouter := ghast.NewRouter()
apiRouter.Get("/users", handler)
app.Route("/api", apiRouter)
```

---

## [0.4.0] - 2026-02-22

### Features

**JSON Support**:

- `w.JSON(statusCode int, data interface{})` - Automatic JSON marshaling with `application/json` content-type
- `w.JSONPretty(statusCode int, data interface{})` - Pretty-printed JSON for development
- `r.JSON(v interface{})` - Parse request body as JSON into provided value

**Content Negotiation**:

- Response format shortcuts: `Plain()`, `HTML()`, `JSON()`
- Automatic content-type headers based on response method
- Multiple format support for same endpoint

**Enhanced Response Writer**:

- All convenience methods return the `ResponseWriter` for chaining
- Chainable `Status()` and `SetHeader()` methods
- Improved error handling for malformed JSON

### Examples

New comprehensive examples added:

- JSON API endpoints with request/response marshaling
- Content negotiation patterns
- Error handling with JSON responses

### Documentation

- Updated API reference with JSON methods
- Added JSON handling guide to main documentation
- Examples for common JSON patterns (list, create, update, delete)

### Bug Fixes

- Fixed Content-Type header encoding issues
- Improved robustness of JSON parsing with malformed input
- Better error messages for JSON unmarshaling failures

---

## [0.3.0] - 2026-02-21

### Features

**Middleware Ecosystem**:

- `middleware/cors.go`: Cross-Origin Resource Sharing support with configurable origins
- `middleware/recovery.go`: Panic recovery to prevent server crashes on unexpected errors
- `middleware/request_id.go`: Automatic request ID injection for tracing
- `middleware/rate_limit.go`: Simple per-IP rate limiting
- `middleware/response_time.go`: Automatic response timing tracking

**Global and Path-Specific Middleware**:

- Apply middleware globally with `app.Use(middleware Middleware)`
- Apply middleware to routers with `router.Use(middleware Middleware)`
- Apply middleware to individual routes with `app.Get(path, handler, middleware1, middleware2, ...)`
- Middleware composition with predictable execution order

**Middleware Chain Utilities**:

- `ChainMiddleware(handler Handler, middlewares []Middleware) Handler` - Utility for composing multiple middleware
- Middleware execution follows decorative pattern for clean control flow

### Examples

- Basic logging middleware
- Authentication middleware pattern
- CORS configuration examples
- Rate limiting setup
- Composing multiple middleware on single route

### Documentation

- Comprehensive middleware reference guide
- Best practices for middleware design
- Examples for each built-in middleware
- Middleware execution order documentation

### Bug Fixes

- Fixed middleware execution order inconsistency
- Improved error propagation through middleware chains
- Fixed global middleware not applying to all routes

---

## [0.2.0] - 2026-02-20

### Features

**Dynamic Routing with Path Parameters**:

- Support for `:param` syntax in route paths (e.g., `/users/:id`)
- Automatic path parameter extraction
- `r.Param(key string)` method to retrieve parameters
- Case-sensitive parameter matching

**Enhanced Router**:

- Trie-based route matching for efficient lookups
- Support for exact path matching (fast path)
- Regex compilation for dynamic routes (lazy evaluated per request)
- Proper 404 handling when no route matches

**Router Prefix Support**:

- `app.NewRouter(prefix string)` creates sub-routers with path prefixes
- Clean API for route organization
- Example: `/api` router can contain routes that resolve to `/api/users`, `/api/posts`, etc.

### Examples

Added examples showcasing:

- Single parameter routes: `/users/:id`
- Multiple parameters: `/users/:id/posts/:postId`
- RESTful API patterns
- Parameter validation at handler level

### Documentation

- Dynamic routing guide with examples
- Parameter extraction documentation
- Route matching algorithm explanation
- Performance considerations for regex routes

### Testing

- Unit tests for exact path matching
- Unit tests for dynamic parameter extraction
- Edge cases: overlapping routes, special characters in parameters
- 404 behavior verification

### Bug Fixes

- Fixed parameter extraction with query strings
- Improved route matching priority (exact before dynamic)
- Fixed UTF-8 handling in path parameters

---

## [0.1.0] - 2026-02-18

### Initial Release

This is the foundation version of Ghast, featuring core HTTP server capabilities built from first principles.

#### Core Components

**TCP Server Foundation** (`server.go`):

- TCP listener on configurable port
- Connection acceptance loop with goroutine-per-connection concurrency model
- Graceful connection handling with defer cleanup
- Connection keep-alive support for persistent connections

**HTTP Request Parsing** (`request.go`):

- Complete HTTP/1.1 request parsing from raw TCP data
- Request line parsing (method, path, HTTP version)
- Header parsing with case-insensitive lookup
- Body reading with Content-Length respect
- Query parameter extraction
- Path and query parameter mapping

**HTTP Response Formatting** (`response.go`):

- HTTP/1.1 response generation with status line, headers, and body
- Response headers management
- Chainable response methods for intuitive API
- Multiple response format shortcuts:
  - `JSON(statusCode, data)` - JSON response with automatic marshaling
  - `Plain(statusCode, data)` - Plain text response
  - `HTML(statusCode, data)` - HTML response with proper content-type
  - `Send([]byte)` - Raw byte response
  - `SendString(string)` - String response

**Handler Interface** (`handler.go`):

- `Handler` interface with `ServeHTTP(ResponseWriter, *Request)` method
- `HandlerFunc` adapter for function-based handlers
- Familiar interface matching Go's `net/http` patterns
- Type assertion and interface satisfaction for flexible handler composition

**Routing** (`router.go`):

- `Router` interface for handler registration and request dispatching
- Path-based handler mapping
- HTTP method routing (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- Exact path matching with O(1) lookup performance
- 404 handling for unmatched routes
- Router interface satisfies Handler interface for composition

**Middleware System** (`middleware.go`):

- `Middleware` function type: `func(Handler) Handler`
- Decorator pattern for handler wrapping
- Global middleware with `app.Use(middleware)`
- Router-level middleware with `router.Use(middleware)`
- Route-specific middleware support
- Middleware composition and chaining utilities

**Error Handling** (`error.go`):

- `HTTPError` struct for structured error responses
- `Error(rw ResponseWriter, statusCode, message)` for JSON errors
- `ErrorString(rw ResponseWriter, statusCode, message)` for text errors

**HTTP Status Codes** (`http_status.go`):

- Comprehensive HTTP status code constants
- Common status codes for easy reference

#### Type System

```go
type Server struct
type Router interface
type Handler interface
type Middleware func(Handler) Handler
type ResponseWriter interface
type HandlerFunc func(ResponseWriter, *Request)
type Request struct
type HandlerBuilder struct
```

#### Key Design Decisions

1. **Handler Interface**: Follows Go standard library patterns for familiarity and compatibility
2. **Middleware Pattern**: Decorator pattern for composable, reusable handler wrapping
3. **ResponseWriter Abstraction**: Separates HTTP protocol concerns from business logic
4. **Goroutine-per-Connection**: Lightweight concurrency for handling multiple simultaneous clients
5. **Interface-Based Design**: Promotes loose coupling and testability

#### Express.js-Inspired API

Familiar DSL for Node.js developers transitioning to Go:

```go
app := ghast.New()
app.Get("/hello", handler)
app.Post("/api/users", handler)
app.Use(loggingMiddleware)
app.Listen(":8080")
```

#### Examples

Included examples demonstrate:

- Basic "Hello World" server
- Grouped routes with prefix patterns
- Custom middleware creation
- Request/response handling patterns
- Error handling approaches

#### Documentation

- Comprehensive API reference with all public types and methods
- Getting started guide with quick start examples
- Middleware reference guide
- Architecture overview

#### Testing

- Unit tests for request parsing (`ghast_test.go`)
- Integration tests for routing (`router_test.go`)
- Server lifecycle tests (`server_test.go`)
- Coverage includes core parsing, routing, and response generation

#### Code Quality

- Clean, readable codebase
- Full godoc comments on all public APIs
- Organized file structure reflecting functional areas
- Clear separation of concerns
- Comprehensive examples and documentation

---

## Legend

- **Foundation** = Core component
- **Feature** = Feature implementation
- **Enhancement** = New capability
- **Tools** = Tooling / Utilities
- **Bugfix** = Bug fix
- **Docs** = Documentation

---

## Release Notes Format

Each version documents:

- **Major Changes**: Breaking changes and architectural decisions
- **Features**: New functionality added
- **Bug Fixes**: Issues resolved
- **Breaking Changes**: API changes requiring user code migration
- **Migration Guide**: How to update code from previous version
- **Examples**: Real-world usage patterns showcasing new features
- **Documentation**: Updates to guides and API references

---

## Version Philosophy

Ghast follows Semantic Versioning:

- **Major (v1.0.0)**: Production-ready status, stable API guarantees
- **Minor (v0.x.0)**: New features, backward compatible
- **Patch (v0.0.x)**: Bug fixes and documentation only

Current status: **Pre-1.0 (Learning & Experimentation)** — APIs may change between versions, but effort is made to document migrations.

---

**Last Updated**: February 23, 2026
