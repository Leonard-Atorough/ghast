package ghast

import (
	"bytes"
	"testing"
)

// TestRouterExactPathMatching tests that exact paths are matched correctly.
func TestRouterExactPathMatching(t *testing.T) {
	router := NewRouter()
	handlerCalled := false

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		handlerCalled = true
	})

	router.Get("/users", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/users", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if !handlerCalled {
		t.Error("exact path handler was not called")
	}
}

// TestRouterParameterizedPathMatching tests that routes with parameters are matched.
func TestRouterParameterizedPathMatching(t *testing.T) {
	router := NewRouter()
	handlerCalled := false

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		handlerCalled = true
	})

	router.Get("/users/:id", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/users/123", Headers: make(map[string]string), Params: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if !handlerCalled {
		t.Error("parameterized path handler was not called")
	}
}

// TestRouterParameterExtraction tests that route parameters are correctly extracted.
func TestRouterParameterExtraction(t *testing.T) {
	router := NewRouter()
	var extractedParams map[string]string

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		extractedParams = r.Params
	})

	router.Get("/users/:id", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/users/123", Headers: make(map[string]string), Params: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if extractedParams == nil {
		t.Fatal("params not set")
	}

	if id, ok := extractedParams["id"]; !ok || id != "123" {
		t.Errorf("parameter extraction failed: got %v, want map[id:123]", extractedParams)
	}
}

// TestRouterMultipleParameters tests parameter extraction with multiple parameters.
func TestRouterMultipleParameters(t *testing.T) {
	router := NewRouter()
	var extractedParams map[string]string

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		extractedParams = r.Params
	})

	router.Get("/users/:userId/posts/:postId", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{
		Method:  "GET",
		Path:    "/users/456/posts/789",
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}

	router.ServeHTTP(rw, req)

	if extractedParams == nil {
		t.Fatal("params not set")
	}

	if userId, ok := extractedParams["userId"]; !ok || userId != "456" {
		t.Errorf("userId parameter extraction failed: got %v", extractedParams)
	}

	if postId, ok := extractedParams["postId"]; !ok || postId != "789" {
		t.Errorf("postId parameter extraction failed: got %v", extractedParams)
	}
}

// TestRouterParameterIsolation tests that parameters from one route don't affect another.
func TestRouterParameterIsolation(t *testing.T) {
	router := NewRouter()
	var firstParams, secondParams map[string]string

	handler1 := HandlerFunc(func(w ResponseWriter, r *Request) {
		firstParams = r.Params
	})

	handler2 := HandlerFunc(func(w ResponseWriter, r *Request) {
		secondParams = r.Params
	})

	router.Get("/users/:id", handler1)
	router.Get("/posts/:postId", handler2)

	// First request to /users/:id route
	mockConn1 := &MockConnection{}
	rw1 := NewResponseWriter(mockConn1)
	req1 := &Request{
		Method:  "GET",
		Path:    "/users/100",
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}
	router.ServeHTTP(rw1, req1)

	// Second request to /posts/:postId route
	mockConn2 := &MockConnection{}
	rw2 := NewResponseWriter(mockConn2)
	req2 := &Request{
		Method:  "GET",
		Path:    "/posts/200",
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}
	router.ServeHTTP(rw2, req2)

	if val, ok := firstParams["id"]; !ok || val != "100" {
		t.Errorf("first route params incorrect: %v", firstParams)
	}

	if val, ok := secondParams["postId"]; !ok || val != "200" {
		t.Errorf("second route params incorrect: %v", secondParams)
	}
}

// TestRouter404NotFound tests that non-existent routes return 404.
func TestRouter404NotFound(t *testing.T) {
	router := NewRouter()

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/nonexistent", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	output := mockConn.writeBuffer.String()
	if !bytes.Contains([]byte(output), []byte("404")) {
		t.Error("404 response not found in output")
	}
}

// TestRouterExactPathPriority tests that exact paths are matched before regex routes.
func TestRouterExactPathPriority(t *testing.T) {
	router := NewRouter()
	var whichHandler string

	exactHandler := HandlerFunc(func(w ResponseWriter, r *Request) {
		whichHandler = "exact"
	})

	paramHandler := HandlerFunc(func(w ResponseWriter, r *Request) {
		whichHandler = "param"
	})

	// Register both exact and parameterized versions
	router.Get("/users/me", exactHandler)
	router.Get("/users/:id", paramHandler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/users/me", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if whichHandler != "exact" {
		t.Errorf("exact path should take priority: got %s", whichHandler)
	}
}

// TestRouterHTTPMethods tests different HTTP verbs register correctly.
func TestRouterHTTPMethods(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		register func(Router, Handler)
	}{
		{"GET", "GET", func(r Router, h Handler) { r.Get("/test", h) }},
		{"POST", "POST", func(r Router, h Handler) { r.Post("/test", h) }},
		{"PUT", "PUT", func(r Router, h Handler) { r.Put("/test", h) }},
		{"DELETE", "DELETE", func(r Router, h Handler) { r.Delete("/test", h) }},
		{"PATCH", "PATCH", func(r Router, h Handler) { r.Patch("/test", h) }},
		{"HEAD", "HEAD", func(r Router, h Handler) { r.Head("/test", h) }},
		{"OPTIONS", "OPTIONS", func(r Router, h Handler) { r.Options("/test", h) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter()
			handlerCalled := false

			handler := HandlerFunc(func(w ResponseWriter, r *Request) {
				handlerCalled = true
			})

			tt.register(router, handler)

			mockConn := &MockConnection{}
			rw := NewResponseWriter(mockConn)
			req := &Request{Method: tt.method, Path: "/test", Headers: make(map[string]string)}

			router.ServeHTTP(rw, req)

			if !handlerCalled {
				t.Errorf("%s handler was not called", tt.method)
			}
		})
	}
}

// TestRouterGlobalMiddleware tests that global middleware is applied to all routes.
func TestRouterGlobalMiddleware(t *testing.T) {
	router := NewRouter()
	middlewareCalled := false

	middleware := func(next Handler) Handler {
		return HandlerFunc(func(w ResponseWriter, r *Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	}

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {})

	router.Use(middleware)
	router.Get("/test", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/test", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if !middlewareCalled {
		t.Error("global middleware was not applied")
	}
}

// TestRouterPathSpecificMiddleware tests that path-specific middleware applies only to specified paths.
func TestRouterPathSpecificMiddleware(t *testing.T) {
	router := NewRouter()
	var middlewarePaths []string

	middleware := func(path string) Middleware {
		return func(next Handler) Handler {
			return HandlerFunc(func(w ResponseWriter, r *Request) {
				middlewarePaths = append(middlewarePaths, path)
				next.ServeHTTP(w, r)
			})
		}
	}

	handler1 := HandlerFunc(func(w ResponseWriter, r *Request) {})
	handler2 := HandlerFunc(func(w ResponseWriter, r *Request) {})

	router.UsePath("/users", middleware("/users"))
	router.Get("/users", handler1)
	router.Get("/posts", handler2)

	// Request to /users (should trigger middleware)
	mockConn1 := &MockConnection{}
	rw1 := NewResponseWriter(mockConn1)
	req1 := &Request{Method: "GET", Path: "/users", Headers: make(map[string]string)}
	router.ServeHTTP(rw1, req1)

	// Request to /posts (should not trigger middleware)
	mockConn2 := &MockConnection{}
	rw2 := NewResponseWriter(mockConn2)
	req2 := &Request{Method: "GET", Path: "/posts", Headers: make(map[string]string)}
	router.ServeHTTP(rw2, req2)

	if len(middlewarePaths) != 1 || middlewarePaths[0] != "/users" {
		t.Errorf("path-specific middleware not applied correctly: got %v", middlewarePaths)
	}
}

// TestRouterParameterWithSpecialCharacters tests parameters containing hyphens and underscores.
func TestRouterParameterWithSpecialCharacters(t *testing.T) {
	router := NewRouter()
	var extractedParams map[string]string

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		extractedParams = r.Params
	})

	router.Get("/files/:file-id", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{
		Method:  "GET",
		Path:    "/files/my-file-123",
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}

	router.ServeHTTP(rw, req)

	if extractedParams == nil {
		t.Fatal("params not set")
	}

	if fileId, ok := extractedParams["file-id"]; !ok || fileId != "my-file-123" {
		t.Errorf("parameter extraction failed: got %v", extractedParams)
	}
}

// TestRouterChainingMethods tests that router methods return the router for chaining.
func TestRouterChainingMethods(t *testing.T) {
	router := NewRouter()

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {})
	middleware := func(next Handler) Handler {
		return next
	}

	// These should all chain successfully
	result := router.Get("/users", handler).
		Post("/users", handler).
		Use(middleware).
		UsePath("/posts", middleware)

	if result == nil {
		t.Error("method chaining failed")
	}
}
