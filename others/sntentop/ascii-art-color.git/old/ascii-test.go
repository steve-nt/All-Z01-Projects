package main

import (
	"testing"
)

// TestReturnAsciiCodeInt tests the ReturnAsciiCodeInt function.
func TestReturnAsciiCodeInt(t *testing.T) {
	result := ReturnAsciiCodeInt("A")
	expected := []int{33} // 'A' corresponds to 33 in ASCII offset by -32

	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("ReturnAsciiCodeInt(\"A\") = %v; want %v", result, expected)
	}
}

// TestReturnStringToEndlineArray tests the ReturnStringToEndlineArray function.
func TestReturnStringToEndlineArray(t *testing.T) {
	input := "hello\\nworld"
	result := ReturnStringToEndlineArray(input)
	expected := []string{"hello", "\\n", "world"}

	if !equalSlices(result, expected) {
		t.Errorf("ReturnStringToEndlineArray(%q) = %v; want %v", input, result, expected)
	}
}

// Helper function to compare two slices for equality.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestPrintMultipleCharacter ensures that PrintMultipleCharacter doesn't crash.
func TestPrintMultipleCharacter(t *testing.T) {
	asciiTemplates := [][]string{
		{"@@@@@", "@   @", "@   @", "@   @", "@@@@@", "@   @", "@   @", "@   @"},
	}
	PrintMultipleCharacter("A", asciiTemplates) // Just testing for crashes
}
