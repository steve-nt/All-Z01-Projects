package services

import (
	"os"
	"strings"
	"testing"
)

func setupTestService() *AsciiArtWeb {
	// Ensure we're in the project root directory for tests
	originalDir, _ := os.Getwd()
	if strings.Contains(originalDir, "services") {
		os.Chdir("..")
	}

	service := NewAsciiArtWeb()
	if err := service.LoadBanners(); err != nil {
		panic("Failed to load banners for testing: " + err.Error())
	}
	return service
}

// Test LoadBanners function
func TestLoadBanners(t *testing.T) {
	service := setupTestService()

	if len(service.banners) == 0 {
		t.Error("LoadBanners should load at least one banner")
	}

	// Check if standard banner exists
	if _, exists := service.banners["standard"]; !exists {
		t.Error("Standard banner should be loaded")
	}
}

// Test Generate function
func TestGenerate(t *testing.T) {
	service := setupTestService()

	tests := []struct {
		name       string
		text       string
		banner     string
		shouldFail bool
	}{
		{"Valid generation", "Hello", "standard", false},
		{"Single character", "A", "standard", false},
		{"With space", "A B", "standard", false},
		{"With newline", "A\nB", "standard", false},
		{"Invalid banner", "Hello", "nonexistent", true},
		{"Empty text", "", "standard", false}, // Generate allows empty, validation handles this
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Generate(tt.text, tt.banner)

			if tt.shouldFail {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == "" && tt.text != "" {
					t.Error("Expected non-empty result for non-empty input")
				}
			}
		})
	}
}

// Test GetAvailableBanners function
func TestGetAvailableBanners(t *testing.T) {
	service := setupTestService()

	banners := service.GetAvailableBanners()

	if len(banners) == 0 {
		t.Error("GetAvailableBanners should return at least one banner")
	}

	// Check if standard banner is in the list
	found := false
	for _, banner := range banners {
		if banner == "standard" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Standard banner should be in available banners list")
	}
}

// Test ValidateInput function
func TestValidateInput(t *testing.T) {
	service := setupTestService()

	tests := []struct {
		name       string
		text       string
		banner     string
		shouldFail bool
	}{
		{"Valid input", "Hello", "standard", false},
		{"Non-ASCII character", "Hello\x80", "standard", true},
		{"Valid newline", "Hello\nWorld", "standard", false},
		{"Valid tab", "Hello\tWorld", "standard", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateInput(tt.text, tt.banner)

			if tt.shouldFail {
				if err == nil {
					t.Error("Expected validation error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error: %v", err)
				}
			}
		})
	}
}
