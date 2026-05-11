package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	a "platform.zone01.gr/git/santonop/SampleAsciiWeb/internal/ascii"
)

// PageData struct to pass dynamic data to the template
type PageData struct {
	AsciiArt string
}

// ErrorPageData struct for passing error data to error.html
type ErrorPageData struct {
	HttpStatusCode int
	ErrorMessage   string
}

func main() {
	// Serve static files like CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Route for GET /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			renderErrorPage(w, http.StatusNotFound, "Page not found") //404 page not found
			return
		}

		tmpl, err := template.ParseFiles("templates/ascii.art.html")
		if err != nil {
			renderErrorPage(w, http.StatusNotFound, "Home page template not found") //404 page not found
			return
		}

		data := PageData{}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	})

	// Route for POST /ascii-art
	http.HandleFunc("/ascii-art", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			renderErrorPage(w, http.StatusMethodNotAllowed, "Method not allowed") //415 method not allowed
			return
		}

		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			renderErrorPage(w, http.StatusBadRequest, "Unable to parse form") // 400 Bad request
			return
		}

		textInput := strings.ReplaceAll(r.FormValue("text"), "\r", "\n")
		bannerValue := strings.TrimSpace(r.FormValue("banner"))
		bannerPath := "assets/banners/" + bannerValue + ".txt"

		// Validate text input
		if strings.TrimSpace(textInput) == "" {
			renderErrorPage(w, http.StatusBadRequest, "Text is required") //400 Bad request
			return
		}

		// Validate banner input
		if bannerValue == "" {
			renderErrorPage(w, http.StatusBadRequest, "Banner is required") //400 Bad request
			return
		}

		// Check if banner file exists
		if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
			renderErrorPage(w, http.StatusNotFound, "Banner file not found") //404 not found
			return
		}

		// Generate ASCII art
		asciiArt, err := a.GenerateTextToAscii(textInput, bannerPath)
		if err != nil {
			renderErrorPage(w, http.StatusInternalServerError, "Failed to generate ASCII art: "+err.Error())//500 internal server error
			return
		}

		data := PageData{
			AsciiArt: asciiArt,
		}

		tmpl, err := template.ParseFiles("templates/ascii.art.html")
		if err != nil {
			renderErrorPage(w, http.StatusNotFound, "Template file not found")//404 not found
			return
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	})

	// Start the server
	log.Println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// renderErrorPage renders not_found.html for 404, error.html for other errors
func renderErrorPage(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)

	var tmplPath string
	var data interface{}

	if statusCode == http.StatusNotFound {
		tmplPath = "templates/not_found.html"
		data = nil
	} else {
		tmplPath = "templates/error.html"
		data = ErrorPageData{
			HttpStatusCode: statusCode,
			ErrorMessage:   errorMessage,
		}
	}

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError) //500 Internal Server Error
		return
	}

	tmpl.Execute(w, data)
}
