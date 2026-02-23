package examples

import (
	"log"

	"github.com/Leonard-Atorough/ghast"
)

// GroupedRoutesExampleApp demonstrates how to organize routes into logical groups
// using route prefixes. This is useful for organizing endpoints by feature or API version.
func GroupedRoutesExampleApp() {
	app := ghast.New()

	// Root endpoint
	app.Get("/", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]string{"message": "API Root"})
	}))

	// ===== User Routes =====
	// Create a router for all user-related endpoints, then mount it at /users
	userRouter := ghast.NewRouter()

	userRouter.Get("", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]interface{}{
			"users": []map[string]string{
				{"id": "1", "name": "Alice"},
				{"id": "2", "name": "Bob"},
			},
		})
	}))

	userRouter.Get("/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		userID := r.Params["id"]
		w.JSON(200, map[string]interface{}{
			"id":    userID,
			"name":  "John Doe",
			"email": "john@example.com",
		})
	}))

	userRouter.Post("/:id/update", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		userID := r.Params["id"]
		w.JSON(200, map[string]string{
			"message": "User " + userID + " updated",
		})
	}))

	// Mount the user router under /users
	app.Route("/users", userRouter)

	// ===== Product Routes =====
	// Create a router group for all product-related endpoints at /products
	// This also shows how you can define the route first then associate it with a base path later
	productRouter := ghast.NewRouter()

	productRouter.Get("", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]interface{}{
			"products": []map[string]string{
				{"id": "1", "name": "Laptop"},
				{"id": "2", "name": "Mouse"},
			},
		})
	}))

	productRouter.Get("/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		productID := r.Params["id"]
		w.JSON(200, map[string]interface{}{
			"id":    productID,
			"name":  "Laptop Pro",
			"price": 1299.99,
		})
	}))

	// Mount the product router under /products
	app.Route("/products", productRouter)

	// ===== Admin Routes =====
	// Create a router for admin-only endpoints, then mount it at /admin
	adminRouter := ghast.NewRouter()

	adminRouter.Get("/stats", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		w.JSON(200, map[string]interface{}{
			"totalUsers":    1000,
			"totalProducts": 500,
		})
	}))

	adminRouter.Get("/stats/:id", ghast.HandlerFunc(func(w ghast.ResponseWriter, r *ghast.Request) {
		statID := r.Params["id"]
		w.JSON(200, map[string]interface{}{
			"id":    statID,
			"value": 42,
		})
	}))

	// Mount the admin router under /admin
	app.Route("/admin", adminRouter)

	// Start the server on port 8081
	log.Printf("Starting grouped routes example on :8081")
	log.Printf("Try: GET /users, GET /users/123, GET /products, GET /admin/stats, etc.")
	if err := app.Listen(":8081"); err != nil {
		log.Fatal(err)
	}
}
