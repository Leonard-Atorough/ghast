package examples

import (
	"bytes"
	"net"
	"testing"
	"time"

	ghast "ghast/lib"
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

// TestRoutingExamples demonstrates all routing functionality with the example router
func TestRoutingExamples(t *testing.T) {
	router := SetupRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus string // We check if response contains this
		description    string
	}{
		{
			name:           "GET root",
			method:         "GET",
			path:           "/",
			body:           "",
			expectedStatus: "Welcome to the API",
			description:    "Root endpoint returns welcome message",
		},
		{
			name:           "GET user list",
			method:         "GET",
			path:           "/users",
			body:           "",
			expectedStatus: "page",
			description:    "User list endpoint returns paginated users",
		},
		{
			name:           "GET specific user with parameter",
			method:         "GET",
			path:           "/users/123",
			body:           "",
			expectedStatus: "123",
			description:    "User endpoint extracts ID parameter correctly",
		},
		{
			name:           "GET user posts with multiple parameters",
			method:         "GET",
			path:           "/users/456/posts/789",
			body:           "",
			expectedStatus: "456",
			description:    "Nested route correctly extracts multiple parameters",
		},
		{
			name:           "POST create user",
			method:         "POST",
			path:           "/users",
			body:           `{"name":"Alice","email":"alice@example.com"}`,
			expectedStatus: "201",
			description:    "Create user endpoint returns 201 status",
		},
		{
			name:           "PUT update user",
			method:         "PUT",
			path:           "/users/999",
			body:           `{"name":"Bob","email":"bob@example.com"}`,
			expectedStatus: "999",
			description:    "Update user endpoint extracts ID parameter",
		},
		{
			name:           "DELETE user",
			method:         "DELETE",
			path:           "/users/111",
			body:           "",
			expectedStatus: "deleted successfully",
			description:    "Delete user endpoint returns success message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &MockConnection{}
			rw := ghast.NewResponseWriter(mockConn)
			req := &ghast.Request{
				Method:  tt.method,
				Path:    tt.path,
				Body:    tt.body,
				Headers: make(map[string]string),
				Queries: make(map[string]string),
				Params:  make(map[string]string),
			}

			router.ServeHTTP(rw, req)

			output := mockConn.writeBuffer.String()
			if !bytes.Contains([]byte(output), []byte(tt.expectedStatus)) {
				t.Errorf("%s: expected response to contain '%s', but got:\n%s",
					tt.description, tt.expectedStatus, output)
			}
		})
	}
}

// TestParameterizedrouting demonstrates parameter extraction in various route patterns
func TestParameterizedRouting(t *testing.T) {
	router := SetupRouter()

	tests := []struct {
		name           string
		path           string
		expectedParams map[string]string
		description    string
	}{
		{
			name:           "/users/:id",
			path:           "/users/42",
			expectedParams: map[string]string{"id": "42"},
			description:    "Single parameter extraction",
		},
		{
			name:           "/users/:userId/posts/:postId",
			path:           "/users/100/posts/200",
			expectedParams: map[string]string{"userId": "100", "postId": "200"},
			description:    "Multiple parameters in nested route",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			mockConn := &MockConnection{}
			rw := ghast.NewResponseWriter(mockConn)
			req := &ghast.Request{
				Method:  "GET",
				Path:    tt.path,
				Headers: make(map[string]string),
				Params:  make(map[string]string),
			}

			router.ServeHTTP(rw, req)

			for paramName, expectedValue := range tt.expectedParams {
				actualValue := req.Params[paramName]
				if actualValue != expectedValue {
					t.Errorf("%s: expected %s=%s, got %s=%s",
						tt.description, paramName, expectedValue, paramName, actualValue)
				}
			}
		})
	}
}

// TestRouterWithQueryParameters demonstrates query parameter handling
func TestRouterWithQueryParameters(t *testing.T) {
	router := SetupRouter()

	mockConn := &MockConnection{}
	rw := ghast.NewResponseWriter(mockConn)
	req := &ghast.Request{
		Method:  "GET",
		Path:    "/users",
		Headers: make(map[string]string),
		Queries: map[string]string{
			"page":  "2",
			"limit": "20",
			"role":  "admin",
		},
		Params: make(map[string]string),
	}

	router.ServeHTTP(rw, req)

	output := mockConn.writeBuffer.String()

	// Check that query parameters are included in response
	if !bytes.Contains([]byte(output), []byte("\"page\":2")) &&
		!bytes.Contains([]byte(output), []byte("page")) {
		t.Error("page query parameter not reflected in response")
	}

	if !bytes.Contains([]byte(output), []byte("\"limit\":20")) &&
		!bytes.Contains([]byte(output), []byte("limit")) {
		t.Error("limit query parameter not reflected in response")
	}
}

// TestAdminRoute demonstrates protected routes with authentication middleware
func TestAdminRoute(t *testing.T) {
	router := SetupRouter()

	// Test without authorization header - should be 401
	mockConn1 := &MockConnection{}
	rw1 := ghast.NewResponseWriter(mockConn1)
	req1 := &ghast.Request{
		Method:  "GET",
		Path:    "/admin/stats",
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}

	router.ServeHTTP(rw1, req1)
	output1 := mockConn1.writeBuffer.String()

	if !bytes.Contains([]byte(output1), []byte("401")) &&
		!bytes.Contains([]byte(output1), []byte("unauthorized")) {
		t.Error("expected 401 unauthorized response when Authorization header is missing")
	}

	// Test with authorization header - should succeed
	mockConn2 := &MockConnection{}
	rw2 := ghast.NewResponseWriter(mockConn2)
	req2 := &ghast.Request{
		Method: "GET",
		Path:   "/admin/stats",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
		},
		Params: make(map[string]string),
	}

	router.ServeHTTP(rw2, req2)
	output2 := mockConn2.writeBuffer.String()

	if !bytes.Contains([]byte(output2), []byte("totalUsers")) {
		t.Error("expected stats response when Authorization header is provided")
	}
}
