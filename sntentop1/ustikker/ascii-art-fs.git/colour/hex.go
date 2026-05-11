package asciiart

import (
	"fmt"
	"strconv"
)

// Function to convert hex color to ANSI
func HexToAnsi(hex string) (string, error) {
	if len(hex) != 7 || hex[0] != '#' {
		return "", fmt.Errorf("invalid hex format")
	}

	r, err := strconv.ParseInt(hex[1:3], 16, 0)
	if err != nil {
		return "", err
	}
	g, err := strconv.ParseInt(hex[3:5], 16, 0)
	if err != nil {
		return "", err
	}
	b, err := strconv.ParseInt(hex[5:7], 16, 0)
	if err != nil {
		return "", err
	}

	return RgbToAnsi(int(r), int(g), int(b)), nil
}
