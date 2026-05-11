package handlers

import (
	"gtracker/internal/models"
	"gtracker/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"
)

// FilterArtistsHandler handles the request to filter artists
func FilterArtistsHandler(w http.ResponseWriter, r *http.Request) {
	var filteredArtists []models.Artist

	// Filter by creation date
	minCreationDateStr := r.URL.Query().Get("minCreationDate")
	maxCreationDateStr := r.URL.Query().Get("maxCreationDate")
	var minCreationDate, maxCreationDate int
	var err error
	if minCreationDateStr != "" {
		minCreationDate, err = strconv.Atoi(minCreationDateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := services.BadRequestTemplate.Execute(w, nil); err != nil {
				http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
				log.Printf("Error executing noresults template: %v", err)
				return
			}
			return
		}
	} else {
		minCreationDate = 0
	}

	if maxCreationDateStr != "" {
		maxCreationDate, err = strconv.Atoi(maxCreationDateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := services.BadRequestTemplate.Execute(w, nil); err != nil {
				http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
				log.Printf("Error executing noresults template: %v", err)
				return
			}
			return
		}
	} else {
		maxCreationDate = time.Now().Year()
	}

	filteredArtists = services.FilterArtistsByCreationDate(minCreationDate, maxCreationDate)

	// Filter by first album date
	minAlbumDateStr := r.URL.Query().Get("minAlbumDate")
	maxAlbumDateStr := r.URL.Query().Get("maxAlbumDate")
	if minAlbumDateStr != "" || maxAlbumDateStr != "" {
		var minAlbumDate, maxAlbumDate int
		if minAlbumDateStr != "" {
			minAlbumDate, err = strconv.Atoi(minAlbumDateStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := services.BadRequestTemplate.Execute(w, nil); err != nil {
					http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
					log.Printf("Error executing noresults template: %v", err)
					return
				}
				return
			}
		} else {
			minAlbumDate = 0
		}
		if maxAlbumDateStr != "" {
			maxAlbumDate, err = strconv.Atoi(maxAlbumDateStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := services.BadRequestTemplate.Execute(w, nil); err != nil {
					http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
					log.Printf("Error executing noresults template: %v", err)
					return
				}
				return
			}
		} else {
			maxAlbumDate = time.Now().Year()
		}
		filteredArtists = services.Intersect(filteredArtists, services.FilterArtistsByFirstAlbumDate(minAlbumDate, maxAlbumDate))
	}

	// Filter by number of members
	membersStr := r.URL.Query()["members"]
	if len(membersStr) > 0 {
		var memberCounts []int
		for _, m := range membersStr {
			if m == "7+" {
				memberCounts = append(memberCounts, 7)
			} else {
				count, err := strconv.Atoi(m)
				if err == nil {
					memberCounts = append(memberCounts, count)
				}
			}
		}
		filteredArtists = services.Intersect(filteredArtists, services.FilterArtistsByMembers(memberCounts))
	}

	// Filter by concert locations
	locations := r.URL.Query()["locations"]
	if len(locations) > 0 {
		filteredArtists = services.Intersect(filteredArtists, services.FilterArtistsByConcertLocations(locations))
	}

	if len(filteredArtists) == 0 {

		if err := services.NoResultsTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render no results page", http.StatusInternalServerError)
			log.Printf("Error executing noresults template: %v", err)
		}
		return
	}

	if err := services.IndexTemplate.Execute(w, filteredArtists); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
