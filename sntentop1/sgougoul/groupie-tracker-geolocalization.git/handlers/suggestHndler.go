package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"sgougoupractice/suggestions"
)

// Returns suggestions using fuzzing logic with every key press
func SuggestionsHandler(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("q")
	if len(query) > 100 {
		return &HTTPError{
			Status:  http.StatusBadRequest,
			Message: "bad request",
		}

	}
	suggestions.CacheLock.RLock()
	allSuggestions := make([]suggestions.TypeSuggestion, len(suggestions.SuggestionsCache))
	copy(allSuggestions, suggestions.SuggestionsCache)
	if len(allSuggestions) == 0 {
		return &HTTPError{
			Status:  http.StatusInternalServerError,
			Message: "server-error",
		}
	}
	suggestions.CacheLock.RUnlock()
	filtered := suggestions.FilterSuggestions(query, allSuggestions, 2)
	var builder strings.Builder
	builder.WriteString(`<datalist id="artistSuggestions">`)
	for _, s := range filtered {
		builder.WriteString(fmt.Sprintf(`<option value="%s - %s"></option>`, s.Label, s.Type))
	}
	builder.WriteString(`</datalist>`)
	w.Header().Set("Content-type", "text/html")
	fmt.Fprint(w, builder.String())
	return nil
}
