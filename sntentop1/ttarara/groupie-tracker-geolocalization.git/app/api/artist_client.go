package api

import (
	"fmt"
	"groupie-tracker-geolocalization/app/models"
)

type ArtistClient struct {
	*Client
}

// NewArtistClient initializes and returns a new client for artist-related API calls
func NewArtistClient() *ArtistClient {
	return &ArtistClient{NewClient()}
}

// GetAllArtists retrieves the list of all artists from the API
func (c *ArtistClient) GetAllArtists() (models.Artists, error) {
	var artists models.Artists
	err := c.get("/artists", &artists)
	return artists, err
}

// GetArtistByID fetches details of a specific artist using their ID
func (c *ArtistClient) GetArtistByID(id int) (*models.Artist, error) {
	var artist models.Artist
	err := c.get(fmt.Sprintf("/artists/%d", id), &artist)
	return &artist, err
}

// GetArtistRelations retrieves relationship data for a given artist ID
func (c *ArtistClient) GetArtistRelations(id int) (*models.Relation, error) {
	var relation models.Relation
	err := c.get(fmt.Sprintf("/relation/%d", id), &relation)
	return &relation, err
}
