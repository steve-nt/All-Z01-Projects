package api

import (
	"fmt"
	"groupie-tracker-geolocalization/app/models"
)

type LocationClient struct {
	*Client
}

// NewLocationClient initializes and returns a new client for location-related API calls
func NewLocationClient() *LocationClient {
	return &LocationClient{NewClient()}
}

// GetAllLocations retrieves all available locations from the API
func (c *LocationClient) GetAllLocations() (*models.Locations, error) {
	var locations models.Locations
	err := c.get("/locations", &locations)
	return &locations, err
}

// GetLocationByID fetches details of a specific location by its ID
func (c *LocationClient) GetLocationByID(id int) (*models.Locations, error) {
	var location models.Locations
	err := c.get(fmt.Sprintf("/locations/%d", id), &location)
	return &location, err
}
