package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestLoggingMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t)

	excludePaths := map[string]bool{
		"/health": true,
	}

	middleware := LoggingMiddleware(logger, excludePaths)

	// Create a simple handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware(nextHandler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that request was processed
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoggingMiddleware_ExcludedPath(t *testing.T) {
	logger := zaptest.NewLogger(t)

	excludePaths := map[string]bool{
		"/health": true,
	}

	middleware := LoggingMiddleware(logger, excludePaths)

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Verify handler was still called
	if !handlerCalled {
		t.Error("next handler was not called for excluded path")
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestLoggingMiddleware_CapturesStatusCode(t *testing.T) {
	logger := zaptest.NewLogger(t)

	excludePaths := map[string]bool{}

	middleware := LoggingMiddleware(logger, excludePaths)

	testCases := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
		{"BadRequest", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})

			handler := middleware(nextHandler)

			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.statusCode)
			}
		})
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapped := newResponseWriter(rr)

	wrapped.WriteHeader(http.StatusNotFound)

	if wrapped.statusCode != http.StatusNotFound {
		t.Errorf("wrapped response writer has wrong status code: got %v want %v", wrapped.statusCode, http.StatusNotFound)
	}

	if rr.Code != http.StatusNotFound {
		t.Errorf("underlying response writer has wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}
}

func TestResponseWriter_DefaultStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	wrapped := newResponseWriter(rr)

	// Default status should be 200
	if wrapped.statusCode != http.StatusOK {
		t.Errorf("wrapped response writer has wrong default status code: got %v want %v", wrapped.statusCode, http.StatusOK)
	}
}

func TestLoggingMiddleware_WithUserAgent(t *testing.T) {
	logger := zaptest.NewLogger(t)

	excludePaths := map[string]bool{}

	middleware := LoggingMiddleware(logger, excludePaths)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "test-agent/1.0")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func BenchmarkLoggingMiddleware(b *testing.B) {
	logger := zap.NewNop()

	excludePaths := map[string]bool{}

	middleware := LoggingMiddleware(logger, excludePaths)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
