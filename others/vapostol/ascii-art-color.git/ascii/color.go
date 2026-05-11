package ascii

import "fmt"

// ANSI color codes map
var AnsiColors = map[string]string{
	"black":  "\033[30m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"white":  "\033[37m",
	"orange": "\033[38;5;214m",
	"reset":  "\033[0m",
}

// ColorToAnsi converts color names or codes to ANSI escape codes
func ColorToAnsi(color string) (string, error) {
	if code, exists := AnsiColors[color]; exists {
		return code, nil
	}
	return "", fmt.Errorf("invalid color format")
}
