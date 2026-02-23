package examples

import (
	"log"

	ghast "ghast/lib"
)

// BasicExampleApp demonstrates basic Ghast usage with simple routes and JSON responses.
// This is the simplest way to get started with Ghast.
//
// To run this example, rename this file to main.go in a new directory, run `go mod tidy`,
// and then `go run main.go`.
func BasicExampleApp() {
	// Create a new Ghast application
	app := ghast.New()

	// Simple root endpoint that returns JSON
	app.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"message": "Welcome to Ghast!",
			"version": "0.1.0",
		})
	}))

	// Plain text endpoint
	app.Get("/hello", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.Status(200).SetHeader("Content-Type", "text/plain")
		w.SendString("Hello, World!")
	}))

	// Route with path parameter
	app.Get("/users/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		userID := r.Params["id"]
		w.JSON(200, map[string]interface{}{
			"id":   userID,
			"name": "John Doe",
		})
	}))

	// POST endpoint that echoes the request body
	app.Post("/echo", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{
			"received": r.Body,
		})
	}))

	// Start the server on port 8080
	log.Printf("Starting basic example server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
