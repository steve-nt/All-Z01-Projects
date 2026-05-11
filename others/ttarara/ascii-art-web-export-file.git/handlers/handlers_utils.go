package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func PrepareTemplateData(title, pageClass string) map[string]interface{} { //To give the title to a page and select the class for Index.html to access the right css style
	return map[string]interface{}{
		"Title":     title,
		"PageClass": pageClass,
		"First":     "",
		"Error":     "",
	}
}

func RenderTemplateFiles(w http.ResponseWriter, templateName string, data interface{}) {
	//w.Header().Set("Content-Type", "text/html")
	log.Println("Rendering template:", templateName)

	tmpl, err := template.ParseFiles("templates/"+templateName, "templates/Index.html")
	if err != nil {
		log.Printf("Status 500: Error loading templates: %v\n", err)
		// na mpainei sto handler tou [Status 500] - handleServerError()
		return
	}
	log.Println("Status 200: Loading templates successfully")

	err = tmpl.ExecuteTemplate(w, "Index", data)
	if err != nil {
		log.Printf("Status 500: Error rendering template: %v\n", err)
		// na mpainei sto handler tou [Status 500] - handleServerError()
	}
}

// _________________________________________________________________________________

func inputValidation(input string) (string, int) {
	input = strings.TrimSpace(input)

	if input == "" {
		return "Input cannot be empty or just whitespace.", 400
	}
	if len(input) >= 128 {
		return "Input too long. Maximum allowed length is 128 characters.", 400
	}
	if !asciiValidation(input) {
		return "Input contains invalid characters. Only ASCII characters, tabs, and newlines are allowed.", 400
	}
	return strings.ReplaceAll(input, "\r\n", "\\n"), 200
}

func asciiValidation(input string) bool {
	for _, char := range input {
		if (char < 32 || char > 126) && char != 9 && char != 10 && char != 13 {
			return false
		}
	}
	return true
}

// _________________________________________________________________________________

func asciiArt(argument string, fonts string) (string, int) {
	bannerStyle, err := os.ReadFile("static/banners/" + fonts + ".txt")
	if err != nil {
		return "", 500
	}

	normalizedBanner := normalizeNewlines(string(bannerStyle))
	lines := strings.Split(normalizedBanner, "\n")
	str := ""
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

func normalizeNewlines(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n") // Replace Windows-style line endings (\r\n) with Unix-style (\n)
	content = strings.ReplaceAll(content, "\r", "\n")   // Replace old Mac-style line endings (\r) with Unix-style (\n)
	return content
}
