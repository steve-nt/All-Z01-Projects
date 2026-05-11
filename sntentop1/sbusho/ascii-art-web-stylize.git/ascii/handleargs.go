package ascii

import (
	"errors"
	"fmt"
	"strings"
)

func HandleArgs(args []string) error {
	if len(args) < 1 {
		return errors.New(" Usage: go run . [STRING] [BANNER(optional)]")
	}

	// Map of banner types to their corresponding file paths
	bannerFiles := map[string]string{
		"standard":   "banners/standard.txt",
		"shadow":     "banners/shadow.txt",
		"thinkertoy": "banners/thinkertoy.txt",
		"other":      "banners/other.txt",
	}

	// Default banner type and file path
	bannerType := "standard"
	bannerFile := bannerFiles[bannerType]

	// Check if the last argument is a valid banner type and update accordingly
	if _, exists := bannerFiles[args[len(args)-1]]; exists {
		bannerType = args[len(args)-1]
		bannerFile = bannerFiles[bannerType]
		args = args[:len(args)-1] // Remove the banner type from args
	}

	input := strings.Join(args[1:], " ")
	input = strings.Replace(input, `\n`, "\n", -1)

	asciiArt, err := GenerateAsciiArt(input, bannerFile)
	if err != nil {
		return fmt.Errorf("error: failed generate ASCII art: %w", err)
	}

	// asciiArt += "\n"

	fmt.Println(asciiArt)

	return nil
}
