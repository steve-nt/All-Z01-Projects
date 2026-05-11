package groupie_tracker_search

import (
	"sync"
	"time"
)

var (
	cachedArtists []Artist
	cacheMu       sync.RWMutex       // Mutex to manage concurrent access to the cache
	cacheTime     time.Time          // Timestamp of the last cache update
	cacheDuration = 30 * time.Second // Duration before refreshing the cache
)

// Fetch and cache data if necessary
func getCachedArtists() ([]Artist, error) {
	cacheMu.RLock()
	// Return cached data if it's still fresh
	if time.Since(cacheTime) < cacheDuration {
		cachedData := cachedArtists
		cacheMu.RUnlock()
		return cachedData, nil
	}
	cacheMu.RUnlock()

	// Acquire write lock to refresh cache
	cacheMu.Lock()
	defer cacheMu.Unlock()

	// Double-check cache expiration in case of a race condition
	if time.Since(cacheTime) < cacheDuration {
		return cachedArtists, nil
	}

	// Fetch fresh data
	artists, err := FetchAllData()
	if err != nil {
		return nil, err
	}

	// Update the cache
	cachedArtists = artists
	cacheTime = time.Now()

	// //Filter location
	// mapping := GetLocationMapping(artists)

	// // For a consistent output, sort the country keys.
	// var countries []string
	// for country := range mapping {
	// 	countries = append(countries, country)
	// }
	// sort.Strings(countries)

	// // Print the mapping.
	// fmt.Println("Mapping of countries to unique cities:")
	// for _, country := range countries {
	// 	fmt.Printf("%s: %v\n", country, mapping[country])
	// }

	return cachedArtists, nil
}
