package main

import (
	"net/http"
	"testing"
	"time"
)

// Change this if your server runs on a different port
var baseURL = "http://localhost:8080"

// ✅ Wait for the server to start before testing
func waitForServer() {
	for i := 0; i < 10; i++ {
		_, err := http.Get(baseURL)
		if err == nil {
			return // Server is up
		}
		time.Sleep(500 * time.Millisecond) // Wait for 0.5s before retrying
	}
}

// ✅ Test Home Page
func TestHomePage(t *testing.T) {
	waitForServer()
	resp, err := http.Get(baseURL + "/")
	if err != nil {
		t.Fatalf("❌ Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %v", resp.StatusCode)
	}
}

// ✅ Test Artist Page
func TestArtistPage(t *testing.T) {
	resp, err := http.Get(baseURL + "/Artist/1") // Test for artist ID 1
	if err != nil {
		t.Fatalf("❌ Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %v", resp.StatusCode)
	}
}

// ✅ Test 404 Page
func Test404Page(t *testing.T) {
	resp, err := http.Get(baseURL + "/random-page") // Invalid page
	if err != nil {
		t.Fatalf("❌ Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("❌ Expected status 404, got %v", resp.StatusCode)
	}
}

// ✅ Test CSS File
func TestCssHandler(t *testing.T) {
	resp, err := http.Get(baseURL + "/frontend/css/index.css")
	if err != nil {
		t.Fatalf("❌ Failed to fetch CSS file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %v", resp.StatusCode)
	}
}

// ✅ Test Image File
func TestImageHandler(t *testing.T) {
	resp, err := http.Get(baseURL + "/frontend/images/sample.jpg")
	if err != nil {
		t.Fatalf("❌ Failed to fetch image file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Expected status 200, got %v", resp.StatusCode)
	}
}
