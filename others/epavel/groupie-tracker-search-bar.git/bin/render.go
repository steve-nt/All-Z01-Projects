package bin

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// renderTemplate is a helper function to load and render html templates
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Define a custom function to replace spaces with hyphens within html templates
	funcMap := template.FuncMap{
		"replaceSpaces": func(s string) string {
			return strings.ReplaceAll(s, " ", "-")
		},
	}
	// Load the home template
	t, err := template.New(tmpl).Funcs(funcMap).ParseFiles("templates/" + tmpl)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// renderErrorTemplate is a helper function to load and render error html templates
func renderErrorTemplate(w http.ResponseWriter, message string, suggestions []Suggestions) {
	funcMap := template.FuncMap{
		"replaceSpaces": func(s string) string {
			return strings.ReplaceAll(s, " ", "-")
		},
	}
	// Create a new template instance
	t := template.New("error.html").Funcs(funcMap)
	t, err := t.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Error loading error template: %v", err)
		http.Error(w, "Failed to load error template", http.StatusInternalServerError)
		return
	}
	data := struct {
		Message     string
		Suggestions []Suggestions
	}{
		Message:     message,
		Suggestions: suggestions,
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Failed to render error template", http.StatusInternalServerError)
	}
}
