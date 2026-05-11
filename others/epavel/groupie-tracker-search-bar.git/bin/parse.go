package bin

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

// these functions are used to parse the data from the API
// and make it more human-readable etc.

// Sort the concerts of a given town by date
func (c *Concerts) Sort() {
	for country, countryData := range c.Countries {
		for town, townData := range countryData.Towns {
			sort.Slice(townData.Dates, func(i, j int) bool {
				date1, _ := time.Parse("02-01-2006", townData.Dates[i])
				date2, _ := time.Parse("02-01-2006", townData.Dates[j])
				return date1.Before(date2)
			})
			c.Countries[country].Towns[town] = townData
		}
	}
}

// Parse the Relations struct to a custom Concerts struct(deep nested array)
func (r *Relation) Parse() Concerts {
	result := Concerts{
		Id:        r.Id,
		Countries: make(map[string]Country),
	}
	for key, value := range r.Relations { // key is the location, value is the date
		loc := key
		index := strings.Index(key, "-") // left side is the town, right side is the country
		loc = strings.ReplaceAll(loc, "_", " ")
		country := strings.Title(loc[index+1:])
		town := strings.Title(loc[:index])
		if country == "Usa" {
			country = "USA"
		}
		if country == "Uk" {
			country = "UK"
		}
		dates := value
		if _, ok := result.Countries[country]; !ok { // if the country doesn't exist in the map
			result.Countries[country] = Country{ // create a new country
				Towns: make(map[string]Town), // with an empty map of towns
			}
		}
		if _, ok := result.Countries[country].Towns[town]; !ok { // if the town doesn't exist in the map
			result.Countries[country].Towns[town] = Town{ // create a new town
				Dates: make([]string, 0), // with an empty slice of dates
			}
		}
		result.Countries[country].Towns[town] = Town{
			Dates: append(result.Countries[country].Towns[town].Dates, dates...), // append the dates to the town
		}
	}
	result.Sort()
	return result
}

// Parse the locations and manipulate them to be more human-readable
func (l *Location) Parse() {
	newLocations := make([]string, 0)
	for _, location := range l.Locations {
		newLocation := strings.ReplaceAll(location, "-", ", ")
		newLocation = strings.ReplaceAll(newLocation, "_", " ")
		index := strings.Index(newLocation, "usa")
		if index == -1 {
			index = strings.Index(newLocation, "uk")
		}
		if index != -1 {
			newLocation = strings.Title(newLocation[:index]) + strings.ToUpper(newLocation[index:])
		} else {
			newLocation = strings.Title(newLocation)
		}
		newLocations = append(newLocations, newLocation)
	}
	l.Locations = newLocations
}

// Parse the Dates and remove the asterisks
func (d *Date) Parse() {
	newDates := make([]string, 0)
	for _, date := range d.Dates {
		newDate := strings.ReplaceAll(date, "*", "")
		newDates = append(newDates, newDate)
	}
	d.Dates = newDates
}

// shuffleArtists shuffles the list of artists for the home page or suggestions in error page
func shuffleArtists(artists []Artist) {
	// Shuffle the list of artists
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(artists), func(i, j int) { artists[i], artists[j] = artists[j], artists[i] })
}

func sortArtists(artists []Artist) {
	sort.Slice(artists, func(i, j int) bool {
		return artists[i].Name < artists[j].Name
	})
}

// TransformLocationsToFiltered transforms the list of locations to a list of FilteredLocation for the frontend
func TransformLocationsToFiltered(locations []Location) []FilteredLocation {
	// Map to hold country as key and list of towns
	countryMap := make(map[string]map[string]FilterTown)

	for _, loc := range locations {
		for _, fullLocation := range loc.Locations {
			parts := strings.Split(fullLocation, "-") // Split "town-country" format
			if len(parts) != 2 {
				continue // Skip invalid entries
			}

			town := strings.Title(strings.ReplaceAll(parts[0], "_", " "))    // Format town name
			country := strings.Title(strings.ReplaceAll(parts[1], "_", " ")) // Format country name

			// Initialize country in the map if it doesn't exist
			if _, exists := countryMap[country]; !exists {
				countryMap[country] = make(map[string]FilterTown)
			}

			// Add town only if it doesn't already exist
			if _, exists := countryMap[country][town]; !exists {
				countryMap[country][town] = FilterTown{
					Town: town,
				}
			}
		}
	}

	// Convert map to slice of FilteredLocation
	var result []FilteredLocation
	for country, townsMap := range countryMap {
		// Convert town map to slice
		var towns []FilterTown
		for _, town := range townsMap {
			towns = append(towns, town)
		}

		result = append(result, FilteredLocation{
			Country: country,
			Towns:   towns,
		})
	}

	return result
}
