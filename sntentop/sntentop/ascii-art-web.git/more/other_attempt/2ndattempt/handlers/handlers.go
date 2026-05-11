// handlers.go
package handlers

import (
	"ascii-art-web/ascii"
	"html/template"
	"io/ioutil"
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

	// 1. Read the banner file
	fileLines := ascii.ReadTxt(banner)
	if fileLines == nil {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}

	// 2. Generate ASCII art
	asciiTemplates := ascii.Return2DASCIIArray(fileLines)
	output := ascii.ReturnAllStringASCII(text, asciiTemplates)

	// 3. Save the output to a file (optional)
	err := ioutil.WriteFile("ascii_art_output.txt", []byte(output), 0644)
	if err != nil {
		http.Error(w, "Error saving output to file", http.StatusInternalServerError)
		return
	}

	// 4. Render the result.html template
	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, output) // Pass the output to the template
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
