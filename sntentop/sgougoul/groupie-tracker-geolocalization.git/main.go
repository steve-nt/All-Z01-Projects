package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"sgougoupractice/fetch"
	nominatim "sgougoupractice/geocoding/client"
	fadapter "sgougoupractice/geocoding/fetch"
	geocode "sgougoupractice/geocoding/service"
	"sgougoupractice/handlers"
	"sgougoupractice/routes"
	"sgougoupractice/suggestions"
)

var MapQuestKey string

func init() {
	if MapQuestKey == "" {
		log.Fatal("MAPQUEST_KEY not set (build with -ldflags)")
	}
}

// Main function to start the server
func main() {

	suggestions.InitCache()
	go suggestions.RefreshCache(5 * time.Minute)

	mqClient, err := nominatim.NewMQ(MapQuestKey, "GTRACKER", nil)

	if err != nil {
		log.Fatalf("invalid  client: %v", err)
	}

	geoSvc := geocode.NewService(mqClient)
	fetcher := fadapter.FetcherFunc(func(ctx context.Context, id int) (*fadapter.RawLocation, error) {
		loc, err := fetch.FetchLocations(id)
		if err != nil {
			return nil, err
		}
		return &fadapter.RawLocation{
			ID:        loc.ID,
			Locations: loc.Locations,
		}, nil
	})

	h := routes.NewHandler(fetcher, geoSvc)

	// Middleware to recover from panics
	recoveryMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered from panic: %v", err)
					handlers.ErrorHandler(w, r, http.StatusInternalServerError, "Internal Server Error.")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
	r := routes.RouteHandler(h)
	wrappedRouter := recoveryMiddleware(r)

	log.Println("Starting server on :8080")
	log.Println("Server is running at: http://localhost:8080/")

	// Start the server and handle potential startup errors
	if err := http.ListenAndServe(":8080", wrappedRouter); err != nil {
		log.Fatal(err)
	}
}
