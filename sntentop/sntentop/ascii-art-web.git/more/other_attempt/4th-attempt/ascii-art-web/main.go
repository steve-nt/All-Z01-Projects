package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

// ASCII banners map
var banners = map[string]string{
	"standard":   "banners/standard.txt",
	"shadow":     "banners/shadow.txt",
	"thinkertoy": "banners/thinkertoy.txt",
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	fmt.Println("Server is running on http://localhost:8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Parse form data
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Simulate Internal Server Error when user enters "error500"
	if strings.ToLower(text) == "error500" {
		renderError500(w)
		return
	}

	// Validate input
	if text == "" || banner == "" {
		http.Error(w, "Text and banner are required", http.StatusBadRequest)
		return
	}

	// Load ASCII banner and generate output
	result, err := generateASCIIArt(text, banners[banner])
	if err != nil {
		renderError500(w)
		return
	}

	// Render result
	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, map[string]string{"Result": result})
}

func renderError500(w http.ResponseWriter) {
	// Parse the 500 error template
	tmpl, err := template.ParseFiles("templates/500.html")
	if err != nil {
		// Fallback to a simple plain-text error message if template fails
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Write the HTTP status code and render the template
	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Execute(w, nil)
}

// generateASCIIArt generates the ASCII art based on the input text and banner file path.
func generateASCIIArt(inputText, bannerPath string) (string, error) {
	// Read the banner file
	content, err := ioutil.ReadFile(bannerPath)
	if err != nil {
		return "", errors.New("failed to load banner file")
	}

	lines := strings.Split(string(content), "\n")
	output := ""

	// Split input text into lines
	inputLines := strings.Split(inputText, "\n")

	// Process each line of input text
	for _, line := range inputLines {
		// ASCII art contains 8 lines per character block
		for row := 0; row < 8; row++ { // Process 8 rows for each character
			for _, char := range line {
				// Only allow printable ASCII characters (32 to 126)
				if char < 32 || char > 126 {
					continue // Skip invalid characters
				}

				// Find the starting line for this character
				startLine := int(char-32) * 9
				if startLine+row < len(lines) {
					output += lines[startLine+row]
				}
			}
			output += "\n" // Move to the next line of ASCII art
		}
	}

	return output, nil
}
