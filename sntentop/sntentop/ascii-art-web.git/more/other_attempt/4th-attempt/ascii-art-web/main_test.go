package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if !strings.Contains(rr.Body.String(), "ASCII Art Generator") {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestAsciiArtHandler(t *testing.T) {
	// Prepare a POST request with form data
	formData := "text=TEST&banner=standard"
	req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(asciiArtHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if !strings.Contains(rr.Body.String(), "T") { // Example check
		t.Errorf("handler did not return expected ASCII art")
	}
}

func TestAsciiArtHandlerError(t *testing.T) {
	// Simulate Internal Server Error
	formData := "text=error500&banner=standard"
	req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader(formData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(asciiArtHandler)

	handler.ServeHTTP(rr, req)

	// Check for the 500 Internal Server Error
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler did not return Internal Server Error, got %v", status)
	}
}

func TestGenerateASCIIArt(t *testing.T) {
	// Define test cases
	tests := []struct {
		input    string
		banner   string
		expected string
	}{
		{"A", "banners/standard.txt", "  A  \n A A \nAAAAA\nA   A\n"},
		{"B", "banners/standard.txt", "BBBB \nB   B\nBBBB \nB   B\nBBBB \n"},
		{"", "banners/standard.txt", ""},              // Test empty input
		{"invalid-banner", "banners/invalid.txt", ""}, // Invalid banner
	}

	for _, test := range tests {
		result, err := generateASCIIArt(test.input, test.banner)

		// If banner is invalid, expect an error
		if test.banner == "banners/invalid.txt" && err == nil {
			t.Errorf("Expected error for invalid banner, got nil")
		}

		// Compare the result with expected output
		if result != test.expected {
			t.Errorf("For input '%s', expected:\n%s\nGot:\n%s", test.input, test.expected, result)
		}
	}
}
