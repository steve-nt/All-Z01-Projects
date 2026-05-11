package usecase

import (
	"fmt"
	"strings"

	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/adapter/repository"
	"platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/domain"
)

const (
	charHeight = 8
)

type AsciiUsecase struct {
	repo repository.BannerRepositoryInterface
}

func NewAsciiUsecase(repo repository.BannerRepositoryInterface) *AsciiUsecase {
	return &AsciiUsecase{repo: repo}
}

// Not making use of AsciiUsecase, maybe make this a domain.Banner method?
func (u *AsciiUsecase) ConvertTextToAscii(asciiTextRequest *domain.ASCIITextRequest) (string, error) {
	domainBanner, err := u.repo.LoadBanner(asciiTextRequest.Banner)
	if err != nil {
		return "", err
	}

	// Split input into lines (handle \n)
	text := strings.ReplaceAll(asciiTextRequest.Text, "\\n", "\n")
	lines := strings.Split(text, "\n")
	var output strings.Builder

	if len(text) == 0 {
		return "", fmt.Errorf("empty text")
	}

	for _, line := range lines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		// Build each of the 8 ASCII lines
		asciiLines := make([]string, charHeight)
		for _, char := range line {
			art, exists := domainBanner.Lines[char]
			if !exists {
				return "", fmt.Errorf("character %q not supported in banner", char)
			}
			for i := 0; i < charHeight; i++ {
				asciiLines[i] += art[i]
			}
		}

		// Append the ASCII lines to the output
		for _, l := range asciiLines {
			output.WriteString(l + "\n")
		}
	}

	return output.String(), nil
}
