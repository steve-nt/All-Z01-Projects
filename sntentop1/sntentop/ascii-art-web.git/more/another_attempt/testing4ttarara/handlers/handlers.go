package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var tpl *template.Template
var str string       // Stores the ASCII art result temporarily
var history []string // Stores the history of all submissions and errors

func InitializeTemplates(templatePath string, logger *log.Logger) (*template.Template, error) {
	pattern := templatePath + "/*.html"
	parsedTemplates, err := template.ParseGlob(pattern)
	if err != nil {
		logger.Printf("Failed to parse templates from path %s: %v", pattern, err)
		return nil, err
	}

	logger.Printf("Templates successfully parsed from path: %s", pattern)
	tpl = parsedTemplates
	return tpl, nil
}

func SetTemplates(t *template.Template) {
	tpl = t
}

// Index serves the home page or the 404 page if the path is invalid
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		logEvent("404 Error", "Page not found: "+r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		err := tpl.ExecuteTemplate(w, "404.html", nil)
		if err != nil {
			fmt.Println("Error rendering 404 template:", err)
		}
		return
	}

	err := tpl.ExecuteTemplate(w, "index.html", struct {
		First string
		Error string
	}{
		First: "",
		Error: "",
	})
	if err != nil {
		handleServerError(w, "Failed to render index template")
	}
}

// Processor handles the /ascii-art route
func Processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logEvent("405 Error", "Invalid HTTP method: "+r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("ascii-data")
	font := r.FormValue("fonts")

	// Validate and process input
	final, renderStatus := render(text)
	if renderStatus == 400 {
		logEvent("400 Error", "Invalid Input: "+text)
		handleBadRequest(w, "Invalid Input: "+text)
		return
	}

	if font != "standard" && font != "shadow" && font != "thinkertoy" {
		logEvent("400 Error", "Invalid Font Selected: "+font)
		handleBadRequest(w, "Invalid font type. Choose between 'standard', 'shadow', or 'thinkertoy'.")
		return
	}

	// Generate ASCII art
	data, statusCode := asciiArt(final, font)
	if statusCode == 500 {
		logEvent("500 Error", "Failed to read font file for font: "+font)
		handleServerError(w, "Failed to generate ASCII art")
		return
	}

	// Log success and render the result
	logEvent("200: Success", fmt.Sprintf("Font: %s, Input: %s", font, text))
	err := tpl.ExecuteTemplate(w, "index.html", struct {
		First string
		Error string
	}{
		First: data,
		Error: "",
	})
	if err != nil {
		handleServerError(w, "Failed to render result template")
	}
}

// asciiArt generates the ASCII art using the specified font
func asciiArt(argument string, fonts string) (string, int) {
	//func asciiArt(argument string, fonts string, readFile func(string) ([]byte, error)) (string, int) {
	banner, err := os.ReadFile("static/fonts/" + fonts + ".txt")
	if err != nil {
		return "", 500
	}

	lines := strings.Split(string(banner), "\n")
	if fonts == "thinkertoy" {
		lines = strings.Split(string(banner), "\r\n")
	}

	str = ""
	myLines := strings.Split(strings.ReplaceAll(argument, "\r", ""), "\\n")
	for _, line := range myLines {
		for k := 0; k < 8; k++ {
			for i := 0; i < len(line); i++ {
				if int(line[i]) < 32 || int(line[i]) > 126 {
					return "", 500
				}
				str += lines[(int(line[i])-32)*9+1+k]
			}
			str += "\n"
		}
	}
	return str, 200
}

// render validates and processes the input text
func render(input string) (string, int) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "Input cannot be empty or just whitespace.", 400
	}

	if len(input) >= 128 {
		return "Input too long. Maximum allowed length is 128 characters.", 400
	}

	if !validInput(input) {
		return "Input contains invalid characters. Only ASCII characters, tabs, and newlines are allowed.", 400
	}

	return strings.ReplaceAll(input, "\r\n", "\\n"), 200
}

// validInput checks if the input contains valid ASCII characters
func validInput(input string) bool {
	for _, char := range input {
		if (char < 32 || char > 126) && char != 9 && char != 10 && char != 13 {
			return false
		}
	}
	return true
}

// handleBadRequest renders the 400 error page
func handleBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	err := tpl.ExecuteTemplate(w, "400.html", struct {
		Message string
	}{
		Message: message,
	})
	if err != nil {
		fmt.Println("Error rendering 400 template:", err)
	}
}

// handleServerError renders the 500 error page
func handleServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	err := tpl.ExecuteTemplate(w, "500.html", struct {
		Message string
	}{
		Message: message,
	})
	if err != nil {
		fmt.Println("Error rendering 500 template:", err)
	}
}

// logEvent logs an event (success or error) to the history and prints the entire history
func logEvent(eventType, message string) {
	entry := fmt.Sprintf("[%s] %s: %s", time.Now().Format("2006-01-02 15:04:05"), eventType, message)
	history = append(history, entry)

	// Print the full history
	fmt.Println("\n===== Submission and Error History =====")
	for _, h := range history {
		fmt.Println(h)
	}
}
