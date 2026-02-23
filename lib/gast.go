package ghast

const Version = "0.1.0"

// New creates and returns a new Server instance, ready for route registration and listening.
// This is the primary entry point for the Ghast framework.
//
// Usage:
//
//	app := ghast.New()
//	app.Get("/", handler)
//	app.Listen(":8080")
//
// For path-prefixed sub-routers:
//
//	apiRouter := app.NewRouter("/api")
//	apiRouter.Get("/users", handler)
func New() *Server {
	return NewServer()
}
