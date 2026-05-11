package handlers

import "groupie-tracker-geolocalization/app/models"

var (
	artistsData models.Artists
)

// SetArtistsData sets the artists data for the handlers package
func SetArtistsData(artists models.Artists) {
	artistsData = artists
}

// GetArtistsData gets the artists data
func GetArtistsData() models.Artists {
	return artistsData
}
