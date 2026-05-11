package tests

import (
	"go-reloaded/pipeline"
	"reflect"
	"testing"
)

func TestApplyTransfomations_Basic(t *testing.T) {
	// Basic end-to-end scenario: articles, quoted-case change and punctuation spacing.
	input := []string{"a", "apple", "said", "\"hello\"", "world", "."}
	expected := []string{"an", "apple", "said", "\"Hello\"", "world", "."}
	result := pipeline.ApplyTransformations(input)
	// Assert the entire token slice matches expected.
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
func TestApplyTransfomations_QuotesAndCase(t *testing.T) {
	// Quoted short word should be uppercased inside quotes.
	input := []string{"she", "said", "\"hi\"", "there", "."}
	expected := []string{"she", "said", "\"HI\"", "there", "."}
	result := pipeline.ApplyTransformations(input)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
func TestApplyTransfomations_PunctuationSpacing(t *testing.T) {
	// Punctuation attachments should be combined with previous word.
	input := []string{"Hello", ",", "world", "!"}
	expected := []string{"Hello", "world!"}
	result := pipeline.ApplyTransformations(input)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
func TestApplyTransformations_EmptyInput(t *testing.T) {
	// Empty input should return empty output.
	input := []string{}
	expected := []string{}
	result := pipeline.ApplyTransformations(input)
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}
