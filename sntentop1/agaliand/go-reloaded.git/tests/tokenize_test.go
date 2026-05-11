package tests

import (
	"reflect"
	"testing"

	"go-reloaded/pipeline"
)

// TestTokenize_Basic verifies the tokenizer splits words and special symbols.
// Input: runes containing words, angle brackets and a double-quoted word.
// Expectation: tokens should separate words and emit angle brackets and quote characters as tokens.
func TestTokenize_Basic(t *testing.T) {
	input := []rune(`Hello <world> "Go"`)
	expected := []string{"Hello", "<", "world", ">", "\"", "Go", "\""}

	// Call the tokenizer and compare the returned slice with expected tokens.
	result := pipeline.Tokenize(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestTokenize_WithSpaces ensures consecutive spaces are treated as a single separator.
func TestTokenize_WithSpaces(t *testing.T) {
	input := []rune("Hello   world")
	expected := []string{"Hello", "world"}

	result := pipeline.Tokenize(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestTokenize_OnlySpecials checks tokenizer behavior with only special characters.
func TestTokenize_OnlySpecials(t *testing.T) {
	input := []rune(`<>""`)
	expected := []string{"<", ">", "\"", "\""}

	result := pipeline.Tokenize(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestTokenize_EmptyInput verifies that empty input yields no tokens.
func TestTokenize_EmptyInput(t *testing.T) {
	input := []rune("")
	expected := []string{}

	result := pipeline.Tokenize(input)

	if len(result) != len(expected) {
		t.Errorf("Expected %v tokens, got %v", len(expected), len(result))
	}

}
