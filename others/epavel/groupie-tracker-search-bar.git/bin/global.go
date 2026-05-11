package bin

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// this file contains global variables and types that are used across multiple files
// it also contains functions that are used to populate these global variables

type Artist struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	StartYear  int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Members    []string `json:"members"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

type LocationsResponse struct {
	Index []Location `json:"index"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

type FilterTown struct {
	Town string `json:"town"`
}

type FilteredLocation struct {
	Country string       `json:"country"`
	Towns   []FilterTown `json:"towns"`
}

type Town struct {
	Dates []string
}

type Country struct {
	Towns map[string]Town
}

type Concerts struct {
	Id        int
	Countries map[string]Country
}

type Suggestions struct {
	Name  string
	Image string
}

type ArtistPageData struct {
	Artist       Artist
	Dates        Date
	Locations    Location
	Relations    Concerts
	Maps         string
	ConcertsJSON string
}

type HomePageData struct {
	Artists        []Artist
	Message        string
	Total          int
	Current        int
	NextPagination int
	Shuffle        string
	Filter         Filters
	HasFilters     bool
	HasSearchQuery bool
}

type Filters struct {
	CreationDate FilterDate
	AlbumDate    FilterDate
	Members      map[int]bool
	Locations    map[string]map[string]bool
}

type FilterDate struct {
	Min int
	Max int
}

var (
	UniqueLocations        []string
	MinCreationYear        int
	MaxCreationYear        int
	MinAlbumYear           int
	MaxAlbumYear           int
	AllLocations           []Location
	ErrInvalidAlbumDate    = errors.New("invalid album date")
	ErrInvalidCreationDate = errors.New("invalid creation date")
	ErrInvalidFilter       = errors.New("invalid filter")
	mu                     sync.Mutex     // Mutex to protect the global variables
	wg                     sync.WaitGroup // WaitGroup to wait for all goroutines to finish
)

// generates suggestions for the error page
func generateSuggestions() []Suggestions {
	var suggestions []Suggestions
	artists, err := FetchAndCacheArtists()
	if err != nil {
		log.Fatalf("Error fetching artists: %v", err)
	}
	// Shuffle the list of artists
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(artists), func(i, j int) { artists[i], artists[j] = artists[j], artists[i] })

	// Pick the first 3 artists from the shuffled list
	for i := 0; i < 3 && i < len(artists); i++ {
		suggestions = append(suggestions, Suggestions{Name: artists[i].Name, Image: artists[i].Image})
	}

	return suggestions
}

// fetches and caches the artists data, crucial for filter initialization and artist search
func PopulateUniqueLocations() {
	wg.Add(1) // Increment the WaitGroup counter

	artists, err := FetchAndCacheArtists()
	if err != nil {
		log.Fatalf("Error fetching artists: %v", err)
	}
	locations, err := FetchAndCacheAllLocations()
	if err != nil {
		log.Fatalf("Error fetching locations: %v", err)
	}
	AllLocations = locations
	MinCreationYear, MaxCreationYear = artists[0].StartYear, artists[0].StartYear
	MinAlbumYear, _ = strconv.Atoi(strings.Split(artists[0].FirstAlbum, "-")[2])
	MaxAlbumYear = MinAlbumYear
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine completes

		locationSet := make(map[string]struct{})
		for _, artist := range artists {
			album, _ := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
			if artist.StartYear < MinCreationYear {
				MinCreationYear = artist.StartYear
			}
			if artist.StartYear > MaxCreationYear {
				MaxCreationYear = artist.StartYear
			}
			if album < MinAlbumYear {
				MinAlbumYear = album
			}
			if album > MaxAlbumYear {
				MaxAlbumYear = album
			}
			for _, location := range AllLocations[artist.Id-1].Locations {
				locationSet[location] = struct{}{}
			}
		}

		mu.Lock()
		for location := range locationSet {
			UniqueLocations = append(UniqueLocations, location)
		}
		// Write unique locations to file
		writeLinesToFile("uniquelocations.txt", UniqueLocations)

		// Write filter data to file
		writeFiltersToFile("filters.txt")

		// Write all locations to file
		writeAllLocationsToFile("alllocations.txt", AllLocations)
		mu.Unlock()
	}()
}

// writeLinesToFile writes a slice of lines to a file in the core directory
func writeLinesToFile(filename string, lines []string) {
	filePath := getCoreFilePath(filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	for _, line := range lines {
		file.WriteString(line + "\n")
	}
}

// writeFiltersToFile writes the filter data to a file in the core directory
func writeFiltersToFile(filename string) {
	filePath := getCoreFilePath(filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString("MinCreationYear: " + strconv.Itoa(MinCreationYear) + "\n")
	file.WriteString("MaxCreationYear: " + strconv.Itoa(MaxCreationYear) + "\n")
	file.WriteString("MinAlbumYear: " + strconv.Itoa(MinAlbumYear) + "\n")
	file.WriteString("MaxAlbumYear: " + strconv.Itoa(MaxAlbumYear) + "\n")
}

// writeAllLocationsToFile writes all locations data to a file in the core directory
func writeAllLocationsToFile(filename string, locations []Location) {
	filePath := getCoreFilePath(filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	for _, location := range locations {
		locationsStr := strings.Join(location.Locations, "+")
		file.WriteString(fmt.Sprintf("Artist ID: %d, Locations: %s\n", location.Id, locationsStr))
	}
}

func GetLocationsFromFile() []Location {
	filePath := getCoreFilePath("alllocations.txt")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var locations []Location
	for {
		var id int
		var locationsStr string
		_, err := fmt.Fscanf(file, "Artist ID: %d, Locations: %s\n", &id, &locationsStr)
		if err != nil {
			break
		}
		locations = append(locations, Location{Id: id, Locations: strings.Split(locationsStr, "&")})
	}
	return locations
}

// getCoreFilePath constructs the absolute path to a file in the core directory
func getCoreFilePath(filename string) string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	// Construct the absolute path to the core directory
	return filepath.Join(cwd, "core", filename)
}

// fetches and caches a map of unique locations, crucial for filter initialization
func useUniqueLocationsMap() map[string]map[string]bool {

	wg.Wait() // Wait for the WaitGroup to complete

	mu.Lock()
	defer mu.Unlock()

	UniqueLocationsMap := make(map[string]map[string]bool)
	for _, location := range UniqueLocations {
		parts := strings.Split(location, "-")
		if len(parts) == 2 {
			city := parts[0]
			country := parts[1]
			if _, exists := UniqueLocationsMap[country]; !exists {
				UniqueLocationsMap[country] = make(map[string]bool)
			}
			UniqueLocationsMap[country][city] = false
		}
	}

	return UniqueLocationsMap
}
