package repository

import (
	"bufio"
	"os"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/domain"
)

const (
	spaceASCII = 32
	tidleASCII = 126
	charHeight = 8
)

// Define an interface
type BannerRepositoryInterface interface {
	LoadBanner(name string) (*domain.Banner, error)
}

// Ensure BannerRepository implements the interface
type BannerRepository struct {
	basePath string
}

func NewBannerRepository(path string) *BannerRepository {
	return &BannerRepository{
		basePath: path,
	}
}

// LoadBanner now implements the interface
func (r *BannerRepository) LoadBanner(name string) (*domain.Banner, error) {
	bannerPath := "../../" + r.basePath + name + ".txt"
	banner, err := os.Open(bannerPath)
	if err != nil {
		return nil, err
	}
	defer banner.Close()

	asciiMap := make(map[rune][]string)
	char := []string{} // Holds multiple lines representing the character

	scanner := bufio.NewScanner(banner)
	startOfCharacter := false // Switch to indicate the NEXT line is the start of a character
	sliceNum := 0             // Number of currently processing line
	key := ' '                // Indicates the first writable character in ascii

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 && !startOfCharacter { // Found start of character text
			startOfCharacter = true // Set switch on
			continue
		} else {
			char = append(char, line) // Fill the slice with subsequent lines
			sliceNum++                // Mark how many lines we appended
		}
		if sliceNum == charHeight { // Stop when 8 lines appended
			asciiMap[key] = char     // Push it to the map
			key++                    // Change to the next ascii character
			startOfCharacter = false // Set switch off
			sliceNum = 0             // Reset the counter of lines we already have
			char = []string{}        // Clear the text we have appended
		}
	}

	return &domain.Banner{Name: name, Lines: asciiMap}, nil
}
