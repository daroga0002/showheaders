package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShowHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add some test headers
	req.Header.Set("X-Test-Header", "test-value")
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Accept", "text/html")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/html; charset=utf-8")
	}

	// Check that response contains our headers
	body := rr.Body.String()

	if !strings.Contains(body, "X-Test-Header") {
		t.Error("response does not contain X-Test-Header")
	}

	if !strings.Contains(body, "test-value") {
		t.Error("response does not contain test-value")
	}

	if !strings.Contains(body, "User-Agent") {
		t.Error("response does not contain User-Agent")
	}

	if !strings.Contains(body, "test-agent") {
		t.Error("response does not contain test-agent")
	}

	// Check that HTML structure is present
	if !strings.Contains(body, "<html>") {
		t.Error("response does not contain HTML opening tag")
	}

	if !strings.Contains(body, "</html>") {
		t.Error("response does not contain HTML closing tag")
	}

	if !strings.Contains(body, "<table>") {
		t.Error("response does not contain table element")
	}
}

func TestShowHeaders_DisplaysMethod(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/test-path", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ShowHeaders)

			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			if !strings.Contains(body, method) {
				t.Errorf("response does not contain request method %s", method)
			}
		})
	}
}

func TestShowHeaders_DisplaysURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/test-path?query=value", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "/test-path") {
		t.Error("response does not contain request URL path")
	}
}

func TestShowHeaders_MultipleHeaderValues(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add header with multiple values
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "deflate")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()

	// Check that both values are present
	if !strings.Contains(body, "gzip") {
		t.Error("response does not contain first Accept-Encoding value")
	}

	if !strings.Contains(body, "deflate") {
		t.Error("response does not contain second Accept-Encoding value")
	}
}

func TestShowHeaders_JSONFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/?format=json", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Test-Header", "test-value")
	req.Header.Set("User-Agent", "test-agent")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse JSON response
	var response HeadersResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	// Check method
	if response.Method != "GET" {
		t.Errorf("wrong method in response: got %v want %v", response.Method, "GET")
	}

	// Check headers are present
	if response.Headers == nil {
		t.Fatal("headers map is nil")
	}

	if values, ok := response.Headers["X-Test-Header"]; !ok {
		t.Error("X-Test-Header not found in response")
	} else if len(values) == 0 || values[0] != "test-value" {
		t.Errorf("wrong X-Test-Header value: got %v want %v", values, []string{"test-value"})
	}
}

func TestShowHeaders_PlainFormat(t *testing.T) {
	req, err := http.NewRequest("GET", "/?format=plain", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Test-Header", "test-value")
	req.Header.Set("User-Agent", "test-agent")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/plain; charset=utf-8")
	}

	body := rr.Body.String()

	// Check that response contains expected content
	if !strings.Contains(body, "Method: GET") {
		t.Error("response does not contain method")
	}

	if !strings.Contains(body, "X-Test-Header: test-value") {
		t.Error("response does not contain X-Test-Header")
	}

	if !strings.Contains(body, "User-Agent: test-agent") {
		t.Error("response does not contain User-Agent")
	}

	if !strings.Contains(body, "--- Headers ---") {
		t.Error("response does not contain headers separator")
	}
}

func TestShowHeaders_DefaultHTMLFormat(t *testing.T) {
	// Test that no format parameter returns HTML
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "text/html; charset=utf-8")
	}

	body := rr.Body.String()
	if !strings.Contains(body, "<html>") {
		t.Error("response does not contain HTML tag")
	}
}

func TestShowHeaders_UnknownFormat(t *testing.T) {
	// Test that unknown format parameter returns HTML (default)
	req, err := http.NewRequest("GET", "/?format=unknown", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowHeaders)

	handler.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("handler returned wrong content type for unknown format: got %v want %v", contentType, "text/html; charset=utf-8")
	}
}
