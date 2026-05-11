package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"
)

func LoadDataFromURL[T any](URL string, result *T, maxRetries int) (*T, error) {

	var err error

	retries := maxRetries
	for retries > 0 {
		//context with 1 second timeout
		err = attemptRequest(URL, result)
		if err != nil {
			time.Sleep(time.Second * 2)
			retries--
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("failed to fetch link after %d retries: %v", maxRetries, err)
}

func attemptRequest[T any](URL string, result *T) error {
	var response *http.Response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return fmt.Errorf("request creation failed: %w", err)
	}
	response, err = http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}

	return nil
}

func FormatMapKeys(mainmap map[string][]string) map[string][]string {
	// Iterate through the map
	mainmap2 := make(map[string][]string)
	for location := range mainmap {
		mainmap2[FixKey(location)] = mainmap[location]
	}
	return mainmap2
}

func FixKey(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = CleanUpTitle(s)
	s = strings.ReplaceAll(s, "-", " - ")
	s = strings.ReplaceAll(s, "Usa", "USA")
	if strings.HasSuffix(s, "Uk") {
		s = strings.ReplaceAll(s, "Uk", "UK")
	}
	return s
}

func CleanUpTitle(s string) string {
	var nextTitle bool
	sRune := []rune(s)
	for i := 0; i < len(sRune); i++ {
		if i == 0 {
			sRune[i] = sRune[i] - 32
		}
		if sRune[i] == ' ' || sRune[i] == '-' {
			nextTitle = true
		} else if sRune[i] != ' ' && sRune[i] != '-' && nextTitle {
			sRune[i] = sRune[i] - 32
			nextTitle = false
		}
	}
	return string(sRune)
}

// sorts location/dates map by date
func SortDates(datesMap map[string][]string) map[string][]string {
	for key, dateSlice := range datesMap {
		for i := 0; i < len(dateSlice); i++ {
			parts := strings.Split(dateSlice[i], "-")
			if len(parts) == 1 {
				parts = strings.Split(parts[0], "/")
			}
			// fmt.Println(parts)
			dateSlice[i] = parts[2] + "-" + parts[1] + "-" + parts[0]
		}

		slices.Sort(dateSlice)
		// slices.Reverse(dateSlice)

		for i := 0; i < len(dateSlice); i++ {
			parts := strings.Split(dateSlice[i], "-")
			dateSlice[i] = parts[2] + "/" + parts[1] + "/" + parts[0]
		}
		datesMap[key] = dateSlice
	}
	return datesMap
}

// sorts location/dates map by location
func SortLocations(datesLocations map[string][]string) []string {
	locations := make([]string, 0, len(datesLocations))
	for location := range datesLocations {
		locations = append(locations, location)
	}
	sort.Slice(locations, func(i int, j int) bool {
		iparts := strings.Split(locations[i], " - ")
		jparts := strings.Split(locations[j], " - ")
		if iparts[1] == jparts[1] {
			return iparts[0] < jparts[0]
		}
		return iparts[1] < jparts[1]
	})
	return locations
}

// custom hasPrefix
func SameEnough(a, b string) bool {
	return strings.HasPrefix(strings.ToLower(a), strings.ToLower(b))
}

// splits string by non alphanumeric separators
func SplitByWords(str string) []string {
	str = strings.ToLower(str)
	var builder strings.Builder
	for _, r := range str {
		if IsAlphaNumeric(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	parts := strings.Split(builder.String(), " ")
	newSlice := []string{}
	for _, str := range parts {
		if str == "" {
			continue
		}
		newSlice = append(newSlice, str)
	}
	return newSlice
}

func IsAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}
