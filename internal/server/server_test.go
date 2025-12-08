package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"

	"showheaders/internal/config"
	"showheaders/internal/handlers"
)

func TestNew(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cfg := &config.Config{
		Port:     8080,
		Hostname: "",
	}

	srv := New(cfg, logger)

	if srv == nil {
		t.Fatal("New() returned nil")
	}

	if srv.httpServer == nil {
		t.Error("httpServer is nil")
	}

	if srv.logger == nil {
		t.Error("logger is nil")
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cfg := &config.Config{
		Port:     0, // Use port 0 to let OS assign a free port
		Hostname: "localhost",
	}

	srv := New(cfg, logger)

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Check for startup errors
	select {
	case err := <-serverErr:
		t.Errorf("Server start error: %v", err)
	default:
		// No error, as expected
	}
}

func TestServer_HealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Health)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response handlers.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("wrong health status: got %v want %v", response.Status, "healthy")
	}
}

func TestServer_HeadersEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Test", "value")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ShowHeaders)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("wrong content type: got %v want %v", contentType, "text/html; charset=utf-8")
	}
}
