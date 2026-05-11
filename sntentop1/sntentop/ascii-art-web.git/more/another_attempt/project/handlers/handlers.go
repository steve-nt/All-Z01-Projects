package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// [ Page Handlers ] ------------------------------------------------------------------------

func PageHOME(w http.ResponseWriter, r *http.Request) {
	data := prepareData(".HOME", "home-page")
	log.Println("Handling_Request: Go_to_Page: HOME")
	renderTemplate(w, "home.html", data)
}

func PageCONVERTER(w http.ResponseWriter, r *http.Request) {
	data := prepareData(".CONVERTER", "converter-page")
	log.Println("Handling_Request: Go_to_Page: CONVERTER")
	renderTemplate(w, "converter.html", data)
}

// [ Template Parsing ] --------------------------------------------------------------------

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	log.Println("Rendering template:", templateName)

	tmpl, err := template.ParseFiles("templates/"+templateName, "templates/index.html")
	if err != nil {
		log.Printf("Error loading templates: %v\n", err)
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		return
	}
	log.Println("Templates loaded successfully.")

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func prepareData(title, pageClass string) map[string]interface{} {
	return map[string]interface{}{
		"Title":     title,
		"PageClass": pageClass,
	}
}

// [ Input Parsing ] --------------------------------------------------------------------

func InputHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to process form", http.StatusBadRequest)
			return
		}

		// Get the input values
		text := r.FormValue("text")
		banner := r.FormValue("banner")

		// Print the inputs to the terminal
		fmt.Printf("Text: %s\n", text)
		fmt.Printf("Banner: %s\n", banner)

		// Respond back to the user
		fmt.Fprintf(w, "Input received. Check the terminal for details.")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
