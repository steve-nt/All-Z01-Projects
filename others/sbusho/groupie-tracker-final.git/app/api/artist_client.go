package api

import (
	"fmt"
	"groupie-tracker-final/app/models"
)

type ArtistClient struct {
	*Client
}

func NewArtistClient() *ArtistClient {
	return &ArtistClient{NewClient()}
}

func (c *ArtistClient) GetAllArtists() (models.Artists, error) {
	var artists models.Artists
	err := c.get("/artists", &artists)
	return artists, err
}

func (c *ArtistClient) GetArtistByID(id int) (*models.Artist, error) {
	var artist models.Artist
	err := c.get(fmt.Sprintf("/artists/%d", id), &artist)
	return &artist, err
}

func (c *ArtistClient) GetArtistRelations(id int) (*models.Relation, error) {
	var relation models.Relation
	err := c.get(fmt.Sprintf("/relation/%d", id), &relation)
	return &relation, err
}
