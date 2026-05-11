package server

import (
	"encoding/json"
	"fmt"
	"groupie/fetch"
	"os"
)

var artists []fetch.Artist

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func saveData(filename string, artists []fetch.Artist) error {
	dataBytes, err := json.MarshalIndent(artists, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, dataBytes, 0644)
}

func dataLoad(filename string) ([]fetch.Artist, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func Fetching(apiURL string, filename string) error {
	var err error
	if fileExists(filename) {
		fmt.Println("Using cached data")
		artists, err = dataLoad(filename)
		if err != nil {
			return err
		}
		return nil
	}
	fetch.FetchArtists(apiURL, func(fetchedArtists []fetch.Artist) {
		artists = fetchedArtists
		saveData(filename, artists)
	})
	return nil
}
