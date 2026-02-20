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

	// Create the router and register routes
	router := ghast.NewRouter()

	loggingMiddleware := func(next ghast.Handler) ghast.Handler {
		return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
			log.Printf("Request: %s %s", r.Method, r.Path)
			next.ServeHTTP(w, r)
		})
	}

	// Apply global middleware
	router.Use(loggingMiddleware)
	router.Use(middleware.RecoveryMiddleware(middleware.Options{Log: true}))
	router.Use(middleware.RequestIDMiddleware(middleware.RequestIDOptions{HeaderName: "X-Request-ID"}))

	// Register example routes
	router.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "Welcome to Ghast!"})
	}))

	router.Get("/hello", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.Status(200).SetHeader("Content-Type", "text/plain")
		w.SendString("Hello, World!")
	}))

	router.Get("/users/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		userID := r.Params["id"]
		w.JSON(200, map[string]string{"userId": userID})
	}))

	router.Get("/admin/:id/stat/:statId", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		id := r.Params["id"]
		statId := r.Params["statId"]
		w.JSON(200, map[string]string{"adminId": id, "statId": statId})
	}))

	router.Post("/echo", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"received": r.Body})
	}))

	AppRouter := ghast.NewRouter()

	// Create server and start listening
	server := ghast.NewServer().
		AddRouter(ghast.RouterPath{Router: router}).
		AddRouter(ghast.RouterPath{Path: "/app", Router: AppRouter})

	// Validate that we can add a method onto a route after it's been added to the server
	AppRouter.Get("/details", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "Hello from App Router!"})
	}))

	if err := server.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
