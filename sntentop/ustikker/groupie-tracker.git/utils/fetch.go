package utils

import (
	"encoding/json"
	"net/http"
)

// FetchData fetches JSON data from a URL and unmarshals it into the target interface.
func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the JSON response directly into the target
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}
