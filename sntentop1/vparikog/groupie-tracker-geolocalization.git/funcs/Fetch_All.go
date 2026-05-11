package groupie_tracker_search

import (
	"fmt"
	"log"
	"sync"
)

func FetchAllData() ([]Artist, error) {
	apiData, err := FetchAPI()
	if err != nil {
		log.Printf("Failed to fetch main API: %v", err)
		return nil, fmt.Errorf("failed to fetch main API: %w", err)
	}

	var (
		artists       []Artist
		locationsResp LocationsResponse
		datesResp     DatesResponse
		relationsResp RelationsResponse
		errors        []error

		wg sync.WaitGroup // WaitGroup to manage Go routines
		mu sync.Mutex     // Mutex to prevent race conditions on shared resources
	)

	// Helper function for concurrent data fetching
	fetch := func(url string, target interface{}) {
		defer wg.Done()
		if err := FetchData(url, target); err != nil {
			mu.Lock()
			errors = append(errors, fmt.Errorf("failed to fetch %s: %w", url, err))
			mu.Unlock()
		}
	}

	// Add tasks to the WaitGroup
	wg.Add(4)
	go fetch(apiData.Artists_API, &artists)
	go fetch(apiData.Locations_API, &locationsResp)
	go fetch(apiData.ConcertDates_API, &datesResp)
	go fetch(apiData.Relations_API, &relationsResp)

	// Wait for all fetches to complete
	wg.Wait()

	// If there were any errors, return them
	if len(errors) > 0 {
		log.Printf("Errors occurred while fetching data: %v", errors)
		return nil, fmt.Errorf("data fetch errors: %v", errors)
	}

	// Map the additional data to each artist
	for i, artist := range artists {
		// Assign Locations
		for _, loc := range locationsResp.Index {
			if loc.ID == artist.ID {
				artists[i].Locations = loc
				break
			}
		}

		// Assign Dates
		for _, date := range datesResp.Index {
			if date.ID == artists[i].ID {
				// Directly assign the slice of dates
				artists[i].Dates = date
				break
			}
		}

		// Assign Relations
		for _, rel := range relationsResp.Index {
			if rel.ID == artist.ID {
				artists[i].Relations = rel
				break
			}
		}
	}

	log.Printf("Data successfully fetched and mapped.")
	return artists, nil
}
