package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"ascii-art-web/services"
)

func setupTestHandler() *AsciiHandler {
	// Ensure we're in the project root directory for tests
	originalDir, _ := os.Getwd()
	if strings.Contains(originalDir, "handlers") {
		os.Chdir("..")
	}

	service := services.NewAsciiArtWeb()
	if err := service.LoadBanners(); err != nil {
		panic("Failed to load banners for testing: " + err.Error())
	}
	return NewAsciiHandler(service)
}

// Test NewAsciiHandler constructor
func TestNewAsciiHandler(t *testing.T) {
	service := services.NewAsciiArtWeb()
	handler := NewAsciiHandler(service)

	if handler == nil {
		t.Error("NewAsciiHandler should return a non-nil handler")
	}

	if handler.service != service {
		t.Error("Handler should store the provided service")
	}
}

// Test HandleHome function
func TestHandleHome(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"Valid GET request", http.MethodGet, "/", http.StatusOK},
		{"Invalid method", http.MethodPost, "/", http.StatusMethodNotAllowed},
		{"Invalid path", http.MethodGet, "/invalid", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.HandleHome(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test HandleAsciiArt function
func TestHandleAsciiArt(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		method         string
		text           string
		banner         string
		expectedStatus int
	}{
		{"Valid POST", http.MethodPost, "Hello", "standard", http.StatusOK},
		{"Empty text", http.MethodPost, "", "standard", http.StatusBadRequest},
		{"Invalid method", http.MethodGet, "Hello", "standard", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("text", tt.text)
			form.Add("banner", tt.banner)

			req := httptest.NewRequest(tt.method, "/ascii-art", strings.NewReader(form.Encode()))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			w := httptest.NewRecorder()

			handler.HandleAsciiArt(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Test HandleErrors function
func TestHandleErrors(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"Bad Request", http.StatusBadRequest, "Bad Request"},
		{"Not Found", http.StatusNotFound, "Not Found"},
		{"Internal Server Error", http.StatusInternalServerError, "Internal Server Error"},
		{"Method Not Allowed", http.StatusMethodNotAllowed, "Method Not Allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			handler.HandleErrors(w, tt.statusCode, tt.message)

			if w.Code != tt.statusCode {
				t.Errorf("Expected status %d, got %d", tt.statusCode, w.Code)
			}

			body := w.Body.String()
			if !strings.Contains(body, "Error") {
				t.Error("Error page should contain 'Error' text")
			}
		})
	}
}
