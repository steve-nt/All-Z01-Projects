package repository

import (
	"testing"
)

func TestNewBannerRepository(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		validBanners []string
		expected     *BannerRepository
	}{
		{
			name:         "Test1",
			path:         "existing/path",
			validBanners: []string{"banner1", "banner2", "banner3"},
			expected: &BannerRepository{
				basePath: "existing/path",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := NewBannerRepository(testCase.path)
			if result.basePath != testCase.expected.basePath {
				t.Errorf("Expected basePath: %s, got %s", testCase.expected.basePath, result.basePath)
			}
		})
	}
}

var mockBannerRepository = &BannerRepository{
	basePath: "../assets/banners/",
}

func TestLoadBanner(t *testing.T) {
	testCases := []struct {
		name        string
		banner      string
		expectedErr bool
	}{
		{
			name:        "Standard banner",
			banner:      "standard",
			expectedErr: false,
		},
		{
			name:        "Shadow banner",
			banner:      "shadow",
			expectedErr: false,
		},
		{
			name:        "Thinkertoy banner",
			banner:      "thinkertoy",
			expectedErr: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := mockBannerRepository.LoadBanner(testCase.banner)
			if testCase.expectedErr && err == nil {
				t.Error("Expected error, got none")
			}
			if !testCase.expectedErr && err != nil {
				t.Errorf("Did not expect error and got: %v", err)
			}
			for i := spaceASCII; i <= tidleASCII; i++ {
				char := rune(i)
				if _, exists := result.Lines[char]; !exists {
					t.Errorf("character %v not found in map", char)
				}
				if len(result.Lines[char]) != charHeight {
					t.Errorf("character %v has %v lines, need 8", char, len(result.Lines[char]))
				}
			}
		})
	}
}
