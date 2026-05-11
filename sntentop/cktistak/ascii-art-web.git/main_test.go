package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"ascii-art-web/handlers"
	"ascii-art-web/services"
)

func setupTestHandler() *handlers.AsciiHandler {
	service := services.NewAsciiArtWeb()
	if err := service.LoadBanners(); err != nil {
		panic("Failed to load banners for testing: " + err.Error())
	}
	return handlers.NewAsciiHandler(service)
}

// Test GET request to home page
func TestHandleHome_GET(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.HandleHome(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "ASCII Art") {
		t.Error("Expected page to contain 'ASCII Art' title")
	}

	if !strings.Contains(body, "Generate") {
		t.Error("Expected page to contain 'Generate' button")
	}
}

// Test POST request with valid data
func TestHandleAsciiArt_POST_Valid(t *testing.T) {
	handler := setupTestHandler()

	form := url.Values{}
	form.Add("text", "Hello")
	form.Add("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.HandleAsciiArt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Hello") {
		t.Error("Expected response to contain input text 'Hello'")
	}
}

// Test POST request with empty text (should return 400)
func TestHandleAsciiArt_POST_EmptyText(t *testing.T) {
	handler := setupTestHandler()

	form := url.Values{}
	form.Add("text", "")
	form.Add("banner", "standard")

	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.HandleAsciiArt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Test GET request to invalid path (should return 404)
func TestHandleHome_GET_InvalidPath(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	w := httptest.NewRecorder()

	handler.HandleHome(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

// Test invalid HTTP method on home (should return 405)
func TestHandleHome_InvalidMethod(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.HandleHome(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

// Test invalid HTTP method on ascii-art (should return 405)
func TestHandleAsciiArt_InvalidMethod(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	w := httptest.NewRecorder()

	handler.HandleAsciiArt(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
