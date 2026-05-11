package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const bannerHeight = 8

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ascii-art.go [STRING] [BANNER]")
		os.Exit(1)
	}

	inputString := os.Args[1]
	bannerName := "standard"
	if len(os.Args) > 2 {
		bannerName = os.Args[2]
	}

	bannerFilePath := filepath.Join("banners", bannerName+".txt")
	bannerMap, err := loadBanner(bannerFilePath)
	if err != nil {
		fmt.Printf("Error: Banner '%s' not found.\n", bannerName)
		availableBanners, listErr := getAvailableBanners("banners")
		if listErr != nil {
			fmt.Println("Error retrieving available banners:", listErr)
			os.Exit(1)
		}
		fmt.Println("Available banners are:")
		for _, b := range availableBanners {
			fmt.Println("-", b)
		}
		os.Exit(1)
	}

	// Handle empty input string
	if inputString == "" || inputString == "\\n" {
		fmt.Println()
		return
	}

	// Split the input string on newline characters
	lines := strings.Split(inputString, "\\n")

	var output []string

	// Process each line
	for idx, line := range lines {
		if line == "" {
			// If the line is empty, add an empty line to the output
			output = append(output)
		} else {
			asciiArtLines, err := renderString(line, bannerMap)
			if err != nil {
				fmt.Println("Error rendering string:", err)
				os.Exit(1)
			}
			output = append(output, asciiArtLines...)
		}

		// Add an empty line between lines if it's not the last line
		if idx < len(lines)-1 {
			output = append(output)
		}
	}

	// Remove any leading empty lines
	output = trimLeadingEmptyLines(output)

	// Find the maximum line length
	maxLineLength := 0
	for _, line := range output {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}

	// Pad empty lines with spaces to match maxLineLength
	for i, line := range output {
		if line == "" {
			output[i] = strings.Repeat(" ", maxLineLength)
		}
	}

	// Print the output lines
	for _, line := range output {
		fmt.Println(line)
	}

	// Print one lines of spaces with length maxLineLength
	fmt.Println(strings.Repeat(" ", maxLineLength))
}

// loadBanner reads the banner file and returns a map from characters to their ASCII art.
func loadBanner(filePath string) (map[rune][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open banner file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bannerMap := make(map[rune][]string)
	var currentChar rune = 32 // ASCII space
	var charLines []string

	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if lineCount < bannerHeight {
			charLines = append(charLines, line)
			lineCount++
		}
		if lineCount == bannerHeight {
			bannerMap[currentChar] = charLines
			charLines = []string{}
			lineCount = 0
			currentChar++
			scanner.Scan() // Consume the empty line between characters
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading banner file: %w", err)
	}

	return bannerMap, nil
}

// renderString converts the input string into ASCII art using the banner map.
func renderString(input string, bannerMap map[rune][]string) ([]string, error) {
	outputLines := make([]string, bannerHeight)

	for _, char := range input {
		art, ok := bannerMap[char]
		if !ok {
			return nil, fmt.Errorf("character '%c' not found in banner", char)
		}
		for i := 0; i < bannerHeight; i++ {
			outputLines[i] += art[i]
		}
	}

	return outputLines, nil
}

// getAvailableBanners lists all the banners available in the banners directory.
func getAvailableBanners(dirPath string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var banners []string
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(name, ".txt") {
				bannerName := strings.TrimSuffix(name, ".txt")
				banners = append(banners, bannerName)
			}
		}
	}
	return banners, nil
}

// trimLeadingEmptyLines removes empty lines from the beginning of the output.
func trimLeadingEmptyLines(lines []string) []string {
	startIndex := 0
	for i, line := range lines {
		if line != "" {
			startIndex = i
			break
		}
	}
	return lines[startIndex:]
}
