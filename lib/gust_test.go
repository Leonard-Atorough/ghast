package gust

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// MockConnection implements net.Conn for testing
type MockConnection struct {
	writeBuffer bytes.Buffer
}

func (m *MockConnection) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *MockConnection) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *MockConnection) Close() error {
	return nil
}

func (m *MockConnection) LocalAddr() net.Addr {
	return nil
}

func (m *MockConnection) RemoteAddr() net.Addr {
	return nil
}

func (m *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

// TestRouterGet tests that GET routes are registered and matched correctly
func TestRouterGet(t *testing.T) {
	router := NewRouter()
	handlerCalled := false

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		handlerCalled = true
	})

	router.Get("/test", handler)

	// Manually verify route was registered
	// (This requires accessing private router struct - in real tests would use interface methods)
	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/test", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if !handlerCalled {
		t.Error("GET handler was not called")
	}
}

// TestRouterPost tests that POST routes are registered
func TestRouterPost(t *testing.T) {
	router := NewRouter()
	handlerCalled := false

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		handlerCalled = true
	})

	router.Post("/create", handler)

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "POST", Path: "/create", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	if !handlerCalled {
		t.Error("POST handler was not called")
	}
}

// TestRouter404 tests that non-existent routes return 404
func TestRouter404(t *testing.T) {
	router := NewRouter()

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/nonexistent", Headers: make(map[string]string)}

	router.ServeHTTP(rw, req)

	// Verify status code was set to 404
	// We can check by looking at what was written
	output := mockConn.writeBuffer.String()
	if !bytes.Contains([]byte(output), []byte("404")) {
		t.Error("404 response not found in output")
	}
}

// TestMiddlewareChaining tests that middleware is applied in correct order
func TestMiddlewareChaining(t *testing.T) {
	callOrder := []string{}

	middleware1 := func(next Handler) Handler {
		return HandlerFunc(func(w ResponseWriter, r *Request) {
			callOrder = append(callOrder, "m1-before")
			next.ServeHTTP(w, r)
			callOrder = append(callOrder, "m1-after")
		})
	}

	middleware2 := func(next Handler) Handler {
		return HandlerFunc(func(w ResponseWriter, r *Request) {
			callOrder = append(callOrder, "m2-before")
			next.ServeHTTP(w, r)
			callOrder = append(callOrder, "m2-after")
		})
	}

	handler := HandlerFunc(func(w ResponseWriter, r *Request) {
		callOrder = append(callOrder, "handler")
	})

	// Chain: m1 -> m2 -> handler
	chainedHandler := middleware1(middleware2(handler))

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{Method: "GET", Path: "/", Headers: make(map[string]string)}

	chainedHandler.ServeHTTP(rw, req)

	// Expected order: m1-before -> m2-before -> handler -> m2-after -> m1-after
	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}

	if len(callOrder) != len(expected) {
		t.Errorf("Call order length mismatch: got %d, want %d", len(callOrder), len(expected))
		return
	}

	for i, v := range callOrder {
		if v != expected[i] {
			t.Errorf("Call order mismatch at index %d: got %s, want %s", i, v, expected[i])
		}
	}
}

// TestRequestQueryParam tests Query() method
func TestRequestQueryParam(t *testing.T) {
	req := &Request{
		Queries: map[string]string{
			"name":  "alice",
			"email": "alice@example.com",
		},
	}

	if req.Query("name") != "alice" {
		t.Errorf("Query('name') failed")
	}

	if req.Query("nonexistent") != "" {
		t.Errorf("Query('nonexistent') should return empty string")
	}
}

// TestRequestJSON tests JSON unmarshaling
func TestRequestJSON(t *testing.T) {
	type testData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	req := &Request{
		Body: `{"name":"bob","age":30}`,
	}

	var data testData
	err := req.JSON(&data)

	if err != nil {
		t.Errorf("JSON unmarshaling failed: %v", err)
	}

	if data.Name != "bob" || data.Age != 30 {
		t.Errorf("JSON data mismatch: got %+v", data)
	}
}

// TestResponseJSON tests JSON marshaling and response
func TestResponseJSON(t *testing.T) {
	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)

	testData := map[string]string{"message": "test"}
	err := rw.JSON(200, testData)

	if err != nil {
		t.Errorf("JSON response failed: %v", err)
	}

	output := mockConn.writeBuffer.String()

	if !bytes.Contains([]byte(output), []byte("application/json")) {
		t.Error("Content-Type header not set to application/json")
	}

	if !bytes.Contains([]byte(output), []byte("\"message\":\"test\"")) {
		t.Error("JSON body not found in response")
	}
}

// TestResponseChainingStatus tests that Status() returns the ResponseWriter for chaining
func TestResponseStatus(t *testing.T) {
	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)

	rw.Status(201).SetHeader("X-Custom", "value")

	// Verify header was set (indicates chaining worked)
	headers := rw.Header()
	if headers["X-Custom"] != "value" {
		t.Error("Response chaining failed")
	}
}

// TestHandlerFunc tests that HandlerFunc implements Handler interface
func TestHandlerFunc(t *testing.T) {
	called := false
	hf := HandlerFunc(func(w ResponseWriter, r *Request) {
		called = true
	})

	mockConn := &MockConnection{}
	rw := NewResponseWriter(mockConn)
	req := &Request{}

	hf.ServeHTTP(rw, req)

	if !called {
		t.Error("HandlerFunc did not call the underlying function")
	}
}
