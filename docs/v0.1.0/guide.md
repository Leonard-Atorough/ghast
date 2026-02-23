---
title: "Routing and Middleware"
description: "Learn how to define routes and handlers in Ghast, including support for multiple routers, route parameters, and middleware."
authors: ["Leonardo"]
dateCreated: "2024-02-20"
dateUpdated: "2024-02-20"
---

# Routing

Routing is a core feature of any web framework, and Ghast provides a powerful yet simple routing system. The router matches incoming HTTP requests to the appropriate handler functions based on the request method and path.

You create an application using `ghast.New()`, which gives you a Server instance with a built-in root router. You can then register routes directly on the server or create sub-routers with path prefixes:

```go
import "ghast"

app := ghast.New()

// Register routes on the root router
app.Get("/", func(w ghast.ResponseWriter, r *ghast.Request) {
    w.JSON(200, map[string]string{"message": "Hello"})
})

// Create a sub-router for API routes
apiRouter := app.NewRouter("/api")
apiRouter.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Params["id"]
    // Handle GET /api/users/:id
})

// Start the server
app.Listen(":8080")
```

This creates routes at `/` and `/api/users/:id`, demonstrating both root-level routing and sub-router organization.

## Route methods

Ghast provides convenience methods for each HTTP method (GET, POST, PUT, DELETE, HEAD, OPTIONS) on the `Router` struct. These methods allow you to define routes in a more expressive way:

```go
app.Get("/users/:id", getUserHandler)
app.Post("/users", createUserHandler)
app.Put("/users/:id", updateUserHandler)
app.Delete("/users/:id", deleteUserHandler)
app.Head("/users/:id", headUserHandler)
app.Options("/users", optionsUserHandler)
```

## Route Paths

Route paths can include static segments (e.g., `/users`) and dynamic segments (e.g., `/:id`). Dynamic segments are denoted by a colon (`:`) followed by the parameter name. You can access these parameters in your handler using the `Param` method on the request:

```go
app.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Params["id"]
    // Use the id parameter
})
```

Route paths can also include query parameters, which can be accessed using the `Query` method on the request:

```go
app.Get("/search", func(w ghast.ResponseWriter, r *ghast.Request) {
    q := r.Queries["q"]
    // Use the q query parameter
})
```

This could match a request like `/search?q=ghast` and allow you to access the `q` query parameter in your handler.

## Route Handlers

Route handlers are functions that take a `ResponseWriter` and a `Request` as parameters. They are responsible for processing the incoming request and sending a response back to the client. For example:

```go
func getUserHandler(w ghast.ResponseWriter, r *ghast.Request) {
    id := r.Params["id"]
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

## Example

```go
package main

import (
    "log"
    "ghast"
)

func main() {
    app := ghast.New()
    
    // Root router
    app.Get("/hello", func(w ghast.ResponseWriter, r *ghast.Request) {
        w.JSON(200, map[string]string{"message": "Hello, World!"})
    })
    
    // Sub-router for API routes
    apiRouter := app.NewRouter("/api")
    apiRouter.Get("/users/:id", func(w ghast.ResponseWriter, r *ghast.Request) {
        id := r.Params["id"]
        user := map[string]string{"id": id, "name": "Alice"}
        w.JSON(200, user)
    })
    
    // Start the server
    if err := app.Listen(":8080"); err != nil {
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

## Sub-Routers

Ghast supports organizing routes into sub-routers with path prefixes, allowing you to group related routes by functionality or module. Use `server.NewRouter(prefix)` to create a sub-router:

```go
// Root router
app.Get("/", homeHandler)

// API routes with /api prefix
apiRouter := app.NewRouter("/api")
apiRouter.Get("/users", getUsersHandler)
apiRouter.Get("/users/:id", getUserByIDHandler)
apiRouter.Post("/users", createUserHandler)

// Admin routes with /admin prefix
adminRouter := app.NewRouter("/admin")
adminRouter.Get("/dashboard", adminDashboardHandler)
adminRouter.Get("/users", adminUsersHandler)
```

Requests are routed based on the longest matching prefix, so:
- `GET /` → handled by root router
- `GET /api/users` → handled by apiRouter
- `GET /admin/dashboard` → handled by adminRouter

Sub-routers can be defined in separate packages for better modularity. Simply create a function that takes a Router and registers routes on it:

```go
// main.go
package main
import (
    "log"
    "ghast"
    "myapp/api"
    "myapp/admin"
)

func main() {
    app := ghast.New()
    app.Get("/", homeHandler)
    
    // Create and initialize sub-routers from packages
    apiRouter := app.NewRouter("/api")
    api.Init(apiRouter)
    
    adminRouter := app.NewRouter("/admin")
    admin.Init(adminRouter)
    
    if err := app.Listen(":8080"); err != nil {
        log.Fatal(err)
    }
}

// api/router.go
package api
import "ghast"

func Init(router ghast.Router) {
    router.Get("/users", getUsersHandler)
    router.Post("/users", createUserHandler)
    router.Get("/users/:id", getUserByIDHandler)
}

// admin/router.go
package admin
import "ghast"

func Init(router ghast.Router) {
    router.Get("/dashboard", dashboardHandler)
    router.Get("/users", adminUsersHandler)
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

You can then apply this middleware at different levels:

```go
// Apply middleware to the server (affects all routers)
app.Use(customMiddleware)

// Apply middleware to a specific router (affects all routes on that router)
apiRouter.Use(customMiddleware)

// Apply middleware to a specific route (affects only that route)
apiRouter.Get("/users/:id", handler, customMiddleware)  // Pass middleware as optional parameters
```

### A note on the order of middleware execution

When a request comes in, the server will first match the request path to the appropriate router based on the longest prefix match. Once the router is selected, the server will execute any server-level middleware in the reverse order they were added, then any router-level middleware for that router, and finally any route-specific middleware for the matched route. After all middleware has been executed, the final route handler will be called to generate the response. This allows you to have a flexible and powerful middleware system that can be applied at different levels of your application architecture. Example:

```go
// Server-level middleware
app.Use(loggingMiddleware)
app.Use(recoveryMiddleware)

// Router-level middleware
apiRouter.Use(authMiddleware)

// Route-specific middleware
apiRouter.Get("/users/:id", userHandler, userMiddleware)

// When a request comes in for /api/users/123, the execution order will be:
// 1. Recovery middleware (server-level)
// 2. Logging middleware (server-level)
// 3. Auth middleware (router-level)
// 4. User middleware (route-specific)
```

Middleware gives us a way to access and modify the request and response objects at different stages of the request handling process, allowing us to implement cross-cutting concerns like logging, authentication, and error handling in a clean and modular way.
