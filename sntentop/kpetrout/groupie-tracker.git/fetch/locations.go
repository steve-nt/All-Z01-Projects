package fetch

import (
	"fmt"
	"groupie/geo"
	"strings"
	"sync"
)

var LocCache = make(map[string]*Location)
var LocMutex sync.Mutex

type Location struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

func LocationParse(relations map[string][]string) []*Location {
	locs := make([]*Location, 0, len(relations))
	for relation := range relations {
		locs = append(locs, cachedLocations(relation))
	}
	return locs
}

func cachedLocations(name string) *Location {
	if loc, exists := LocCache[name]; exists {
		return loc
	}

	newLoc := &Location{Name: name}
	LocMutex.Lock()
	LocCache[name] = newLoc
	LocMutex.Unlock()
	return newLoc
}

func FetchCoords(locs []*Location) error {
	for _, loc := range locs {
		parts := strings.Split(loc.Name, ", ")
		if len(parts) < 2 {
			continue
		}

		city, country := parts[0], parts[1]

		LocMutex.Lock()
		if cached, ok := LocCache[loc.Name]; ok && cached.Lat != 0 && cached.Lng != 0 {
			fmt.Printf("Using cached data for %s\n", loc.Name)
			loc.Lat, loc.Lng = cached.Lat, cached.Lng
			LocMutex.Unlock()
			continue
		}
		LocMutex.Unlock()

		city, country = apiSpecifics(city, country)

		var lat, lng float64
		switch loc.Name {
		case "georgia, usa":
			lat, lng = 32.1656, -82.9001
		case "massachusetts, usa":
			lat, lng = 42.4072, -71.3824
		case "alabama, usa":
			lat, lng = 32.3182, -86.9023
		default:
			var err error
			lat, lng, err = geo.Geocode(city, country)
			if err != nil {
				fmt.Printf("Geocoding failed for %s: %v\n", loc.Name, err)
				continue
			}
		}

		loc.Lat, loc.Lng = lat, lng

		fmt.Println("Geocoding completed for", loc.Name)
		LocMutex.Lock()
		LocCache[loc.Name] = &Location{Name: loc.Name, Lat: lat, Lng: lng}
		fmt.Printf("Cached data for %s\n", loc.Name)
		LocMutex.Unlock()
	}

	return nil
}

func apiSpecifics(city, country string) (string, string) {
	switch country {
	case "uk":
		country = "gb"
	case "french polynesia":
		country = "fp"
	case "new caledonia":
		country = "fr"
	}
	switch city {
	case "north carolina":
		city = "charlotte"
	case "south carolina":
		city = "columbia"
	case "queensland":
		city = "brisbane"
	}
	return city, country
}
