---
title: "Routing and Middleware"
description: "Learn how to define routes and handlers in Ghast, including support for multiple routers, route parameters, and middleware."
authors: ["Leonardo"]
dateCreated: "2024-02-20"
dateUpdated: "2024-02-20"
---

# Routing

Routing is a core feature of any web framework, and Ghast provides a powerful yet simple routing system. The router matches incoming HTTP requests to the appropriate handler functions based on the request method and path.

you define routes using the methods on the `Router` struct. For example:

```go
import "github.com/leonardo/ghast"
server := ghast.NewServer()

router := ghast.NewRouter()
router.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    // Handle GET /users/:id
})

server.AddRouter(ghast.RouterPath{Path: "/api", Router: router})
```

This sets up a route that matches GET requests to `/api/users/:id`, where `:id` is a path parameter that can be accessed in the handler.

### App level routing

Ghast also supports multiple routers on the same server, allowing you to organize your routes into different groups (e.g., API routes, admin routes). You can add multiple routers to the server with different path prefixes:

```go
apiRouter := ghast.NewRouter()
apiRouter.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    // Handle GET /api/users/:id
})
adminRouter := ghast.NewRouter()
adminRouter.Get("/dashboard", func(w ghast.ResponseWriter, r *ghast.Request) {
    // Handle GET /admin/dashboard
})
server := ghast.NewServer().
    AddRouter(ghast.RouterPath{Path: "/api", Router: apiRouter}).
    AddRouter(ghast.RouterPath{Path: "/admin", Router: adminRouter})
```

## Route methods

Ghast provides convenience methods for each HTTP method (GET, POST, PUT, DELETE, HEAD, OPTIONS) on the `Router` struct. These methods allow you to define routes in a more expressive way:

```go
router.Get("/users/:id", getUserHandler)
router.Post("/users", createUserHandler)
router.Put("/users/:id", updateUserHandler)
router.Delete("/users/:id", deleteUserHandler)
router.Head("/users/:id", headUserHandler)
router.Options("/users", optionsUserHandler)
```

There is also a generic `Handle` method that allows you to specify the method as a string:

```go
router.Handle("GET", "/users/:id", getUserHandler)
```

## Route Paths

Route paths can include static segments (e.g., `/users`) and dynamic segments (e.g., `/:id`). Dynamic segments are denoted by a colon (`:`) followed by the parameter name. You can access these parameters in your handler using the `Param` method on the request:

```go
router.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    // Use the id parameter
})
```

Route paths can also include query parameters, which can be accessed using the `Query` method on the request:

```go
router.Get("/search", func(w ghast.ResponseWriter, r *ghast.Request) {
    q := r.Query("q")
    // Use the q query parameter
})
```

This could match a request like `/search?q=ghast` and allow you to access the `q` query parameter in your handler.

## Route Handlers

Route handlers are functions that take a `ResponseWriter` and a `Request` as parameters. They are responsible for processing the incoming request and sending a response back to the client. For example:

```go
func getUserHandler(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    // Fetch user from database using id
    user := getUserFromDB(id)
    if user == nil {
        w.Status(404).JSON(map[string]string{"error": "User not found"})
        return
    }
    w.JSON(user)
}
```

In this example, the `getUserHandler` function retrieves the `id` parameter from the request, fetches the user from the database, and returns a JSON response. If the user is not found, it returns a 404 error.

You can also define handlers inline when setting up routes:

```go
router.Get("/hello", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"message": "Hello, World!"})
})
```

> Note: Ghast 0.1.0 does not support error handling in handlers, so you need to handle errors manually and send appropriate responses.

> Ghast 0.1.0 does not support route groups or nested routes, so all routes are defined at the top level of the router.

> Ghast 0.1.0 doesn't yet support middleware on handlers, but this is planned for a future release. Currently, you can apply middleware at the server level or router level, but not on individual handlers.

## Example

```go
package main

import (
    "log"
    "github.com/leonardo/ghast"
)

func main() {
    router := ghast.NewRouter()
    router.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
        id := r.Param("id")
        user := map[string]string{"id": id, "name": "Alice"}
        w.JSON(200, user)
    })
    server := ghast.NewServer().AddRouter(ghast.RouterPath{Path: "/api", Router: router})
    if err := server.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

## Response Methods

Ghast's `ResponseWriter` interface includes several convenience methods for sending responses and terminating the request. These methods allow you to set status codes, headers, and send JSON or HTML responses in a more expressive way.

| Method                                               | Description                                                                                                                                                    |
| ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Status(code int) ResponseWriter`                    | Sets the HTTP status code for the response. This method is chainable, allowing you to set the status and then send a response in one line.                     |
| `SetHeader(key, value string) ResponseWriter`        | Sets a header key-value pair for the response. This method is also chainable, so you can set multiple headers before sending the response.                     |
| `JSON(statusCode int, data interface{}) error`       | Encodes the given data as JSON and sends it as the response body with the specified status code. It also sets the `Content-Type` header to `application/json`. |
| `HTML(statusCode int, html string) error`            | Sends the given HTML string as the response body with the specified status code. It sets the `Content-Type` header to `text/html`.                             |
| `Send(data []byte) error`                            | Sends the given byte slice as the response body with the current status code and headers.                                                                      |
| `SendString(data string) error`                      | Sends the given string as the response body with the current status code and headers.                                                                          |
| `JSONPretty(statusCode int, data interface{}) error` | Similar to `JSON()`, but formats the JSON with indentation for better readability.                                                                             |

## App Router

Ghast's `Server` struct supports multiple routers, allowing you to organize your routes into different groups (e.g., API routes, admin routes). You can add multiple routers to the server with different path prefixes:

```go
apiRouter := ghast.NewRouter()
apiRouter.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Param("id")
    // Handle GET /api/users/:id
})

adminRouter := ghast.NewRouter()
adminRouter.Get("/dashboard", func(w ghast.ResponseWriter, r *ghast.Request) {
    // Handle GET /admin/dashboard
})

server := ghast.NewServer().
    AddRouter(ghast.RouterPath{Path: "/api", Router: apiRouter}).
    AddRouter(ghast.RouterPath{Path: "/admin", Router: adminRouter})
```

In this example, we create two routers: `apiRouter` for API routes and `adminRouter` for admin routes. We then add both routers to the server with different path prefixes (`/api` and `/admin`). The server will route incoming requests to the appropriate router based on the longest matching path prefix.

The `AddRouter` method allows you to easily organize your routes and keep related routes together in separate routers. This is especially useful for larger applications where you want to group routes by functionality or module. The `AddRouter` method also returns the server instance, allowing you to chain multiple calls to add routers in a fluent style as shown in the example.

A router can be defined in a separate package and imported into the main server setup, allowing for better modularity and separation of concerns in your application architecture. For example, you could have an `api` package that defines the `apiRouter` and an `admin` package that defines the `adminRouter`, and then import both into your main application to set up the server.

```go
// main.go
import (
    "log"
    "github.com/leonardo/ghast"
    "myapp/api"
    "myapp/admin"
)

func main() {
    server := ghast.NewServer().
        AddRouter(ghast.RouterPath{Path: "/api", Router: api.NewRouter()}).
        AddRouter(ghast.RouterPath{Path: "/admin", Router: admin.NewRouter()})
    if err := server.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
}

// In api/router.go
package api
import "github.com/leonardo/ghast"
func NewRouter() ghast.Router {
    router := ghast.NewRouter()
    router.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
        id := r.Param("id")
        user := map[string]string{"id": id, "name": "Alice"}
        w.JSON(200, user)
    })
    return router
}

// In admin/router.go
package admin
import "github.com/leonardo/ghast"
func NewRouter() ghast.Router {
    router := ghast.NewRouter()
    router.Get("/dashboard", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.HTML(200, "<h1>Admin Dashboard</h1>")
    })
    return router
}
```

# Middleware

For those more fammiliar with Express JS, ghast also supports middleware functions. Middleware functions are functions that have access to the request and response objects, and can perform operations before or after the main route handler is executed. You can apply middleware at the server level (affecting all routers) or at the router level (affecting only routes defined on that router). Middleware can be used for tasks such as logging, authentication, error handling, and more.

Ghast middleware functions are declared a little differently than Express middleware. In Ghast, a middleware function is a function that takes a `Handler` and returns a new `Handler`. This allows you to wrap the original handler with additional functionality. This handler can be seen as the equivalent of the `next` function in Express middleware, allowing you to call the next handler in the chain after performing your middleware logic. For example, here is how you can create a simple logging middleware that logs the request path before calling the next handler:

```go
func customMiddleware(next ghast.Handler) ghast.Handler {
    return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
        // Before handler
        log.Println("Request:", r.Path)
        // Call handler
        next.ServeHTTP(w, r)
        // After handler
        log.Println("Response sent")
    })
}
```

You can then apply this middleware to the server or to a specific router:

```go
// Apply middleware to the server (affects all routers)
server.Use(customMiddleware)
// Apply middleware to a specific router
router.Use(customMiddleware)
// You can also apply route-specific middleware
router.UsePath("/users/:id", customMiddleware) // This applies the middleware only to routes that match /users/:id
```

### A note on the order of middleware execution

When a request comes in, the server will first match the request path to the appropriate router based on the longest prefix match. Once the router is selected, the server will execute any server-level middleware in the reverse order they were added, then any router-level middleware for that router, and finally any route-specific middleware for the matched route. After all middleware has been executed, the final route handler will be called to generate the response. This allows you to have a flexible and powerful middleware system that can be applied at different levels of your application architecture. Example:

```go
// Server-level middleware
server.Use(loggingMiddleware)
server.Use(recoveryMiddleware)
// Router-level middleware
router.Use(authMiddleware)
// Route-specific middleware
router.UsePath("/users/:id", userMiddleware)

// When a request comes in for /users/123, the execution order will be:
// 1. Recovery middleware (server-level)
// 2. Logging middleware (server-level)
// 3. Auth middleware (router-level)
// 4. User middleware (route-specific)
```

Middleware gives us a way to access and modify the request and response objects at different stages of the request handling process, allowing us to implement cross-cutting concerns like logging, authentication, and error handling in a clean and modular way.
