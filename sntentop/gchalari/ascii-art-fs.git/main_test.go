package main

import "testing"

func TestParseArgsSingleString(t *testing.T) {
	input, banner, err := parseArgs([]string{"hello"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if input != "hello" {
		t.Fatalf("expected input hello, got %q", input)
	}

	if banner != "standard" {
		t.Fatalf("expected default banner standard, got %q", banner)
	}
}

func TestParseArgsStringAndBanner(t *testing.T) {
	input, banner, err := parseArgs([]string{"hello", "shadow"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if input != "hello" {
		t.Fatalf("expected input hello, got %q", input)
	}

	if banner != "shadow" {
		t.Fatalf("expected banner shadow, got %q", banner)
	}
}

func TestParseArgsTooManyArgs(t *testing.T) {
	_, _, err := parseArgs([]string{"banana", "standard", "abc"})
	if err == nil {
		t.Fatal("expected usage error for too many arguments")
	}

	if err.Error() != usageMessage {
		t.Fatalf("expected usage message, got %q", err.Error())
	}
}

func TestParseArgsRejectsOptionLikeArgument(t *testing.T) {
	_, _, err := parseArgs([]string{"--bad"})
	if err == nil {
		t.Fatal("expected usage error for option-like argument")
	}
}

func TestParseArgsRejectsUnsupportedBanner(t *testing.T) {
	_, _, err := parseArgs([]string{"hello", "invalidbanner"})
	if err == nil {
		t.Fatal("expected usage error for unsupported banner")
	}

	if err.Error() != usageMessage {
		t.Fatalf("expected usage message, got %q", err.Error())
	}
}
