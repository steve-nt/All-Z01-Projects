package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"ascii-art-web/services"
)

// AsciiHandler handles HTTP requests for ASCII art generation.
// It contains a reference to the ASCII art service for processing requests.
type AsciiHandler struct {
	service *services.AsciiArtWeb
}

// PageData contains data passed to HTML templates for rendering the main page.
// It includes user input, generated ASCII art, and available banner options.
type PageData struct {
	InputText   string   // User's input text
	InputBanner string   // Selected banner style
	AsciiArt    string   // Generated ASCII art result
	Banners     []string // List of available banner styles
}

// ErrorData contains data passed to error page templates.
// It includes the HTTP status code and error message to display.
type ErrorData struct {
	StatusCode int    // HTTP status code (400, 404, 500, etc.)
	Message    string // Error message to display to user
}

// NewAsciiHandler creates and returns a new AsciiHandler instance.
// It takes an ASCII art service as dependency for processing requests.
func NewAsciiHandler(service *services.AsciiArtWeb) *AsciiHandler {
	return &AsciiHandler{
		service: service,
	}
}

// HandleHome serves the main page with the ASCII art generation form.
// It only accepts GET requests to the root path "/" and renders the home template
// with available banner options and a default banner selection.
func (a *AsciiHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - only GET allowed
	if r.Method != http.MethodGet {
		a.HandleErrors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Validate URL path - only root path allowed
	if r.URL.Path != "/" {
		a.HandleErrors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	// Load and parse the home page template
	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		a.HandleErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Prepare template data with available banners and default selection
	data := PageData{
		InputBanner: "standard",                      // Default banner selection
		Banners:     a.service.GetAvailableBanners(), // All available banner options
	}

	// Render the template with the prepared data
	if err := pageTemplate.Execute(w, data); err != nil {
		a.HandleErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

// HandleAsciiArt processes form submissions and generates ASCII art.
// It accepts POST requests with form data containing text and banner selection,
// validates the input, generates ASCII art, and returns the result page.
func (a *AsciiHandler) HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method - only POST allowed
	if r.Method != http.MethodPost {
		a.HandleErrors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Parse form data from the request body
	if err := r.ParseForm(); err != nil {
		a.HandleErrors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Extract text input from form data
	text := r.FormValue("text")
	if text == "" {
		a.HandleErrors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Extract banner selection from form data
	banner := r.FormValue("banner")

	// Validate input using service layer validation rules
	if err := a.service.ValidateInput(text, banner); err != nil {
		a.HandleErrors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	// Generate ASCII art using the service layer
	asciiArt, err := a.service.Generate(text, banner)
	if err != nil {
		a.HandleErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Load and parse the home template for displaying results
	pageTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		a.HandleErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	// Prepare template data with input, result, and available banners
	data := PageData{
		InputText:   text,                            // Original user input
		InputBanner: banner,                          // Selected banner style
		AsciiArt:    asciiArt,                        // Generated ASCII art
		Banners:     a.service.GetAvailableBanners(), // All available banners
	}

	// Render the template with the result data
	if err := pageTemplate.Execute(w, data); err != nil {
		a.HandleErrors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (a *AsciiHandler) HandleResources(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		a.HandleErrors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	filePath := strings.TrimPrefix(r.URL.Path, "/static")
	http.ServeFile(w, r, filepath.Join("static", filePath))
}

// HandleErrors renders error pages for various HTTP error conditions.
// It takes a status code and message, loads the error template, and displays
// a formatted error page. Falls back to plain text if template loading fails.
func (a *AsciiHandler) HandleErrors(w http.ResponseWriter, statusCode int, message string) {
	// Attempt to load the error page template
	errorTemplate, err := template.ParseFiles("templates/error.html")
	if err != nil {
		// Fallback to plain text error if template fails to load
		http.Error(w, fmt.Sprintf("Error %d: %s", statusCode, http.StatusText(statusCode)), statusCode)
		return
	}

	// Set the HTTP status code in the response
	w.WriteHeader(statusCode)

	// Render the error template with status code and message
	errorData := ErrorData{
		StatusCode: statusCode,
		Message:    message,
	}

	// Execute template or fallback to plain text if rendering fails
	if err := errorTemplate.Execute(w, errorData); err != nil {
		http.Error(w, fmt.Sprintf("Error %d: %s", statusCode, http.StatusText(statusCode)), statusCode)
		return
	}
}
