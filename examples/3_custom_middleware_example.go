package examples

import (
	"log"
	"time"

	"github.com/Leonard-Atorough/ghast"
)

// CustomMiddlewareExampleApp demonstrates building and using custom middleware
// to add cross-cutting concerns like logging, timing, and authentication.
func CustomMiddlewareExampleApp() {
	app := ghast.New()

	// ===== Custom Middleware: Request Logging =====
	// This middleware logs every request with method, path, and response status
	loggingMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			log.Printf("[REQUEST] %s %s - Starting", r.Method, r.Path)
			next.ServeHTTP(w, r)
			log.Printf("[REQUEST] %s %s - Completed", r.Method, r.Path)
		})
	}

	// ===== Custom Middleware: Request Timing =====
	// This middleware measures how long it takes to process each request
	timingMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			elapsed := time.Since(start)
			log.Printf("[TIMING] %s %s took %v", r.Method, r.Path, elapsed)
		})
	}

	// ===== Custom Middleware: Custom Header Injection =====
	// This middleware adds a custom header to all responses
	headerMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			w.SetHeader("X-Powered-By", "Ghast Framework")
			next.ServeHTTP(w, r)
		})
	}

	// ===== Custom Middleware: Simple Auth Check =====
	// This middleware checks for an authorization header
	authMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			// In a real app, validate the token properly
			if r.Headers["Authorization"] == "" {
				log.Printf("[AUTH] Missing authorization header for %s %s", r.Method, r.Path)
				w.JSON(401, map[string]string{"error": "Missing authorization header"})
				return
			}
			log.Printf("[AUTH] Authorization header present: %s", r.Headers["Authorization"])
			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware globally - these will run for ALL requests
	app.Use(loggingMiddleware)
	app.Use(timingMiddleware)
	app.Use(headerMiddleware)

	// Public endpoint - doesn't require auth
	app.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"message": "Public endpoint - no auth required",
		})
	}))

	// Public health check
	app.Get("/health", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"status": "healthy",
		})
	}))

	// Protected endpoint - requires authentication
	// We wrap it with authMiddleware in addition to the global ones
	protectedHandler := ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"message": "This is a protected endpoint",
			"data":    "Secret information here",
		})
	})
	app.Get("/protected", protectedHandler, authMiddleware)

	// Another protected endpoint
	app.Post("/api/users", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"message": "User created successfully",
		})
	}), authMiddleware)

	// Example of building a handler with multiple middleware using HandlerBuilder
	builder := ghast.NewHandlerBuilder(ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]any{
			"message":     "Built with HandlerBuilder",
			"timestamp":   time.Now().Format(time.RFC3339),
			"request_url": r.Path,
		})
	}))

	// Add multiple middleware to this handler specifically
	builder.Use(authMiddleware)
	builtHandler := builder.Build()

	app.Get("/builder-example", builtHandler)

	log.Printf("Starting custom middleware example on :8082")
	log.Printf("")
	log.Printf("Try these endpoints:")
	log.Printf("  GET  /                    - Public endpoint")
	log.Printf("  GET  /health              - Health check")
	log.Printf("  GET  /protected           - Requires auth header")
	log.Printf("  POST /api/users           - Requires auth header")
	log.Printf("  GET  /builder-example     - Built with HandlerBuilder, requires auth")
	log.Printf("")
	log.Printf("Test with auth: curl -H 'Authorization: Bearer token123' http://localhost:8082/protected")

	if err := app.Listen(":8082"); err != nil {
		log.Fatal(err)
	}
}
