package web

import (
	"html/template"
	"net/http"

	"ascii/ascii"
)

type PageData struct {
	AsciiArt string
	Text     string
	Banner   string
}

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse form data
		text := r.FormValue("text")
		banner := r.FormValue("banner")

		// Generate ASCII art
		asciiArt, err := ascii.GenerateAsciiArt(text, "banners/"+banner+".txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the template with the generated ASCII art
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{
			AsciiArt: asciiArt,
			Text:     text,
			Banner:   banner,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == "GET" { // Handle GET requests

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := PageData{}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
