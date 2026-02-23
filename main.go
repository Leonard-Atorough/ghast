package main

import (
	"flag"
	"log"

	ghast "ghast/lib"
	middleware "ghast/middleware"
)

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	// Create the app and register routes
	app := ghast.New()

	loggingMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			log.Printf("Request: %s %s", r.Method, r.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Apply global middleware
	app.Use(loggingMiddleware)
	app.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))
	app.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{HeaderName: "X-Request-ID"}))

	// Register example routes
	app.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "Welcome to Ghast!"})
	}))

	app.Get("/hello", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.Status(200).SetHeader("Content-Type", "text/plain")
		w.SendString("Hello, World!")
	}))

	app.Get("/users/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		userID := r.Params["id"]
		w.JSON(200, map[string]string{"userId": userID})
	}))

	app.Get("/admin/:id/stat/:statId", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		id := r.Params["id"]
		statId := r.Params["statId"]
		w.JSON(200, map[string]string{"adminId": id, "statId": statId})
	}))

	app.Post("/echo", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"received": r.Body})
	}))

	// Create a sub-router with path prefix
	appRouter := app.NewRouter("/app")

	// Validate that we can add a method onto a route after it's been added to the server
	appRouter.Get("/details", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "Hello from App Router!"})
	}))

	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
