package ghast

import (
	"log"
	"testing"
)

type testHandler struct{}

func (h *testHandler) handleRequest(w ResponseWriter, r *Request) {
	// No-op handler for testing
}

func CreateNewServerWithDefaultConfigTest(t *testing.T) {
	handler := &testHandler{}
	config := &serverConfig{
		Address:                 ":8080",
		HidePort:                false,
		GracefulShutdownTimeout: 30,
		OnShutdownError: func(err error) {
			log.Printf("Error during shutdown: %v", err)
		},
	}
	server := newServer(handler, config)

	if server == nil {
		t.Fatal("Expected NewServer to return a non-nil Server instance")
	}

	if server.config.Address != config.Address {
		t.Errorf("Expected server address to be %s, got %s", config.Address, server.config.Address)
	}

	if server.config.GracefulShutdownTimeout != config.GracefulShutdownTimeout {
		t.Errorf("Expected graceful shutdown timeout to be %d, got %d", config.GracefulShutdownTimeout, server.config.GracefulShutdownTimeout)
	}

	if server.requestHandler == nil {
		t.Error("Expected server to have a non-nil request handler")
	}
}
