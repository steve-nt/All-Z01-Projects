package handlers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/data"
	"groupie-tracker/server/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// LocationData represents the location data for JavaScript
type LocationData struct {
	Original string   `json:"original"`
	Dates    []string `json:"dates"`
}

// FormatLocation formats a location string by replacing underscores and hyphens
func FormatLocation(location string) string {
	// Replace underscores with spaces
	location = strings.ReplaceAll(location, "_", " ")
	// Replace hyphens with commas and spaces
	location = strings.ReplaceAll(location, "-", ", ")
	// Capitalize first letter of each word
	words := strings.Fields(location)
	for i, word := range words {
		words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		if words[i] == "Usa" {
			words[i] = "USA"
		} else if words[i] == "Uk" {
			words[i] = "UK"
		}
	}
	return strings.Join(words, " ")
}

func FormatDate(date string) string {
	// Remove asterisks (*)
	date = strings.ReplaceAll(date, "*", "")
	// Return the date as is without any other modifications
	return date
}

// ToJSON converts data to a JSON string
func ToJSON(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return "{}"
	}
	return string(jsonBytes)
}

// Create template functions map
var funcMap = template.FuncMap{
	"FormatLocation": FormatLocation,
	"FormatDate":     FormatDate,
	"toJSON":         ToJSON,
}

// Load the templates with function map
var templates = template.Must(template.New("").Funcs(funcMap).ParseFiles(
	filepath.Join("templates", "index.html"),
	filepath.Join("templates", "geo.html"),
))

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("404 Not Found - Path: %s", r.URL.Path)

	tmpl, err := template.ParseFiles("templates/404.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "404 - Page not found", http.StatusNotFound)
		return
	}
}

// IndexHandler renders the homepage with artist data
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexes, err := data.FetchAllData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}

	if len(indexes.Artists) == 0 {
		log.Print("No artist data available")
		http.Error(w, "No artist data available", http.StatusNotFound)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", indexes)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered index page - Status: %d", http.StatusOK)
}

func LandingHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFoundHandler(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Failed to load landing page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render landing page", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered landing page - Status: %d", http.StatusOK)
}

// GeoHandler handles the geo page for each artist
func GeoHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("artistID")
	if artistID == "" {
		log.Print("No artist ID provided")
		http.Error(w, "Artist ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(artistID)
	if err != nil {
		log.Printf("Invalid artist ID format: %s", artistID)
		http.Error(w, "Invalid artist ID format", http.StatusBadRequest)
		return
	}

	indexes, err := data.FetchAllData()
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		return
	}

	var artist models.Artist
	var artistLocations []string
	var relations models.Relation
	artistFound := false

	// Find artist and their locations
	for _, a := range indexes.Artists {
		if a.ID == id {
			artist = a
			artistFound = true
			break
		}
	}

	if !artistFound {
		log.Printf("Artist not found with ID: %s", artistID)
		http.Error(w, fmt.Sprintf("Artist with ID %s not found", artistID), http.StatusNotFound)
		return
	}

	// Find locations
	for _, loc := range indexes.Locations {
		if loc.ID == id {
			artistLocations = loc.Locations
			break
		}
	}

	// Find relations (dates for each location)
	for _, rel := range indexes.Relations {
		if rel.ID == id {
			relations = rel
			break
		}
	}

	// Create location data map for JavaScript
	locationDataMap := make(map[string]LocationData)
	for _, loc := range artistLocations {
		formatted := FormatLocation(loc)
		dates := relations.DatesLocations[loc]
		locationDataMap[formatted] = LocationData{
			Original: loc,
			Dates:    dates,
		}
	}

	// Prepare the data to pass to the template
	data := struct {
		Artist       models.Artist
		Locations    []string
		LocationData map[string]LocationData
		Relations    models.Relation
	}{
		Artist:       artist,
		Locations:    artistLocations,
		LocationData: locationDataMap,
		Relations:    relations,
	}

	err = templates.ExecuteTemplate(w, "geo.html", data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully rendered geo page for artist %s - Status: %d", artistID, http.StatusOK)
}
