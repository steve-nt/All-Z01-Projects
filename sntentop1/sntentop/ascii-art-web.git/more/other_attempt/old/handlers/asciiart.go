package handlers

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Templates is exported to be initialized in main.go.
var Templates *template.Template

// HomePageHandler serves the main page.
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	Templates.ExecuteTemplate(w, "index.html", nil)
}

// AsciiArtHandler processes ASCII art generation.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	font := r.FormValue("banner")
	if text == "" || font == "" {
		http.Error(w, "Missing text or font selection", http.StatusBadRequest)
		return
	}

	if !isValidASCIIInput(text) {
		http.Error(w, "Invalid input. Only ASCII characters are allowed.", http.StatusBadRequest)
		return
	}

	asciiArt, err := generateASCIIArt(text, font)
	if err != nil {
		http.Error(w, "Error generating ASCII art: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		AsciiArt string
	}{
		AsciiArt: asciiArt,
	}
	Templates.ExecuteTemplate(w, "result.html", data)
}

// isValidASCIIInput checks if the input contains only printable ASCII characters.
func isValidASCIIInput(input string) bool {
	for _, char := range input {
		if char < 32 || char > 126 { // Check printable ASCII range
			return false
		}
	}
	return true
}

// generateASCIIArt generates ASCII art using the specified font.
func generateASCIIArt(text, font string) (string, error) {
	file, err := os.Open("static/fonts/" + font + ".txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	asciiMap := map[int][]string{}
	asciiIndex := 31
	for _, line := range lines {
		if line == "" {
			asciiIndex++
		} else {
			asciiMap[asciiIndex] = append(asciiMap[asciiIndex], line)
		}
	}

	result := buildASCIIArt(text, asciiMap)
	return result, nil
}

// buildASCIIArt constructs ASCII art line by line.
func buildASCIIArt(text string, asciiMap map[int][]string) string {
	var result strings.Builder
	for i := 0; i < len(asciiMap[32]); i++ { // Rows per character
		for _, char := range text {
			if lines, ok := asciiMap[int(char)]; ok {
				result.WriteString(lines[i])
			}
		}
		result.WriteString("\n")
	}
	return result.String()
}
