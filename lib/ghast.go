package ghast


type Ghast struct {
	// Router is the core component that manages route registration and request handling.
	// It matches incoming requests to the appropriate handlers based on the HTTP method and path.
	Router *Router
	// Server is responsible for managing TCP connections, listening for incoming requests, and delegating request processing to the Router.
	Server *Server
}

// NewGhast creates a new Ghast instance with an initialized Router and Server.
func NewGhast() *Ghast {
	router := NewRouter()
	server := NewServer()
	return &Ghast{
		Router: &router,
		Server: server,
	}
}