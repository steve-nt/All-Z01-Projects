package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func ErrorPage(w http.ResponseWriter, statusCode int) {
	templateFile := fmt.Sprintf("static/status/%d.html", statusCode)
	tmpl, err := template.ParseFiles(templateFile, "static/header.html", "static/footer.html")
	if err != nil {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	PageData := PageData{
		YearNow:           time.Now().Year(),
		DatesAndLocations: getDatesAndLocations(artists),
	}

	w.WriteHeader(statusCode)
	tmpl.Execute(w, PageData)
}
