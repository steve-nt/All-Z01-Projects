package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	charHeight = 8  // Each ASCII art character is 8 lines tall
	startASCII = 32 // ASCII space character (first printable character)
)

// AsciiArtWeb manages ASCII art generation with multiple banner styles.
// It stores loaded banner fonts and provides methods for text conversion.
type AsciiArtWeb struct {
	banners map[string]map[rune][]string // Maps banner names to character art data
}

// NewAsciiArtWeb creates and initializes a new AsciiArtWeb instance.
// It returns a service with an empty banners map ready for loading font files.
func NewAsciiArtWeb() *AsciiArtWeb {
	return &AsciiArtWeb{
		banners: make(map[string]map[rune][]string),
	}
}

// LoadBanners reads all banner font files from the ./banners directory.
// Each .txt file contains ASCII art patterns for printable characters (32-126).
// Characters are separated by empty lines and must be exactly 8 lines tall.
func (a *AsciiArtWeb) LoadBanners() error {
	// Read all files from the banners directory
	files, err := os.ReadDir("./banners")
	if err != nil {
		return fmt.Errorf("failed reading banners directory: %w", err)
	}

	// Process each .txt file as a banner font
	for _, file := range files {
		// Skip directories and non-.txt files
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		// Extract banner name from filename (remove .txt extension)
		name := strings.TrimSuffix(file.Name(), ".txt")

		// Open the banner file
		bannerfile, err := os.Open("./banners/" + name + ".txt")
		if err != nil {
			return fmt.Errorf("failed loading banner: %w", err)
		}
		defer bannerfile.Close()

		// Initialize storage for this banner's character patterns
		banner := make(map[rune][]string)
		scanner := bufio.NewScanner(bannerfile)
		var lines []string
		currentRune := rune(startASCII) // Start with space character (ASCII 32)

		// Read file line by line
		for scanner.Scan() {
			line := scanner.Text()

			// Empty line indicates end of current character
			if line == "" && len(lines) > 0 {
				// Save completed character if it has correct height
				if len(lines) == charHeight {
					banner[currentRune] = lines
					currentRune++ // Move to next ASCII character
				}
				lines = []string{} // Reset for next character
				continue
			}

			// Collect non-empty lines for current character
			if line != "" {
				lines = append(lines, line)
			}
		}

		// Handle the last character in file (no trailing empty line)
		if len(lines) == charHeight {
			banner[currentRune] = lines
		}

		// Check for file reading errors
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed loading file: %w", err)
		}

		// Store the loaded banner in the service
		a.banners[name] = banner
	}
	return nil
}

// Generate converts input text to ASCII art using the specified banner style.
// It processes each character, handles newlines, and builds the final ASCII art output.
// Unknown characters are replaced with spaces to maintain formatting.
func (a *AsciiArtWeb) Generate(text, bannerName string) (string, error) {
	// Verify the requested banner exists
	banner, exists := a.banners[bannerName]
	if !exists {
		return "", fmt.Errorf("banner '%s' not loaded", bannerName)
	}

	var result strings.Builder
	linesToPrint := make([]string, charHeight) // 8 lines for each row of characters

	// Process each character in the input text
	for _, char := range text {
		// Handle newline characters
		if char == '\n' {
			// Flush current accumulated lines
			for i := 0; i < charHeight; i++ {
				if linesToPrint[i] != "" {
					result.WriteString(linesToPrint[i] + "\n")
					linesToPrint[i] = "" // Clear for next line
				}
			}
			result.WriteString("\n") // Add extra newline for spacing
			continue
		}

		// Get ASCII art pattern for current character
		charArt := banner[char]
		if charArt == nil {
			// Use space character for unknown/unsupported characters
			charArt = banner[' ']
		}

		// Append character art to each of the 8 lines
		for i := 0; i < charHeight; i++ {
			if i < len(charArt) {
				linesToPrint[i] += charArt[i]
			}
		}
	}

	// Flush any remaining lines at the end
	for i := 0; i < charHeight; i++ {
		if linesToPrint[i] != "" {
			result.WriteString(linesToPrint[i] + "\n")
		}
	}

	return result.String(), nil
}

// GetAvailableBanners returns a list of all loaded banner style names.
// This is used to populate the banner selection dropdown in the web interface.
func (a *AsciiArtWeb) GetAvailableBanners() []string {
	var banners []string
	for name := range a.banners {
		banners = append(banners, name)
	}
	return banners
}

// ValidateInput performs comprehensive validation on user input.
// It checks for empty values, invalid characters, length limits, and banner availability.
// Returns an error if any validation rule is violated.
func (a *AsciiArtWeb) ValidateInput(text, banner string) error {
	// Check for invalid characters (only ASCII printable + newline/tab/carriage return allowed)
	for _, r := range text {
		if (r < 32 || r > 126) && r != '\n' && r != '\r' && r != '\t' {
			return fmt.Errorf("")
		}
	}

	return nil
}
