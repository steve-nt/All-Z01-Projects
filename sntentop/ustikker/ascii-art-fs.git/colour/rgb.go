package asciiart

import (
	"fmt"
	"regexp"
	"strconv"
)

// Function to convert RGB to ANSI escape code
func RgbToAnsi(r, g, b int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// Function to convert RGB string to ANSI
func RgbToAnsiFromString(rgb string) (string, error) {
	// Update regex pattern to match the new format: rgb255,0,0
	re := regexp.MustCompile(`rgb(\d+),\s*(\d+),\s*(\d+)`)
	matches := re.FindStringSubmatch(rgb)

	if len(matches) != 4 {
		return "", fmt.Errorf("invalid RGB format")
	}

	// Parse the RGB values
	r, err1 := strconv.Atoi(matches[1])
	g, err2 := strconv.Atoi(matches[2])
	b, err3 := strconv.Atoi(matches[3])

	if err1 != nil || err2 != nil || err3 != nil {
		return "", fmt.Errorf("error parsing RGB values")
	}

	return RgbToAnsi(r, g, b), nil
}
