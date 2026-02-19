package examples

import (
	"fmt"
	"strconv"

	ghast "ghast/lib"
)

// ===== Basic Routing Examples =====

// IndexHandler responds to GET /
type IndexHandler struct{}

func (h *IndexHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	rw.JSON(200, map[string]string{
		"message": "Welcome to the API",
		"version": "1.0",
	})
}

// ===== Parameterized Routes Examples =====

// GetUserHandler handles requests to GET /users/:id
type GetUserHandler struct{}

func (h *GetUserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	userID := req.Params["id"]
	if userID == "" {
		rw.Status(400).JSON(400, map[string]string{"error": "user ID is required"})
		return
	}

	// In a real application, you would fetch this from a database
	rw.JSON(200, map[string]interface{}{
		"id":    userID,
		"name":  "John Doe",
		"email": "john@example.com",
	})
}

// ===== Multiple Parameters Example =====

// UserPostHandler handles GET /users/:userId/posts/:postId
type UserPostHandler struct{}

func (h *UserPostHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	userID := req.Params["userId"]
	postID := req.Params["postId"]

	if userID == "" || postID == "" {
		rw.Status(400).JSON(400, map[string]string{"error": "userId and postId are required"})
		return
	}

	rw.JSON(200, map[string]interface{}{
		"userId":  userID,
		"postId":  postID,
		"title":   "Sample Post",
		"content": "This is a sample post from user " + userID,
	})
}

// ===== List Resources Example =====

// UserListHandler handles GET /users with optional query parameters
type UserListHandler struct{}

func (h *UserListHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	// Example query parameters: ?page=1&limit=10&role=admin
	pageStr := req.Query("page")
	limitStr := req.Query("limit")
	role := req.Query("role")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil {
		limit = l
	}

	rw.JSON(200, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"role":  role,
		"users": []map[string]interface{}{
			{"id": "1", "name": "Alice", "role": "admin"},
			{"id": "2", "name": "Bob", "role": "user"},
		},
	})
}

// ===== Create Resource Example =====

// PostUserHandler handles POST /users
type PostUserHandler struct{}

func (h *PostUserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	type UserInput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var user UserInput
	if err := req.JSON(&user); err != nil {
		rw.Status(400).JSON(400, map[string]string{"error": "invalid JSON"})
		return
	}

	if user.Name == "" || user.Email == "" {
		rw.Status(400).JSON(400, map[string]string{"error": "name and email are required"})
		return
	}

	rw.Status(201).JSON(201, map[string]interface{}{
		"id":    "100",
		"name":  user.Name,
		"email": user.Email,
	})
}

// ===== Update Resource Example =====

// PutUserHandler handles PUT /users/:id
type PutUserHandler struct{}

func (h *PutUserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	userID := req.Params["id"]

	type UserUpdate struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var update UserUpdate
	if err := req.JSON(&update); err != nil {
		rw.Status(400).JSON(400, map[string]string{"error": "invalid JSON"})
		return
	}

	rw.JSON(200, map[string]interface{}{
		"id":      userID,
		"name":    update.Name,
		"email":   update.Email,
		"message": fmt.Sprintf("User %s updated successfully", userID),
	})
}

// ===== Delete Resource Example =====

// DeleteUserHandler handles DELETE /users/:id
type DeleteUserHandler struct{}

func (h *DeleteUserHandler) ServeHTTP(rw ghast.ResponseWriter, req *ghast.Request) {
	userID := req.Params["id"]

	rw.Status(204).SetHeader("Content-Type", "application/json")
	rw.SendString(fmt.Sprintf(`{"message":"User %s deleted successfully"}`, userID))
}

// ===== Middleware Examples =====

// LoggingMiddleware logs request information
func LoggingMiddleware(next ghast.Handler) ghast.Handler {
	return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		fmt.Printf("Request: %s %s\n", r.Method, r.Path)
		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware checks for authorization header (simple example)
func AuthenticationMiddleware(next ghast.Handler) ghast.Handler {
	return ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		token := r.Headers["Authorization"]
		if token == "" {
			w.Status(401).JSON(401, map[string]string{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ===== Complete Router Setup Example =====

// SetupRouter creates and configures a router with all example routes
func SetupRouter() ghast.Router {
	router := ghast.NewRouter()

	// Global middleware
	router.Use(LoggingMiddleware)

	// Basic routes
	router.Get("/", &IndexHandler{})

	// User CRUD routes
	router.Get("/users", &UserListHandler{})
	router.Post("/users", &PostUserHandler{})
	router.Get("/users/:id", &GetUserHandler{})
	router.Put("/users/:id", &PutUserHandler{})
	router.Delete("/users/:id", &DeleteUserHandler{})

	// Nested parameter routes
	router.Get("/users/:userId/posts/:postId", &UserPostHandler{})

	// Protected routes with path-specific middleware
	router.UsePath("/admin/stats", AuthenticationMiddleware)
	router.Get("/admin/stats", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]interface{}{
			"totalUsers": 42,
			"totalPosts": 156,
		})
	}))

	return router
}
