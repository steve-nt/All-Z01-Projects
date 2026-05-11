package bin

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContextKey string // Define a custom type for context keys

const (
	DataKey         ContextKey = "artists"
	totalArtistsKey ContextKey = "totalArtists"
	messageKey      ContextKey = "message"
	errorKey        ContextKey = "error"
)

// Middleware function to handle shuffle, pagination, filter, and search queries and using context to pass data to the next handler
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var message string
		pagination, _ := strconv.Atoi(query.Get("pagination"))
		if pagination == 0 {
			pagination = 12
		}
		shuffle := query.Get("shuffle") == "true"
		searchQuery := strings.ReplaceAll(query.Get("query"), "_", " ")
		filters, err := parseFilters(query)
		if err != nil {
			log.Printf("Error parsing filters: %v", err)
			ctx := context.WithValue(r.Context(), errorKey, err)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		cacheKey := generateCacheKey(pagination, shuffle, searchQuery, filters)

		if data, found := queryCache.Get(cacheKey, 1*time.Minute); found {
			// Serve from queryCache
			var cachedData CachedData
			json.Unmarshal(data.([]byte), &cachedData)
			if shuffle {
				shuffleArtists(cachedData.PaginatedArtists)
			}
			ctx := context.WithValue(r.Context(), DataKey, cachedData.PaginatedArtists)
			ctx = context.WithValue(ctx, totalArtistsKey, cachedData.TotalArtists)
			ctx = context.WithValue(ctx, messageKey, cachedData.Message)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// Fetch data and apply filters
		artists, err := FetchAndCacheArtists()
		if err != nil {
			log.Printf("Error fetching artists for URL %s: %v", r.URL.Path, err)
			http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
			return
		}

		filteredArtists := applyFilters(artists, filters)
		originalArtists, err := FetchAndCacheArtists()
		if err != nil {
			log.Printf("Error fetching artists for URL %s: %v", r.URL.Path, err)
			http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
			return
		}
		if len(filteredArtists) != len(originalArtists) {
			message = "Results according to your filters"
		}
		if shuffle {
			shuffleArtists(filteredArtists)
		}

		// Search for artists
		if searchQuery != "" {
			filteredArtists = SearchArtists(filteredArtists, searchQuery)
			message = "Search results for: " + searchQuery
		}
		if message == "" {
			message = "Welcome to Groupie Tracker"
		}
		sortArtists(filteredArtists)
		// Add the total number of artists to the context before pagination
		totalArtists := len(filteredArtists)
		ctx := context.WithValue(r.Context(), totalArtistsKey, totalArtists)

		// Paginate results
		start := 0
		end := pagination
		if end > len(filteredArtists) {
			end = len(filteredArtists)
		}
		paginatedArtists := filteredArtists[start:end]

		// Cache the response
		cachedData := CachedData{
			PaginatedArtists: paginatedArtists,
			TotalArtists:     totalArtists,
			Message:          message,
		}
		responseData, _ := json.Marshal(cachedData)
		queryCache.Set(cacheKey, responseData, 2*time.Minute)

		// Pass the data to the next handler
		ctx = context.WithValue(ctx, DataKey, paginatedArtists)
		ctx = context.WithValue(ctx, messageKey, message)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		next.ServeHTTP(w, r)
	})
}
