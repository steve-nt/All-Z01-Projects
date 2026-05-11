package utils

import (
	"errors"
	"fmt"
)

func Ascii_art(input string, banner string) (string, error) {
	lines, err := ParseBanner(banner)
	if err != nil {
		fmt.Println("Error parsing banner:", err)
		return "", errors.New(err.Error())
	}

	// finding the corresponding index of every character's first line in the standard.txt
	indeces, err := Indexing(input)
	if err != nil {
		fmt.Println("Error indexing:", err)
		return "", errors.New(err.Error())
	}

	// spliting into a new subslice every time a newline character is encountered
	subSlices := NewLineHandling(indeces)
	// printing line by line into the command line
	result := ""
	for _, slice := range subSlices {
		// PrintAscii(slice, lines)
		result += OutputAscii(slice, lines)
	}
	return result, nil
}
