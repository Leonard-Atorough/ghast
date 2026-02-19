package examples

import (
	"encoding/json"
	"fmt"

	ghast "ghast/lib"
)

// HelloHandler responds with a simple "Hello, World!" message.
type HelloHandler struct{}

func (h *HelloHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	rw.Status(200).SetHeader("Content-Type", "text/plain")
	rw.SendString("Hello, World!")
}

// JSONHandler responds with JSON data.
type JSONHandler struct{}

func (h *JSONHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	data := map[string]interface{}{
		"message": "This is JSON",
		"status":  "success",
		"code":    200,
	}
	rw.JSON(200, data)
}

// EchoHandler echoes back the request body.
type EchoHandler struct{}

func (h *EchoHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	if req.Body == "" {
		rw.JSON(400, map[string]string{"error": "no body provided"})
		return
	}
	rw.JSON(200, map[string]string{"echoed": req.Body})
}

// QueryParamHandler demonstrates reading query parameters.
type QueryParamHandler struct{}

func (h *QueryParamHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	name := req.Query("name")
	if name == "" {
		name = "Guest"
	}
	message := fmt.Sprintf("Hello, %s!", name)
	rw.JSON(200, map[string]string{"message": message})
}

// User represents a user resource in the API.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserHandler demonstrates a typical REST API handler for retrieving a user.
type UserHandler struct{}

func (h *UserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	// Simulated database - in real code, this would be actual data access
	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}
	rw.JSON(200, user)
}

// CreateUserHandler demonstrates parsing JSON and creating a resource.
type CreateUserHandler struct{}

func (h *CreateUserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	var user User
	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		ghast.Error(rw, 400, "Invalid JSON")
		return
	}

	// Simulated database save
	user.ID = 2

	rw.JSON(201, user)
}

// NotFoundHandler responds with a 404 error.
type NotFoundHandler struct{}

func (h *NotFoundHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	rw.Status(404)
	rw.SendString("404 Not Found")
}

// HealthCheckHandler is a simple health check endpoint.
type HealthCheckHandler struct{}

func (h *HealthCheckHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	rw.JSON(200, map[string]string{
		"status": "healthy",
	})
}
