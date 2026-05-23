package ascii

import (
	"fmt"
	"os"
	"strings"
)

const (
	firstPrintableASCII = 32
	lastPrintableASCII  = 126
	charHeight          = 8
	blockSize           = 9
)

// LoadBanner reads a banner file by name, for example "standard", from the
// banners/ directory and returns each printable ASCII character's 8-line block.
func LoadBanner(name string) (map[rune][charHeight]string, error) {
	path := fmt.Sprintf("banners/%s.txt", name)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open banner %q: %w", name, err)
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	banner := make(map[rune][charHeight]string)

	for ch := firstPrintableASCII; ch <= lastPrintableASCII; ch++ {
		blockIndex := ch - firstPrintableASCII
		start := 1 + blockIndex*blockSize

		if start+charHeight > len(lines) {
			return nil, fmt.Errorf("invalid banner %q: missing character %q", name, rune(ch))
		}

		var block [charHeight]string
		for row := 0; row < charHeight; row++ {
			block[row] = lines[start+row]
		}

		banner[rune(ch)] = block
	}

	return banner, nil
}

// Render converts input text into ASCII art. Literal "\\n" sequences are treated
// as line breaks, matching how the shell passes arguments such as "Hello\\nThere".
func Render(input string, banner map[rune][charHeight]string) string {
	segments := strings.Split(input, "\\n")

	allEmpty := true
	for _, segment := range segments {
		if segment != "" {
			allEmpty = false
			break
		}
	}

	if allEmpty {
		var sb strings.Builder
		for i := 0; i < len(segments)-1; i++ {
			sb.WriteByte('\n')
		}
		return sb.String()
	}

	var sb strings.Builder

	for _, segment := range segments {
		if segment == "" {
			sb.WriteByte('\n')
			continue
		}

		for row := 0; row < charHeight; row++ {
			for _, ch := range segment {
				block, ok := banner[ch]
				if !ok {
					continue
				}

				sb.WriteString(block[row])
			}

			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
