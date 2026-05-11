package api

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func SendErrorPage(writer http.ResponseWriter, errorType int, message string) {
	tmpl, err := template.ParseFiles("../templates/error.html")
	if err != nil {
		fmt.Println("ERROR:", err)
		http.Error(writer, "500 - Internal Server Super Error", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(errorType)

	if err := tmpl.Execute(writer, struct{ Message string }{Message: message}); err != nil {
		log.Printf("Template execution error: %v", err)
	}

}
