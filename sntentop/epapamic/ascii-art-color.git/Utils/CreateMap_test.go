package utils

import (
	"os"
	"testing"
)

func TestCreateMap(t *testing.T) {
	testCases := []struct {
		name     string
		mockMap  asciiMap
		expected map[rune][]string
	}{
		{
			name:    "Valid File",
			mockMap: asciiMap{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ref, err := os.Open(mapPath)
			if err != nil {
				t.Errorf("Could not find standard.txt at the designated path %v", mapPath)
			}
			testCase.mockMap.ref = ref
			testCase.mockMap.CreateMap()
			for i := 32; i <= 126; i++ {
				char := rune(i)
				if _, exists := testCase.mockMap.content[char]; !exists {
					t.Errorf("character %v not found in map", char)
				}
				if len(testCase.mockMap.content[char]) != 8 {
					t.Errorf("character %v has %v lines, need 8", char, len(testCase.mockMap.content[char]))
				}
			}
		})
	}
}
