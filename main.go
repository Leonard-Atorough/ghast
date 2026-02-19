package main

import (
	"flag"
	"log"

	ghast "ghast/lib"
)

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	// Create the router and register routes
	router := ghast.NewRouter()

	// Apply global middleware
	router.Use(ghast.LoggingMiddleware)
	router.Use(ghast.RecoveryMiddleware)

	// Register example routes
	router.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "Welcome to Ghast!"})
	}))

	router.Get("/hello", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.Status(200).SetHeader("Content-Type", "text/plain")
		w.SendString("Hello, World!")
	}))

	router.Post("/echo", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"received": r.Body})
	}))

	// Create server and start listening
	server := ghast.NewServer(router)
	if err := server.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
