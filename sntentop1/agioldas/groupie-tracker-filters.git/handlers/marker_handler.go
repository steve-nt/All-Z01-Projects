package api

import (
	"fmt"
	"groupie/geocoding"
	"net/http"
	"strconv"
	"time"
)

// handler for map marker requests, can respond multiple times to an SSE, asynchronously as the markers are fetched for an API
func MarkerHandler(writer http.ResponseWriter, request *http.Request) {
	artistIDstr := request.URL.Query().Get("artistID")

	artistID, err := strconv.Atoi(artistIDstr)

	if err != nil {
		fmt.Println("junk data provided!: ", artistIDstr)
		return
	}

	dateLocations := ArtistRelationMap[artistID].DatesLocations

	channel := make(chan geocoding.Marker)

	goroutineCount := len(dateLocations)

	//sending goroutines that will call fetchCoordinates and then put the result into a channel
	for location := range dateLocations {
		go findMarker(channel, location)
	}

	// Set headers for SSE
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")

	for marker := range channel {
		// do stuff with marker from goroutine
		payload := fmt.Sprintf(
			"data: {\"lat\": \"%s\", \"lon\": \"%s\", \"location\": \"%s\", \"finished\": \"%s\"}\n\n",
			marker.Latitude, marker.Longitude, marker.Location, "false",
		)
		_, err := fmt.Fprint(writer, payload)
		if err != nil {
			fmt.Println("error trying to fprint the markers: ", err)
			return
		}
		writer.(http.Flusher).Flush()
		goroutineCount--
		if goroutineCount == 0 {
			close(channel)
		}
		time.Sleep(time.Millisecond * 300)
	}

	//ONE LAST SEND TO TELL THE JAVASCRIPT THAT ALL MARKERS ARE FINISHED
	payload := fmt.Sprintf(
		"data: {\"lat\": \"%s\", \"lon\": \"%s\", \"location\": \"%s\", \"finished\": \"%s\"}\n\n",
		"", "", "", "true",
	)
	_, err = fmt.Fprint(writer, payload)
	if err != nil {
		fmt.Println("error trying to fprint the markers: ", err)
		return
	}
	writer.(http.Flusher).Flush()

}

// goroutine that will returns the marker, when it's ready
func findMarker(channel chan geocoding.Marker, location string) {
	marker, err := geocoding.FetchCoordinates(location)
	if err != nil {
		fmt.Println(location, ": ", err)
		return
	}
	channel <- marker
}
