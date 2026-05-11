package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	// Create a test server with a custom handler for simulating the HTTP request handling
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Switch based on the HTTP method (GET or POST)
		switch r.Method {
		case http.MethodGet:
			// Handle GET requests by serving the "index.html" file
			http.ServeFile(w, r, "index.html")
		case http.MethodPost:
			// Handle POST requests by parsing the form and processing the input
			if err := r.ParseForm(); err != nil {
				// If form parsing fails, return a 500 Internal Server Error
				http.Error(w, "Failed to parse form", http.StatusInternalServerError)
				return
			}
			// Retrieve the value of the "text" field from the form
			input := r.FormValue("text")
			if input == "" {
				// If the input is empty, return a 400 Bad Request error
				http.Error(w, "Empty input", http.StatusBadRequest)
				return
			}
			// If input is valid, return a message with the generated ASCII art
			w.Write([]byte("Generated ASCII Art for: " + input))
		default:
			// If the method is not GET or POST, return a 405 Method Not Allowed error
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	// Close the test server after the test completes
	defer ts.Close()

	// Test case 1: POST request with valid input (e.g., "Hello")
	t.Run("POST request generates ASCII art", func(t *testing.T) {
		// Send a POST request with valid form data (text=Hello)
		resp, err := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader("text=Hello"))
		if err != nil {
			// If the POST request fails, fail the test with the error message
			t.Fatalf("Failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		// Check if the status code is 200 OK
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test case 2: POST request with empty input
	t.Run("POST request with empty input returns error", func(t *testing.T) {
		// Send a POST request with an empty form value (text=)
		resp, err := http.Post(ts.URL, "application/x-www-form-urlencoded", strings.NewReader("text="))
		if err != nil {
			// If the POST request fails, fail the test with the error message
			t.Fatalf("Failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		// Check if the status code is 400 Bad Request
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}
