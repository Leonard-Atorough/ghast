---
title: "Routing and Middleware"
description: "Learn how to define routes and handlers in Ghast 0.5.0, including app-level routing, composable sub-routers, route parameters, and middleware."
authors: ["Leonard"]
dateCreated: "2024-02-23"
dateUpdated: "2024-02-23"
---

# Routing

Routing is a core feature of any web framework, and Ghast provides a simple yet composable routing system. The router matches incoming HTTP requests to the appropriate handler based on the request method and path.

You create an application using `ghast.New()`, which gives you a `Ghast` instance with a built-in root router. You can register routes directly on the app, or build separate `Router` instances and mount them at a path prefix using `app.Route()`.

```go
import "github.com/Leonard-Atorough/ghast"

app := ghast.New()

// Register a route directly on the root router
app.Get("/", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"message": "Hello"})
})

// Build a sub-router and mount it under /users
userRouter := ghast.NewRouter()
userRouter.Get("/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"id": r.Param("id")})
})
app.Route("/users", userRouter)

app.Listen(":8080")
```

This registers two routes: `GET /` on the root router and `GET /users/:id` via the mounted sub-router.

---

## Two ways to register routes

### App-level routes

The `Ghast` instance exposes `Get`, `Post`, `Put`, `Delete`, `Patch`, `Head`, and `Options` as direct methods. Routes registered this way live on the internal root router and match against the full request path.

```go
app.Get("/health", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.Plain(200, "ok")
})

app.Post("/login", loginHandler)
app.Delete("/sessions/:id", logoutHandler)
```

Use app-level routes for top-level, standalone endpoints that don't belong to a particular resource group.

### Router-level routes (sub-routers)

For groups of related endpoints, create a `Router` independently with `ghast.NewRouter()`, register routes on it with **prefix-relative paths**, then mount it onto the app with `app.Route()`.

```go
userRouter := ghast.NewRouter()

userRouter.Get("", listUsersHandler)       // matches GET /users
userRouter.Get("/:id", getUserHandler)     // matches GET /users/:id
userRouter.Post("", createUserHandler)     // matches POST /users
userRouter.Put("/:id", updateUserHandler)  // matches PUT /users/:id
userRouter.Delete("/:id", deleteUserHandler)

app.Route("/users", userRouter)
```

When a request arrives, Ghast strips the mount prefix before passing the request to the sub-router. This means routes inside the router are written **without** the prefix — `""` for the root of the group, `"/:id"` for a parameterised child, and so on.

> **Prefix stripping** — if `userRouter` is mounted at `/users`, a request to `GET /users/42` is handed to the router with the path `/42` (not `/users/42`). Write your sub-router paths accordingly.

---

## Route methods

Ghast provides a method for each standard HTTP verb on both the `Ghast` app and the `Router` interface:

```go
router.Get("/path", handler)
router.Post("/path", handler)
router.Put("/path", handler)
router.Delete("/path", handler)
router.Patch("/path", handler)
router.Head("/path", handler)
router.Options("/path", handler)
```

For non-standard methods use `Handle` directly:

```go
router.Handle("PURGE", "/cache/:key", purgeHandler)
```

All verb methods return the router, so calls can be chained:

```go
ghast.NewRouter().
    Get("", listHandler).
    Post("", createHandler).
    Get("/:id", getHandler)
```

---

## Route paths

### Static paths

```go
app.Get("/about", aboutHandler)
app.Get("/api/v1/status", statusHandler)
```

### Dynamic segments

Prefix any segment with `:` to make it a named parameter. The value is available in the handler via `r.Param("name")`.

```go
app.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")   // "42" for a request to /users/42
    w.JSON(200, map[string]string{"id": id})
})
```

Multiple parameters work the same way:

```go
router.Get("/:userId/posts/:postId", func(w ghast.ResponseWriter, r *ghast.Request) {
    userID := r.Param("userId")
    postID := r.Param("postId")
    // ...
})
```

### Query parameters

Query parameters are parsed automatically and accessed via `r.Query("key")`:

```go
// GET /search?q=ghast&page=2
app.Get("/search", func(w ghast.ResponseWriter, r *ghast.Request) {
    q    := r.Query("q")     // "ghast"
    page := r.Query("page")  // "2"
    // ...
})
```

---

## Route handlers

A handler is any value that implements `ghast.Handler` — that is, any type with a `ServeHTTP(ResponseWriter, *Request)` method. The `ghast.HandlerFunc` adapter lets you use ordinary functions directly:

```go
// Inline function literal
app.Get("/hello", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.Plain(200, "Hello, World!")
})

// Named function (must match the HandlerFunc signature)
func getUserHandler(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    user := fetchUser(id)
    if user == nil {
        ghast.Error(w, 404, "user not found")
        return
    }
    w.JSON(200, user)
}

app.Get("/users/:id", getUserHandler)
```

---

## Composable sub-routers

The real power of v0.5.0's routing model is that routers are completely independent values. Create one anywhere, pass it around, and mount it wherever you need. This makes it straightforward to split routing across packages.

```go
// users/router.go
package users

import "github.com/Leonard-Atorough/ghast"

func Router() ghast.Router {
    r := ghast.NewRouter()
    r.Get("", ListHandler)
    r.Post("", CreateHandler)
    r.Get("/:id", GetHandler)
    r.Put("/:id", UpdateHandler)
    r.Delete("/:id", DeleteHandler)
    return r
}

// products/router.go
package products

import "github.com/Leonard-Atorough/ghast"

func Router() ghast.Router {
    r := ghast.NewRouter()
    r.Get("", ListHandler)
    r.Get("/:id", GetHandler)
    return r
}

// main.go
package main

import (
    "log"
    "github.com/Leonard-Atorough/ghast"
    "myapp/users"
    "myapp/products"
)

func main() {
    app := ghast.New()

    app.Get("/", homeHandler)

    app.Route("/users", users.Router())
    app.Route("/products", products.Router())

    if err := app.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

Requests are dispatched by longest prefix match, so:

| Request           | Handled by                         |
| ----------------- | ---------------------------------- |
| `GET /`           | root router                        |
| `GET /users`      | `users.Router()` → `""` route      |
| `GET /users/42`   | `users.Router()` → `/:id` route    |
| `GET /products/7` | `products.Router()` → `/:id` route |

---

## Response methods

| Method                                               | Description                                                              |
| ---------------------------------------------------- | ------------------------------------------------------------------------ |
| `Status(code int) ResponseWriter`                    | Sets the HTTP status code. Chainable.                                    |
| `SetHeader(key, value string) ResponseWriter`        | Sets a response header. Chainable.                                       |
| `JSON(statusCode int, data interface{}) error`       | Marshals `data` as JSON and sends with `Content-Type: application/json`. |
| `JSONPretty(statusCode int, data interface{}) error` | Same as `JSON` but indented for readability.                             |
| `HTML(statusCode int, html string) error`            | Sends an HTML response with `Content-Type: text/html`.                   |
| `Plain(statusCode int, text string) error`           | Sends a plain text response with `Content-Type: text/plain`.             |
| `Send(data []byte) (int, error)`                     | Writes a raw byte slice as the response body.                            |
| `SendString(s string) (int, error)`                  | Writes a string as the response body.                                    |

---

# Middleware

Middleware are functions that wrap handlers to add behaviour before or after the inner handler runs — logging, authentication, timing, header injection, and so on. In Ghast, a middleware has the signature:

```go
type Middleware func(Handler) Handler
```

It receives the next handler in the chain and returns a new handler that wraps it.

```go
loggingMiddleware := func(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        log.Printf("--> %s %s", r.Method, r.Path)
        next.ServeHTTP(w, r)
        log.Printf("<-- %s %s done", r.Method, r.Path)
    })
}
```

---

## Applying middleware

Ghast supports middleware at three distinct scopes:

### Global middleware (`app.Use`)

Applies to every request handled by the application, regardless of which router matches.

```go
app.Use(loggingMiddleware)
app.Use(recoveryMiddleware)
```

### Router-level middleware (`router.Use`)

Applies to all routes registered on a specific router. Mount-scoped middleware can also be passed as extra arguments to `app.Route()`.

```go
// Applied to every route on this router
adminRouter.Use(authMiddleware)

// — or — applied via app.Route for the same effect
app.Route("/admin", adminRouter, authMiddleware)
```

### Per-route middleware

Pass middleware as additional arguments to any route registration call. They apply only to that route.

```go
app.Get("/protected", protectedHandler, authMiddleware)

router.Post("/sensitive", sensitiveHandler, authMiddleware, auditMiddleware)
```

---

## Middleware execution order

When a request arrives, middleware runs in this order:

1. Global middleware (added via `app.Use`), in the order they were registered
2. Router-level middleware (added via `router.Use` or as arguments to `app.Route`)
3. Per-route middleware (passed alongside the handler registration)
4. The route handler itself

```go
app.Use(timingMiddleware)   // 1st
app.Use(loggingMiddleware)  // 2nd

adminRouter.Use(authMiddleware)              // 3rd (for /admin routes only)
adminRouter.Get("/stats", statsHandler, auditMiddleware)  // 4th, then handler
```

---

## HandlerBuilder

When you need to apply an ordered set of middleware to a single handler without affecting a whole router, `HandlerBuilder` provides a fluent API:

```go
handler := ghast.NewHandlerBuilder(myHandler).
    Use(rateLimitMiddleware).
    Use(authMiddleware).
    Build()

app.Get("/api/data", handler)
```

Middleware passed to `Use` is applied in order — the first call wraps outermost, the last wraps closest to the handler.

---

## Complete example

```go
package main

import (
    "log"
    "time"

    "github.com/Leonard-Atorough/ghast"
)

func timingMiddleware(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s — %v", r.Method, r.Path, time.Since(start))
    })
}

func authMiddleware(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        if r.GetHeader("Authorization") == "" {
            ghast.Error(w, 401, "unauthorized")
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    app := ghast.New()

    // Global middleware — runs on every request
    app.Use(timingMiddleware)

    // Public routes on the root router
    app.Get("/", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, map[string]string{"status": "ok"})
    })

    // User sub-router — no auth required
    userRouter := ghast.NewRouter()
    userRouter.Get("", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, []string{"alice", "bob"})
    })
    userRouter.Get("/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, map[string]string{"id": r.Param("id")})
    })
    app.Route("/users", userRouter)

    // Admin sub-router — auth required for all routes
    adminRouter := ghast.NewRouter()
    adminRouter.Get("/stats", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, map[string]int{"users": 100})
    })
    app.Route("/admin", adminRouter, authMiddleware)

    if err := app.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
}
```
