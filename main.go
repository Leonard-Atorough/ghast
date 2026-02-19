package main

import (
	"flag"
	"log"

	gust "gust/lib"
)

func main() {
	// Parse command-line flags
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	// Create the router and register routes
	router := gust.NewRouter()

	// Apply global middleware
	router.Use(gust.LoggingMiddleware)
	router.Use(gust.RecoveryMiddleware)

	// Register example routes
	router.Get("/", gust.HandlerFunc(func(w gust.ResponseWriter, r *gust.Request) {
		w.JSON(200, map[string]string{"message": "Welcome to Gust!"})
	}))

	router.Get("/hello", gust.HandlerFunc(func(w gust.ResponseWriter, r *gust.Request) {
		w.Status(200).SetHeader("Content-Type", "text/plain")
		w.SendString("Hello, World!")
	}))

	router.Post("/echo", gust.HandlerFunc(func(w gust.ResponseWriter, r *gust.Request) {
		w.JSON(200, map[string]string{"received": r.Body})
	}))

	// Create server and start listening
	server := gust.NewServer(router)
	if err := server.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
