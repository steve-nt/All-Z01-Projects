package bin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// fetchAndDecode fetches and decodes data from the API
func fetchAndDecode[T any](url string, id int) (T, error) {
	url += fmt.Sprintf("/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	var data T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		var zeroValue T
		return zeroValue, fmt.Errorf("error decoding JSON response: %w", err)
	}
	return data, nil
}

// fetchAndDecodeSlice fetches and decodes a slice of data from the API
func fetchAndDecodeSlice[T any](url string) ([]T, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	var data []T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}
	return data, nil
}

// FetchArtists fetches artists data from the API
func FetchArtists() ([]Artist, error) {
	return fetchAndDecodeSlice[Artist]("https://groupietrackers.herokuapp.com/api/artists")
}

// FetchLocations fetches locations data from the API
func FetchLocations(id int) (Location, error) {
	return fetchAndDecode[Location]("https://groupietrackers.herokuapp.com/api/locations", id)
}

// FetchAndCacheAllLocations fetches all locations data from the API or cache
func FetchAndCacheAllLocations() ([]Location, error) {
	cacheKey := "all_locations"
	if data, found := locationsCache.Get(cacheKey, 2*time.Minute); found {
		// Serve from cache
		var locations []Location
		json.Unmarshal(data.([]byte), &locations)
		return locations, nil
	}

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	var locationsResponse LocationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&locationsResponse); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %w", err)
	}

	// Cache the fetched locations data
	responseData, _ := json.Marshal(locationsResponse.Index)
	locationsCache.Set(cacheKey, responseData, 2*time.Minute)

	return locationsResponse.Index, nil
}

// FetchDates fetches dates data from the API
func FetchDates(id int) (Date, error) {
	return fetchAndDecode[Date]("https://groupietrackers.herokuapp.com/api/dates", id)
}

// FetchRelations fetches relations data from the API
func FetchRelations(id int) (Relation, error) {
	return fetchAndDecode[Relation]("https://groupietrackers.herokuapp.com/api/relation", id)
}

// FetchAndCacheArtists fetches artists data from the cache or fetches it from the API
var FetchAndCacheArtists = func() ([]Artist, error) {
	var artists []Artist
	if data, found := artistsCache.Get("artists", 2*time.Minute); found {
		// Serve from artists cache
		json.Unmarshal(data.([]byte), &artists)
	} else {
		// Fetch artists data
		var err error
		artists, err = FetchArtists()
		if err != nil {
			return nil, err
		}
		// Cache the fetched artists data
		responseData, _ := json.Marshal(artists)
		artistsCache.Set("artists", responseData, 2*time.Minute)
	}
	return artists, nil
}

// fetchArtistByName fetches an artist by name from the cache
func fetchArtistByName(name string) (Artist, []Suggestions, error) {
	artists, err := FetchAndCacheArtists()
	if err != nil {
		return Artist{}, nil, err
	}
	for _, artist := range artists {
		if strings.EqualFold(artist.Name, name) {
			return artist, nil, nil
		}
	}
	suggestions := generateSuggestions()
	return Artist{}, suggestions, fmt.Errorf("artist not found")
}
