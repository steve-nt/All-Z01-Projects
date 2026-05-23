package main

import (
	"os"
	"strings"
)

// LoadBanner reads a banner file and returns a map of rune to its 8-line graphical representation.
// The banner file must follow the standard format: 95 characters (ASCII 32-126),
// each 8 lines tall, separated by a blank line.
func LoadBanner(filename string) (map[rune][]string, error) {
	byteSlice, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Trim leading newline in the banners first element (space) and turn []byte -> string
	content := strings.TrimPrefix(string(byteSlice), "\n")

	// Split the string into []string , split on blank line
	characters := strings.Split(content, "\n\n")

	bannerMap := make(map[rune][]string)
	for i, block := range characters {
		lines := strings.Split(block, "\n")
		bannerMap[rune(32+i)] = lines
	}

	return bannerMap, nil
}
