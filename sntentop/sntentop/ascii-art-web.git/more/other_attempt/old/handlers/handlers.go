package handlers

import (
	"ascii-art-web/ascii"
	"html/template"
	"net/http"
)

// HomePageHandler handles the GET request to the home page ("/")
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AsciiArtHandler handles the POST request to "/ascii-art"
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Read the banner file
	fileLines := ascii.ReadTxt(banner)
	if fileLines == nil {
		http.Error(w, "The requested banner file could not be found.", http.StatusNotFound)
		return
	}

	// Generate ASCII art
	asciiTemplates := ascii.Return2DASCIIArray(fileLines)
	output := ascii.ReturnAllStringASCII(text, asciiTemplates)

	// Render the result.html template
	data := map[string]string{
		"AsciiArt": output,
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
