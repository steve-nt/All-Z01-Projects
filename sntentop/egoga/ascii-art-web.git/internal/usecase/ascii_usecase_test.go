package usecase

import (
	"errors"
	"strings"
	"testing"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/domain"
)

// Mock repository implementing the interface
type MockBannerRepository struct{}

func (m *MockBannerRepository) LoadBanner(name string) (*domain.Banner, error) {
	if name == "invalid-banner" {
		return nil, errors.New("banner not found")
	}

	return &domain.Banner{
		Name: name,
		Lines: map[rune][]string{
			'A': {"  A  ", " A A ", "AAAAA", "A   A", "A   A", "A   A", "A   A", "A   A"},
			'B': {"BBBB ", "B   B", "BBBB ", "B   B", "B   B", "BBBB ", "B   B", "BBBB "},
		},
	}, nil
}

func TestConvertTextToAscii(t *testing.T) {
	mockRepo := &MockBannerRepository{} // Now works with the interface
	usecase := NewAsciiUsecase(mockRepo)

	tests := []struct {
		name        string
		request     *domain.ASCIITextRequest
		expected    string
		expectError bool
	}{
		{
			name: "Valid single letter",
			request: &domain.ASCIITextRequest{
				Text:   "A",
				Banner: "standard",
			},
			expected: "  A  \n A A \nAAAAA\nA   A\nA   A\nA   A\nA   A\nA   A\n",
		},
		{
			name: "Valid multiple letters",
			request: &domain.ASCIITextRequest{
				Text:   "AB",
				Banner: "standard",
			},
			expected: "  A  BBBB \n A A B   B\nAAAAABBBB \nA   AB   B\nA   AB   B\nA   ABBBB \nA   AB   B\nA   ABBBB \n",
		},
		{
			name: "Handles new lines",
			request: &domain.ASCIITextRequest{
				Text:   "A\\nB",
				Banner: "standard",
			},
			expected: "  A  \n A A \nAAAAA\nA   A\nA   A\nA   A\nA   A\nA   A\nBBBB \nB   B\nBBBB \nB   B\nB   B\nBBBB \nB   B\nBBBB \n",
		},
		{
			name: "Handles unsupported character",
			request: &domain.ASCIITextRequest{
				Text:   "C",
				Banner: "standard",
			},
			expectError: true,
		},
		{
			name: "Handles empty string",
			request: &domain.ASCIITextRequest{
				Text:   "",
				Banner: "standard",
			},
			expectError: true,
		},
		{
			name: "Handles invalid banner",
			request: &domain.ASCIITextRequest{
				Text:   "A",
				Banner: "invalid-banner",
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := usecase.ConvertTextToAscii(tc.request)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if strings.TrimSpace(result) != strings.TrimSpace(tc.expected) {
					t.Errorf("expected:\n%q\ngot:\n%q", tc.expected, result)
				}
			}
		})
	}
}
