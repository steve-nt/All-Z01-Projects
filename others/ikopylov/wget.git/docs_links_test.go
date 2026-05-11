package main

import (
	"os"
	"strings"
	"testing"
)

func TestREADME_RunbookRequiredLinks(t *testing.T) {
	b, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	s := string(b)

	requiredLinks := []string{
		"[`docs/requirements.md`](docs/requirements.md)",
		"[`docs/audit.md`](docs/audit.md)",
	}
	for _, link := range requiredLinks {
		if !strings.Contains(s, link) {
			t.Fatalf("README.md missing required link %q", link)
		}
	}

	requiredPaths := []string{
		"docs/requirements.md",
		"docs/audit.md",
	}
	for _, p := range requiredPaths {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected %q to exist: %v", p, err)
		}
	}
}

