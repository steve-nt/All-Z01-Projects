package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"html/template"
	"ascii-art-web/utils"
)

var BANNERS = []string{"standard", "shadow", "thinkertoy"}

type PageData struct {
	Banners        []string
	SelectedBanner string
	UserInput      string
	AsciiOutput    string
}

func main() {
	flag := false
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if flag {
			utils.RenderErrorPage(w, "Failed to render template")
			return
		}
		var context PageData
		if r.Method == "GET" {
			placeholderText := "Please enter a sample text"
			asciiPlaceholder, err := utils.Ascii_art(placeholderText, "standard")
			if err != nil {
				asciiPlaceholder = err.Error()
			}

			context = PageData{
				Banners:        BANNERS,
				UserInput:      placeholderText,
				AsciiOutput:    asciiPlaceholder,
				SelectedBanner: "standard",
			}
		} else if r.Method == "POST" {
			err := r.ParseForm()
			fmt.Println(r.Form)
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				log.Println("Error parsing form:", err)
				return
			}

			selectedBanner := r.FormValue("banner-select")
			userInput := r.FormValue("user-input")
			fmt.Println(selectedBanner)

			asciiOutput, err := utils.Ascii_art(userInput, selectedBanner)
			if err != nil {
				asciiOutput = err.Error()
			}

			context = PageData{
				Banners:        BANNERS,
				UserInput:      userInput,
				AsciiOutput:    asciiOutput,
				SelectedBanner: selectedBanner,
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl, err := template.ParseFiles(filepath.Join("templates", "home.html"))
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			log.Println("Error loading template:", err)
			return
		}

		err = tmpl.Execute(w, context)
		if err != nil {
			flag = true
			log.Println("Error rendering template:", err)
		}
	})

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
