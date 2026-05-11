package backend

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

// HandleAllLocations returns a JSON array of *all* unique location strings
func HandleAllLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch the entire "locations" dataset from groupietrackers
	allLocIndex, err := fetchData[LocationsIndex]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, "Failed to fetch location data", http.StatusInternalServerError)
		return
	}

	// Build a set of unique location strings
	uniqueLocs := make(map[string]bool)
	for _, block := range allLocIndex.Index {
		for _, loc := range block.Locations {
			uniqueLocs[loc] = true
		}
	}

	// Convert set to slice
	var allLocs []string
	for locStr := range uniqueLocs {
		allLocs = append(allLocs, locStr)
	}

	// Optionally sort them alphabetically:
	sort.Slice(allLocs, func(i, j int) bool {
		return strings.ToLower(allLocs[i]) < strings.ToLower(allLocs[j])
	})

	// Return as JSON array
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allLocs)
}
