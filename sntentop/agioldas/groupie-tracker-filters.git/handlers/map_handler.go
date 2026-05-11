package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// Handler of map page
func MapHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/map" {
		SendErrorPage(writer, 404, "404 - Page not found")
		return
	}

	artistIDStr := request.URL.Query().Get("artistID")
	artistIDint, err := strconv.Atoi(artistIDStr)
	if err != nil {
		SendErrorPage(writer, 404, "404 - Page not found")
		return
	}

	tmpl, ok := template.ParseFiles("../templates/map.html")
	if ok != nil {
		fmt.Println("ERROR:", ok)
		SendErrorPage(writer, 500, "500 - Internal Server Error")
		return
	}

	//will be sending artistID to map page, so that map page can initiate an SSE request for all the map markers
	data := struct {
		ID      string
		Display string
	}{
		ID:      artistIDStr,
		Display: fmt.Sprintf("%s's Concerts", ArtistMap[artistIDint].Name),
	}

	if err := tmpl.Execute(writer, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
