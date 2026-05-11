package api

import (
	"encoding/json"
	"groupie/utils"
	"net/http"
	"sort"
	"strings"
)

// handler for search suggestions, responds to a get request from javascript when user changes anything in the search input field
func SuggestionsHandler(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query().Get("query")

	suggestions := []string{}

	//go through each artist and keep all the matches
	for _, artist := range Artists {
		res, _ := searchMatch(query, artist)
		suggestions = append(suggestions, res...)
	}

	//put all suggestions into a map to remove any duplicates
	myMap := make(map[string]int)
	for _, location := range suggestions {
		if location[0] >= 'a' && location[0] <= 'z' {
			parts := strings.Split(location, " - ")
			myMap[utils.FixKey(parts[0])+" - "+parts[1]] = 1
		} else {
			myMap[location] = 1
		}
	}

	//put em back to the slice
	suggestions = []string{}
	for key := range myMap {
		suggestions = append(suggestions, key)
	}

	//sort them based on their type and then alphanumerically
	sort.Slice(suggestions, func(i, j int) bool {
		iVal := 10
		switch {
		case strings.HasSuffix(suggestions[i], " - artist/band"):
			iVal = 1
		case strings.HasSuffix(suggestions[i], " - member"):
			iVal = 2
		case strings.HasSuffix(suggestions[i], " - creation date"):
			iVal = 3
		case strings.HasSuffix(suggestions[i], " - first album"):
			iVal = 4
		case strings.HasSuffix(suggestions[i], " - concert location"):
			iVal = 5
		}

		jVal := 10
		switch {
		case strings.HasSuffix(suggestions[j], " - artist/band"):
			jVal = 1
		case strings.HasSuffix(suggestions[j], " - member"):
			jVal = 2
		case strings.HasSuffix(suggestions[j], " - creation date"):
			jVal = 3
		case strings.HasSuffix(suggestions[j], " - first album"):
			jVal = 4
		case strings.HasSuffix(suggestions[j], " - concert location"):
			jVal = 5
		}

		//if same category, sort by name
		if iVal == jVal {
			return suggestions[i] < suggestions[j]
		}

		return iVal < jVal
	})

	jsonData, err := json.Marshal(suggestions)
	if err != nil {
		http.Error(writer, "Failed to generate suggestions", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	writer.Write(jsonData)
}
