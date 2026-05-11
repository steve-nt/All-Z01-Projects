package handlers

import (
	"strings"
	"time"
)

// FormatDisplayTime converts timestamps (including RFC3339 like 2026-01-20T17:22:55Z)
// into a readable local time string.
// - "T" is just the date/time separator in RFC3339.
// - "Z" means UTC ("Zulu" time).
// We convert UTC times to Europe/Athens for display.
func FormatDisplayTime(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// Common layouts we may encounter:
	// - RFC3339 / RFC3339Nano (often with T and Z)
	// - SQLite CURRENT_TIMESTAMP default format
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	var (
		t   time.Time
		err error
	)

	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			break
		}
	}

	// If parse failed, return original
	if err != nil {
		return s
	}

	// Convert to Europe/Athens for display
	if loc, locErr := time.LoadLocation("Europe/Athens"); locErr == nil {
		t = t.In(loc)
	}

	// Nice readable format
	return t.Format("2006-01-02 15:04")
}
