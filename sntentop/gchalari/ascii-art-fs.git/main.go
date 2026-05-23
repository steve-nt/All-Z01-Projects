package main

import (
	"ascii-art/ascii"
	"errors"
	"fmt"
	"os"
	"strings"
)

const usageMessage = "Usage: go run . [STRING] [BANNER]\n\nEX: go run . something standard"

var supportedBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
	"lineart":    true,
}

func usageError() error {
	return errors.New(usageMessage)
}

func parseArgs(args []string) (input string, banner string, err error) {
	banner = "standard"

	switch len(args) {
	case 1:
		if strings.HasPrefix(args[0], "--") {
			return "", "", usageError()
		}
		input = args[0]

	case 2:
		if strings.HasPrefix(args[0], "--") || strings.HasPrefix(args[1], "--") {
			return "", "", usageError()
		}

		if args[1] == "" {
			return "", "", usageError()
		}

		input = args[0]
		banner = args[1]

	default:
		return "", "", usageError()
	}

	if !supportedBanners[banner] {
		return "", "", usageError()
	}

	return input, banner, nil
}

func main() {
	input, bannerName, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	if input == "" {
		return
	}

	banner, err := ascii.LoadBanner(bannerName)
	if err != nil {
		fmt.Println(err)
		return
	}

	result := ascii.Render(input, banner)
	fmt.Print(result)
}
