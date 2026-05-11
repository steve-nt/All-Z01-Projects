package main

import (
	"testing"
)

// TestReadTxt tests the ReadTxt function.
func TestReadTxt(t *testing.T) {
	// Test with a nonexistent file
	result := ReadTxt("nonexistent")
	if result != nil {
		t.Errorf("ReadTxt(\"nonexistent\") = %v; want nil", result)
	}
}

// TestReturn2dASCIIArray tests the Return2dASCIIArray function.
func TestReturn2dASCIIArray(t *testing.T) {
	input := []string{
		"@@@@@", "@   @", "@   @", "@   @", "@@@@@", "@   @", "@   @", "@   @", "",
	}
	result := Return2dASCIIArray(input)
	expected := [][]string{
		{"@@@@@", "@   @", "@   @", "@   @", "@@@@@", "@   @", "@   @", "@   @"},
	}

	if !equal2DSlices(result, expected) {
		t.Errorf("Return2dASCIIArray(%v) = %v; want %v", input, result, expected)
	}
}

// Helper function to compare two 2D slices for equality.
func equal2DSlices(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
