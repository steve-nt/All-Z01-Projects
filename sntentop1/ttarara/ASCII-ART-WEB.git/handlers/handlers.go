package handlers

import (
	"log"
	"net/http"
)

func PgHome(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("HOME", "home-page")
	log.Println("requesting_page: HOME") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgHome.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgProject(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("PROJECT", "project-page")
	log.Println("requesting_page: PROJECT") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgProject.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgTeam(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("TEAM", "team-page")
	log.Println("requesting_page: TEAM") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgTeam.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
}

func PgConverter(w http.ResponseWriter, r *http.Request) {
	data := PrepareTemplateData("CONVERTER", "converter-page")
	if r.Method == http.MethodPost {

		text := r.FormValue("ascii-data") // Handle form submission
		banner := r.FormValue("banner")   // Handle form submission

		final, renderStatus := inputValidation(text) // Validate and process input
		log.Println("Render Status:", renderStatus, "Rendered Text:", final)
		if renderStatus == 400 {
			log.Println("400 Error", "Invalid Input: "+text)
			handleBadRequest(w, "Invalid Input: "+text)
			return
		}

		if renderStatus != 200 {
			data["Error"] = final
		} else {
			output, statusCode := asciiArt(final, banner)
			if statusCode == 500 {
				log.Println("500 Error", "Failed to read font file for font: "+banner)
				handleServerError(w, "Failed to generate ASCII art")
				return
			}

			if statusCode != 200 {
				data["Error"] = "Failed to generate ASCII art. Please check your input and try again."
			} else {
				data["First"] = output // Set the ASCII Art result
			}
		}
	}
	log.Println("requesting-page: CONVERTER") // na ginei me LogEventsRecord
	RenderTemplateFiles(w, "PgConverter.html", data)
	log.Printf("Rendering template %v successfully\n", data["Title"])
	log.Printf("The output: %v\n", data["First"])
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := PrepareTemplateData("ERROR 404", "code-404")
	RenderTemplateFiles(w, "code_404.html", data)
}

func handleBadRequest(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	data := PrepareTemplateData("ERROR 400", "code-400")
	data["Error"] = errorMessage
	RenderTemplateFiles(w, "code_400.html", data)
}

func handleServerError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	data := PrepareTemplateData("ERROR 500", "code-500")
	data["Error"] = errorMessage
	RenderTemplateFiles(w, "code_500.html", data)
}
