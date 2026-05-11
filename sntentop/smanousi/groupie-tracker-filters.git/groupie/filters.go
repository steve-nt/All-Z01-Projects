package groupie

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var locationMap = map[string][]string{
	// 1) “north_carolina-usa” as state, found “charlotte-usa” as city
	"north_carolina": {
		"charlotte",
	},
	"paris": {
		"boulogne_billancourt",
	},
	"toul": {
		"pagney_derriere_barine",
	},
	"auckland": {
		"penrose",
	},
	"buenos_aires": {
		"san_isidro",
	},
	"southend_on_sea": {
		"westcliff_on_sea",
	},
	"metz": {
		"freyming_merlebach",
	},
	"castellon": {
		"burriana",
	},
	"melbourne": {
		"west_melbourne",
	},

	// 2) “georgia-usa” as state, found “atlanta-usa” as city
	"georgia": {
		"atlanta",
	},

	// 3) “california-usa” as state, matched with “los_angeles-usa”, “anaheim-usa”, “oakland-usa”, “del_mar-usa”, “san_francisco-usa”
	"california": {
		"los_angeles",
		"anaheim",
		"oakland",
		"del_mar",
		"san_francisco",
		"inglewood",
		"pico_rivera",
	},

	"los_angeles": {
		"pico_rivera",
	},

	// 4) “arizona-usa” as state, but no city “phoenix-usa” (or similar) found -> empty
	"arizona": {"phoenix"},

	// 5) “texas-usa” as state, matched “houston-usa” and “dallas-usa”
	"texas": {
		"houston",
		"dallas",
	},

	// 6) “nevada-usa” as state, matched “las_vegas-usa”
	"nevada": {
		"las_vegas",
	},

	// 7) “new_york-usa” as state, matched “brooklyn-usa” and “amityville-usa”
	"new_york": {
		"brooklyn",
		"amityville",
	},

	// 8) “illinois-usa” as state, matched “chicago-usa” and “rosemont-usa”
	"illinois": {
		"chicago",
		"rosemont",
	},

	// 9) “maine-usa” as state, no city (e.g., “portland-usa”) found -> empty
	"maine": {"portland"},

	// 10) “florida-usa” as state, no city (e.g., “miami-usa”) in the data -> empty
	"florida": {"miami"},

	// 11) “south_carolina-usa” as state, matched “columbia-usa”
	"south_carolina": {
		"columbia",
	},

	// 12) “michigan-usa” as state, matched “detroit-usa”, “grand_rapids-usa”
	"michigan": {
		"detroit",
		"grand_rapids",
	},

	// 13) “missouri-usa” as state, matched “st_louis-usa”
	"missouri": {
		"st_louis",
	},

	// 14) “alabama-usa” found, but no city “mobile-usa” or “birmingham-usa” – the data’s “birmingham-uk” is a different country -> empty
	"alabama": {"huntsville"},

	// 15) “massachusetts-usa” matched “boston-usa”
	"massachusetts": {
		"boston",
	},

	// 20) “washington-usa” matched “seattle-usa”
	"washington": {
		"seattle",
	},

	// 21) “kansas_city-usa” is ambiguous – some treat “kansas_city” as a city, so “kansas” => ["kansas_city"] if we forcibly interpret the “region” as “kansas.” But the data only shows “kansas_city-usa” (ID=29 #1), no “kansas-usa,” so we skip it here.
}

func FilterArtistHandler(w http.ResponseWriter, r *http.Request, artists *[]Artist) {
	var filters struct {
		YearOfFormation struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"yearOfFormation"`
		FirstAlbum struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"firstAlbum"`
		Members  []int  `json:"members"`
		Location string `json:"location"`
	}
	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&filters); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var filteredArtists []Artist
	for _, artist := range *artists {
		// Filter by year-of-formation
		if filters.YearOfFormation.Start > 0 || filters.YearOfFormation.End > 0 {
			if filters.YearOfFormation.Start > 0 && artist.CreationDate < filters.YearOfFormation.Start {
				continue
			}
			if filters.YearOfFormation.End > 0 && artist.CreationDate > filters.YearOfFormation.End {
				continue
			}
		}

		// Filter by first album (similar approach)
		if filters.FirstAlbum.Start > 0 || filters.FirstAlbum.End > 0 {
			parts := strings.Split(artist.FirstAlbum, "-")
			if len(parts) < 3 {
				continue
			}
			albumYear, err := strconv.Atoi(parts[2])
			if err != nil {
				continue
			}
			if filters.FirstAlbum.Start > 0 && albumYear < filters.FirstAlbum.Start {
				continue
			}
			if filters.FirstAlbum.End > 0 && albumYear > filters.FirstAlbum.End {
				continue
			}
		}

		// Filter by members and location (unchanged)
		if len(filters.Members) > 0 {
			match := false
			for _, num := range filters.Members {
				if len(artist.Members) == num {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		if filters.Location != "" {
			locationMatch := false
			userFilter := strings.ToLower(filters.Location)
			normalizedFilter := strings.ReplaceAll(userFilter, " ", "_")
			for _, loc := range artist.Locations {
				if strings.Contains(strings.ToLower(loc), strings.ToLower(filters.Location)) {
					locationMatch = true
					break
				}
				// Check if the user's filter is a key in locationMap
				if cities, ok := locationMap[normalizedFilter]; ok {
					// For each artist location string (e.g., "Los Angeles, Usa")
					for _, loc := range artist.Locations {
						locLower := strings.ToLower(loc)
						// For each city expected under the user filter region,
						// convert underscores to spaces for comparison.
						for _, city := range cities {
							cityFormatted := strings.ReplaceAll(city, "_", " ")
							if strings.Contains(locLower, cityFormatted) {
								locationMatch = true
								break
							}
						}
						if locationMatch {
							break
						}
					}
				} else {
					// Fallback: direct substring match on the artist's location
					for _, loc := range artist.Locations {
						if strings.Contains(strings.ToLower(loc), normalizedFilter) {
							locationMatch = true
							break
						}
					}
				}
			}
			if !locationMatch {
				continue
			}
		}

		// Ensure Locations, Dates, and Relation are always included
		artistCopy := artist
		if artistCopy.Locations == nil {
			artistCopy.Locations = []string{}
		}
		if artistCopy.Dates == nil {
			artistCopy.Dates = []string{}
		}
		if artistCopy.Relation == nil {
			artistCopy.Relation = map[string]string{}
		}
		filteredArtists = append(filteredArtists, artistCopy)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredArtists)
}
