package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Template cache
var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	templates.ExecuteTemplate(w, "index.html", nil)
}

// ASCII Art handler
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" || banner == "" {
		http.Error(w, "Bad Request: Missing text or banner", http.StatusBadRequest)
		return
	}

	result, err := generateAsciiArt(text, banner)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "result.html", map[string]string{
		"Result": result,
	})
}

// ASCII Art generation function
func generateAsciiArt(text, banner string) (string, error) {
	filePath := fmt.Sprintf("banners/%s.txt", banner)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("banner file not found")
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read banner file")
	}

	lines := strings.Split(string(content), "\n")
	result := ""

	for _, char := range text {
		if char < 32 || char > 126 {
			return "", fmt.Errorf("unsupported character: %q", char)
		}
		start := (int(char) - 32) * 8
		for i := 0; i < 9; i++ {
			result += lines[start+i] + "\n"
		}
	}

	return result, nil
}
