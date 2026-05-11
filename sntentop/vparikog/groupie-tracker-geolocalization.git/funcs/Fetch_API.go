package groupie_tracker_search

import (
	"encoding/json"
	"io"
	"net/http"
)

func FetchAPI() (All_API, error) {
	url := "https://groupietrackers.herokuapp.com/api"
	resp, err := http.Get(url)
	if err != nil {
		return All_API{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return All_API{}, err
	}

	var api All_API
	if err := json.Unmarshal(body, &api); err != nil {
		return All_API{}, err
	}

	return api, nil
}

// FetchData fetches data from a given URL and unmarshals it into the provided target struct
func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}
