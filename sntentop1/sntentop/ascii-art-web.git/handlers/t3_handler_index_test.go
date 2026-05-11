package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock template definitions for testing
var mockTpl = template.Must(template.New("test").Parse(`
{{define "index.html"}}<h1>Welcome to the Index Page</h1>{{end}}
{{define "404.html"}}<h1>404 Not Found</h1>{{end}}
{{define "500.html"}}<h1>500 Internal Server Error</h1>{{end}}
{{define "400.html"}}<h1>400 Bad Request: {{.Message}}</h1>{{end}}
`))

func normalizeWhitespace(s string) string {
	return strings.TrimSpace(s)
}

func TestIndexHandler_IndexPage(t *testing.T) {
	tpl = mockTpl                               // Override the global tpl variable with mockTpl
	req := httptest.NewRequest("GET", "/", nil) // Create a test request recorder
	rec := httptest.NewRecorder()               // Create a test response recorder
	Index(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
	expectedBody := "<h1>Welcome to the Index Page</h1>"
	actualBody := normalizeWhitespace(rec.Body.String())
	if actualBody != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, actualBody)
	}
}

func TestIndexHandler_NotFoundPage(t *testing.T) {
	tpl = mockTpl                                           // Override the global tpl variable with mockTpl
	req := httptest.NewRequest("GET", "/invalid-path", nil) // Create a test request recorder
	rec := httptest.NewRecorder()                           // Create a test response recorder
	Index(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
	}
	expectedBody := "<h1>404 Not Found</h1>"
	actualBody := normalizeWhitespace(rec.Body.String())
	if actualBody != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, actualBody)
	}
}

func TestHandleServerError(t *testing.T) {

	tpl = mockTpl                         // Override the global tpl variable with mockTpl
	rec := httptest.NewRecorder()         // Create a test response recorder
	testMessage := "Something went wrong" // Call handleServerError with a test message
	handleServerError(rec, testMessage)
	// Check the status code
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	// Check the response body
	expectedBody := "<h1>500 Internal Server Error</h1>"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestHandleBadRequest(t *testing.T) {
	tpl = mockTpl                  // Override the global tpl variable with mockTpl
	rec := httptest.NewRecorder()  // Create a test response recorder
	testMessage := "Invalid input" // Call handleBadRequest with a test message
	handleBadRequest(rec, testMessage)
	// Check the status code
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	}
	// Check the response body
	expectedBody := "<h1>400 Bad Request: Invalid input</h1>"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}
