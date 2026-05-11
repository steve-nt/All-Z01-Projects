package bin

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// Initialize filters with default values extracted from the dataset
func (f *Filters) Initialize() {
	f.CreationDate.Min = MinCreationYear
	f.CreationDate.Max = MaxCreationYear
	f.AlbumDate.Min = MinAlbumYear
	f.AlbumDate.Max = MaxAlbumYear
	f.Members = make(map[int]bool, 8)
	for i := 1; i <= 8; i++ {
		f.Members[i] = false
	}
}

// Validate filters to ensure that the values are within the acceptable range
func (f *Filters) Validate() error {
	if f.CreationDate.Min != 0 && f.CreationDate.Max != 0 {
		if f.CreationDate.Min > f.CreationDate.Max || f.CreationDate.Min < MinCreationYear || f.CreationDate.Max > MaxCreationYear {
			return ErrInvalidCreationDate
		}
	}
	if f.AlbumDate.Min != 0 && f.AlbumDate.Max != 0 {
		if f.AlbumDate.Min > f.AlbumDate.Max || f.AlbumDate.Min < MinAlbumYear || f.AlbumDate.Max > MaxAlbumYear {
			return ErrInvalidAlbumDate
		}
	}
	if f.Locations == nil {
		return nil
	}
	knownLocations := useUniqueLocationsMap()
	for location := range f.Locations {
		if _, exists := knownLocations[location]; !exists {
			return fmt.Errorf("invalid location: %s", location)
		}
	}
	return nil
}

// Parse filters from query parameters
func parseFilters(query url.Values) (Filters, error) {
	// Parse filters from query parameters
	var filters Filters
	filters.Initialize()
	if creationDateMin := query.Get("creationDateMin"); creationDateMin != "" {
		filters.CreationDate.Min, _ = strconv.Atoi(creationDateMin)
	}
	if creationDateMax := query.Get("creationDateMax"); creationDateMax != "" {
		filters.CreationDate.Max, _ = strconv.Atoi(creationDateMax)
	}
	if albumDateMin := query.Get("albumDateMin"); albumDateMin != "" {
		filters.AlbumDate.Min, _ = strconv.Atoi(albumDateMin)
	}
	if albumDateMax := query.Get("albumDateMax"); albumDateMax != "" {
		filters.AlbumDate.Max, _ = strconv.Atoi(albumDateMax)
	}
	if members := query.Get("members"); members != "" {
		memberIDs := strings.Split(members, "+")
		for _, memberID := range memberIDs {
			if id, err := strconv.Atoi(memberID); err == nil && id >= 1 && id <= 8 {
				filters.Members[id] = true
			}
		}
	} else {
		filters.Members = make(map[int]bool)
	}
	filters.Locations = make(map[string]map[string]bool)
	if locations := strings.ToLower(query.Get("locations")); locations != "" {
		cities := strings.ToLower(query.Get("cities"))
		locationList := strings.Split(locations, "+")
		uniqueLocationsMap := useUniqueLocationsMap()

		// Create a map to store cities grouped by their country
		cityCountryMap := make(map[string][]string)
		if cities != "" {
			cityList := strings.Split(cities, "+")
			for _, city := range cityList {
				for country, cityMap := range uniqueLocationsMap {
					if _, exists := cityMap[city]; exists {
						cityCountryMap[country] = append(cityCountryMap[country], city)
					}
				}
			}
		}

		// Set boolean values for each country and its cities
		for _, location := range locationList {
			location = strings.ToLower(location)
			if cityMap, exists := uniqueLocationsMap[location]; exists {
				if selectedCities, hasCities := cityCountryMap[location]; hasCities {
					for city := range cityMap {
						if contains(selectedCities, city) {
							uniqueLocationsMap[location][city] = true
						} else {
							uniqueLocationsMap[location][city] = false
						}
					}
				} else {
					for city := range cityMap {
						uniqueLocationsMap[location][city] = true
					}
				}
			}
		}
		filters.Locations = uniqueLocationsMap
	} else {
		filters.Locations = make(map[string]map[string]bool)
	}
	if err := filters.Validate(); err != nil {
		log.Printf("Invalid filters: %v\n", err)
		filters.Initialize()
		return filters, err
	}
	return filters, nil
}

// Apply filters to the list of artists
func applyFilters(artists []Artist, filters Filters) []Artist {
	// Apply filters to the list of artists
	var filteredArtists []Artist
	for _, artist := range artists {
		if filterArtist(artist, filters) {
			filteredArtists = append(filteredArtists, artist)
			continue
		}
	}
	return filteredArtists
}

// Filter artist based on the filters
func filterArtist(artist Artist, filters Filters) bool {
	if artist.StartYear < filters.CreationDate.Min || artist.StartYear > filters.CreationDate.Max {
		return false
	}
	album, err := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
	if err != nil {
		return false
	}
	if album < filters.AlbumDate.Min || album > filters.AlbumDate.Max {
		return false
	}
	members := len(artist.Members)
	if len(filters.Members) != 0 {
		if !filters.Members[members] {
			return false
		}
	}
	if len(filters.Locations) != 0 {
		for _, artistLocation := range AllLocations[artist.Id-1].Locations {
			parts := strings.Split(artistLocation, "-")
			if len(parts) == 2 {
				city := parts[0]
				country := parts[1]

				// Check if the country is in the filters
				if cityMap, exists := filters.Locations[country]; exists {
					// If cities are specified for the country, check if the city is true
					if len(cityMap) > 0 {
						if cityMap[city] {
							return true
						}
					} else {
						// If no cities are specified, the whole country is true
						return true
					}
				}
			}
		}
		return false
	}
	return true
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
