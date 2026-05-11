package core

import (
	"bufio"
	"groupie-tracker/bin"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cacheDuration       = time.Hour
	uniqueLocationsFile = "uniquelocations.txt"
	filtersFile         = "filters.txt"
	allLocationsFile    = "alllocations.txt"
	defaultPort         = 8080
	shutdownTimeout     = 5 * time.Second
)

// ShouldRecreateCacheFiles checks if the cache files need to be recreated.
func ShouldRecreateCacheFiles() bool {
	files := []string{uniqueLocationsFile, filtersFile, allLocationsFile}
	for _, file := range files {
		if !FileExists(file) {
			return true
		}
		info, err := os.Stat(file)
		if err != nil {
			log.Printf("Error stating file %s: %v", file, err)
			return true
		}
		if time.Since(info.ModTime()) > cacheDuration {
			log.Printf("File %s is older than %v, recreating...", file, cacheDuration)
			return true
		}
	}
	return false
}

// FileExists checks if a file exists.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// InitializeGlobalVariables initializes global variables if not cached.
func InitializeGlobalVariables() {
	start := time.Now()

	uniqueLocations := readLinesFromFile(uniqueLocationsFile)
	bin.UniqueLocations = uniqueLocations

	readFiltersFromFile(filtersFile)

	allLocations := readAllLocationsFromFile(allLocationsFile)
	bin.AllLocations = allLocations

	log.Printf("Initialization took %v", time.Since(start))
}

// readLinesFromFile reads lines from a file.
func readLinesFromFile(filename string) []string {
	file, err := os.Open(filename)
	handleError(err, "Error opening "+filename)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	handleError(scanner.Err(), "Error reading "+filename)

	return lines
}

// readFiltersFromFile reads filters from a file for filter initialization.
func readFiltersFromFile(filename string) {
	file, err := os.Open(filename)
	handleError(err, "Error opening "+filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			log.Fatalf("Invalid format in %s: %v", filename, line)
		}
		value, err := strconv.Atoi(parts[1])
		handleError(err, "Error parsing integer from "+filename)
		switch parts[0] {
		case "MinCreationYear":
			bin.MinCreationYear = value
		case "MaxCreationYear":
			bin.MaxCreationYear = value
		case "MinAlbumYear":
			bin.MinAlbumYear = value
		case "MaxAlbumYear":
			bin.MaxAlbumYear = value
		default:
			log.Fatalf("Unknown key in %s: %v", filename, parts[0])
		}
	}
	handleError(scanner.Err(), "Error reading "+filename)
}

// readAllLocationsFromFile reads all locations from a file for filter initialization.
func readAllLocationsFromFile(filename string) []bin.Location {
	file, err := os.Open(filename)
	handleError(err, "Error opening "+filename)
	defer file.Close()

	var allLocations []bin.Location
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", Locations: ")
		if len(parts) != 2 {
			log.Fatalf("Invalid format in %s: %v", filename, line)
		}
		idPart := strings.TrimPrefix(parts[0], "Artist ID: ")
		id, err := strconv.Atoi(idPart)
		handleError(err, "Error parsing artist ID from "+filename)
		locations := strings.Split(strings.Trim(parts[1], "[]"), "+")
		allLocations = append(allLocations, bin.Location{Id: id, Locations: locations})
	}
	handleError(scanner.Err(), "Error reading "+filename)

	return allLocations
}

// handleError logs an error message and exits the program.
func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
