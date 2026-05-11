package handlers

import (
	"html/template" // Provides functions to parse and execute HTML templates securely.
	"log"           // Provides simple logging functions for outputting messages.
	"net/http"      // Provides HTTP client and server implementations.
	"os"            // Provides functions for interacting with the operating system, like reading files.
	"strings"       // Provides functions to manipulate strings.
)

// PrepareTemplateData creates a data map for HTML templates.
// title: The title for the page.
// pageClass: A CSS class to style the page.
func PrepareTemplateData(title, pageClass string) map[string]interface{} {
	//To give the title to a page and select the class for Index.html to access the right css style
	// Returns a map with the provided title, pageClass, and placeholders for additional data.
	return map[string]interface{}{
		"Title":     title,     // The page title.
		"PageClass": pageClass, // CSS class for styling the page.
		"First":     "",        // Placeholder for additional data.
		"Error":     "",        // Placeholder for error messages.
	}
}

// RenderTemplateFiles renders HTML templates and sends them to the HTTP response.
// w: HTTP response writer to send data to the client.
// templateName: The main template file to render.
// data: Dynamic data to inject into the template.
func RenderTemplateFiles(w http.ResponseWriter, templateName string, data interface{}) {
	//w.Header().Set("Content-Type", "text/html")
	// Log the template rendering process.
	log.Println("Rendering template:", templateName)

	// Parse the specified template and a common base template (Index.html).
	tmpl, err := template.ParseFiles("templates/"+templateName, "templates/Index.html")
	if err != nil {
		// Log an error if template parsing fails.
		log.Printf("Status 500: Error loading templates: %v\n", err)
		// na mpainei sto handler tou [Status 500] - handleServerError()
		return
	}
	log.Println("Status 200: Loading templates successfully")

	// Execute the template with the provided data and send it to the HTTP response.
	err = tmpl.ExecuteTemplate(w, "Index", data)
	if err != nil {
		// Log an error if template rendering fails.
		log.Printf("Status 500: Error rendering template: %v\n", err)
		// na mpainei sto handler tou [Status 500] - handleServerError()
	}
}

// _________________________________________________________________________________
// inputValidation validates user input and returns a sanitized version or an error message.
// input: The input string to validate.
// Returns a sanitized string and an HTTP status code.
func inputValidation(input string) (string, int) {
	// Remove leading and trailing spaces from the input.
	input = strings.TrimSpace(input)

	// Check if the input is empty or only whitespace.
	if input == "" {
		return "Input cannot be empty or just whitespace.", 400
	}
	// Ensure the input is not too long (max 128 characters).
	if len(input) >= 128 {
		return "Input too long. Maximum allowed length is 128 characters.", 400
	}
	// Validate that the input contains only valid ASCII characters.
	if !asciiValidation(input) {
		return "Input contains invalid characters. Only ASCII characters, tabs, and newlines are allowed.", 400
	}
	// Replace Windows-style newlines (\r\n) with Unix-style (\n).
	return strings.ReplaceAll(input, "\r\n", "\\n"), 200
}

// asciiValidation checks if a string contains only valid ASCII characters.
// input: The string to validate.
// Returns true if the string is valid, false otherwise.
func asciiValidation(input string) bool {
	// Iterate over each character in the string.
	for _, char := range input {
		// Validate ASCII characters and control characters (tabs, newlines).
		if (char < 32 || char > 126) && char != 9 && char != 10 && char != 13 {
			return false
		}
	}
	return true
}

// _________________________________________________________________________________
// asciiArt generates ASCII art using a specific font and input string.
// argument: The text to convert into ASCII art.
// fonts: The font style to use for the ASCII art.
// Returns the ASCII art string and an HTTP status code.
func asciiArt(argument string, fonts string) (string, int) {
	// Read the font file from the static/banners directory.
	bannerStyle, err := os.ReadFile("static/banners/" + fonts + ".txt")
	if err != nil {
		// Return an error if the font file cannot be read.
		return "", 500
	}
	// Normalize line endings in the font file.
	normalizedBanner := normalizeNewlines(string(bannerStyle))
	lines := strings.Split(normalizedBanner, "\n") // Split the font file into lines.
	str := ""                                      // Initialize the result string.

	// Split the input text into lines based on the custom newline format (\n).
	myLines := strings.Split(strings.ReplaceAll(argument, "\r", ""), "\\n")

	// Iterate over each line of input text.
	for _, line := range myLines {
		// Create ASCII art for each character in the line.
		for k := 0; k < 8; k++ {
			for i := 0; i < len(line); i++ {
				// Ensure each character is a valid ASCII printable character.
				if int(line[i]) < 32 || int(line[i]) > 126 {
					return "", 500
				}
				// Append the corresponding ASCII art row for the character.
				str += lines[(int(line[i])-32)*9+1+k]
			}
			str += "\n" // Add a newline after each row.
		}
	}
	return str, 200 // Return the ASCII art and status code.
}

// normalizeNewlines converts different newline formats to Unix-style (\n).
// content: The string to normalize.
// Returns the normalized string.
func normalizeNewlines(content string) string {
	// Replace Windows-style line endings (\r\n) with Unix-style (\n)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	// Replace old Mac-style line endings (\r) with Unix-style (\n)
	content = strings.ReplaceAll(content, "\r", "\n")
	return content
}
