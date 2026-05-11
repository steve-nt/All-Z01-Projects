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

	// Get the input string and banner name from command-line arguments
	inputString := os.Args[1]
	bannerName := "standard" // Default banner
	if len(os.Args) > 2 {
		bannerName = os.Args[2]
	}

	// Load the banner file
	bannerFilePath := filepath.Join("banners", bannerName+".txt")
	bannerMap, err := createBannerMap(bannerFilePath)
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
	if inputString == "" {
		return
	}
	// Replace occurrences of '\\n' with '\n' to handle newlines
	inputString = strings.ReplaceAll(inputString, "\\n", "\n")

	if inputString == "\n" {
		fmt.Println()
		return
	}
	// Split the input string by newline characters
	lines := strings.Split(inputString, "\n")

	// Process each line separately
	for idx, line := range lines {
		if idx > 1 {
			// Print a new line between lines
			fmt.Println()
		}
		// Handle empty lines
		if line == "" {
			continue
		}
		// Generate ASCII art for the line
		for i := 0; i < bannerHeight; i++ {
			var outputLine string
			for _, ch := range line {
				artLines, ok := bannerMap[ch]
				if !ok {
					// Use space as a placeholder for unsupported characters
					artLines = bannerMap[' ']
				}
				outputLine += artLines[i]
			}
			fmt.Println(outputLine)
		}
	}
}

// createBannerMap loads the banner file and returns a map of rune to ASCII art lines
func createBannerMap(bannerFileName string) (map[rune][]string, error) {
	bannerFile, err := os.Open(bannerFileName)
	if err != nil {
		return nil, err
	}
	defer bannerFile.Close()

	bannerMap := make(map[rune][]string)
	scanner := bufio.NewScanner(bannerFile)

	// ASCII codes from space (32) to tilde (126)
	startChar := 32
	endChar := 126

	// Skip the first line if it's empty
	if scanner.Scan() {
		if scanner.Text() == "" {
			// Continue to the next line
		} else {
			// Reset scanner to the beginning
			bannerFile.Seek(0, 0)
			scanner = bufio.NewScanner(bannerFile)
		}
	}

	// Read the ASCII art for each character
	for asciiCode := startChar; asciiCode <= endChar; asciiCode++ {
		var artLines []string
		for i := 0; i < bannerHeight; i++ {
			if scanner.Scan() {
				artLines = append(artLines, scanner.Text())
			} else {
				return nil, fmt.Errorf("unexpected EOF while reading ASCII art for character %d", asciiCode)
			}
		}
		// Add the character and its ASCII art to the map
		bannerMap[rune(asciiCode)] = artLines
		// Skip the empty line between characters
		scanner.Scan()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return bannerMap, nil
}

// getAvailableBanners lists the available banners in the banners directory
func getAvailableBanners(bannersDir string) ([]string, error) {
	filesDir, err := os.ReadDir(bannersDir)
	if err != nil {
		return nil, err
	}

	var banners []string
	for _, file := range filesDir {
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
