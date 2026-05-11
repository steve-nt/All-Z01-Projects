package groupie_tracker_search

import (
	"fmt"
	"time"
)

func CreationDates(s []Artist) (int, int) {
	max := 0
	min := 2025
	for _, creationDates := range s {
		if creationDates.CreationDate >= max {
			max = creationDates.CreationDate
		}
		if creationDates.CreationDate <= min {
			min = creationDates.CreationDate
		}
	}
	return min, max
}

func FirstAlbum(s []Artist) (int, int) {
	min := 2025
	max := 0
	for _, firstAlbum := range s {
		parseDate, err := time.Parse("02-02-2006", firstAlbum.FirstAlbum)
		if err != nil {
			fmt.Println("Error parsing date:", err)
		}
		year := parseDate.Year()
		if year >= max {
			max = year
		}
		if year <= min {
			min = year
		}
	}
	return min, max
}
