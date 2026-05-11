package asciiart

import (
	"strings"
)

// Function to identify color format and convert to ANSI
func ColorToAnsi(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))

	if strings.HasPrefix(input, "#") {
		return HexToAnsi(input)
	} else if strings.HasPrefix(input, "rgb") {
		return RgbToAnsiFromString(input)
	} else {
		return NamedColorToAnsi(input)
	}
}
