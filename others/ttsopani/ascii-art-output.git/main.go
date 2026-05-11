package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	usageMsg      = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . --output=<fileName.txt> something standard"
	asciiArtLines = 8
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println(usageMsg)
		return
	}

	var outputFile string
	var inputText, banner string
	var err error

	// Parse flags and arguments
	for _, arg := range args {
		if strings.HasPrefix(arg, "--output") {
			outputFile, err = parseOutputFlag(arg)
			if err != nil {
				fmt.Println(usageMsg)
				return
			}
		} else if inputText == "" {
			inputText = arg
		} else {
			banner = arg
		}
	}

	if banner == "" {
		banner = "standard"
	}

	if inputText == "" {
		fmt.Println(usageMsg)
		return
	}

	// Load ASCII art map
	bannerMap, err := loadBannerFile(banner)
	if err != nil {
		fmt.Println("Error loading banner:", err)
		return
	}

	// Generate ASCII art
	asciiArt, err := generateAsciiArt(inputText, bannerMap, asciiArtLines)
	if err != nil {
		fmt.Println("Error generating ASCII art:", err)
		return
	}

	// Output to file or console
	if outputFile != "" {
		err = writeToFile(outputFile, asciiArt)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Printf("ASCII art written to file: %s\n", outputFile)
	} else {
		fmt.Println(asciiArt)
	}
}

func parseOutputFlag(flag string) (string, error) {
	if !strings.HasPrefix(flag, "--output=") {
		return "", errors.New("invalid output flag format")
	}
	return strings.TrimPrefix(flag, "--output="), nil
}

func loadBannerFile(banner string) (map[rune][]string, error) {
	fileName := fmt.Sprintf("%s.txt", banner)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bannerMap := make(map[rune][]string)
	scanner := bufio.NewScanner(file)

	var currentRune rune = 32
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(lines) == asciiArtLines {
				bannerMap[currentRune] = lines
				currentRune++
				lines = nil
			}
		} else {
			lines = append(lines, line)
		}
	}

	if len(lines) == asciiArtLines {
		bannerMap[currentRune] = lines
	}

	return bannerMap, scanner.Err()
}

func generateAsciiArt(input string, bannerMap map[rune][]string, lineHeight int) (string, error) {
	var result []string
	newSlice := strings.Split(input, "\\n")
	for j := 0; j < len(newSlice); j++ {
		for i := 0; i < lineHeight; i++ {
			line := ""
			for _, char := range newSlice[j] {
				if ascii, exists := bannerMap[char]; exists {
					line += ascii[i]
				} else {
					return "", fmt.Errorf("character '%c' not found in ASCII art map", char)
				}
			}
			result = append(result, line)
		}
		result = append(result, "")
	}
	return strings.Join(result, "\n"), nil
}

func writeToFile(fileName, content string) error {
	file, err := os.Create(fileName) // Creates or overwrites the file
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	return writer.Flush()
}
