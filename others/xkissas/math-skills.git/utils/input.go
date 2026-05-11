package utils

import (
	"fmt"
	"os"
)

func CheckFile() string {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	return filename
}

func OpenFile(inputFile string) (*os.File, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return file, nil
}
