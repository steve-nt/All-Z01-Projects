package api

import (
	"fmt"
	"groupie-tracker-final/app/models"
)

type LocationClient struct {
	*Client
}

func NewLocationClient() *LocationClient {
	return &LocationClient{NewClient()}
}

func (c *LocationClient) GetAllLocations() (*models.Locations, error) {
	var locations models.Locations
	err := c.get("/locations", &locations)
	return &locations, err
}

func (c *LocationClient) GetLocationByID(id int) (*models.Locations, error) {
	var location models.Locations
	err := c.get(fmt.Sprintf("/locations/%d", id), &location)
	return &location, err
}
