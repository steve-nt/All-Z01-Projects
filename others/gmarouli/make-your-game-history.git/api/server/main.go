package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Try multiple possible paths for the scores file
func getScoresFilePath() string {
	// Try paths relative to different possible working directories
	possiblePaths := []string{
		"api/server/api/server/data/scores.json", // If run from project root with nested structure
		"api/server/data/scores.json",            // If run from project root
		"data/scores.json",                       // If run from api/server directory
	}
	
	// Check which path exists
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Default to the first path (will create directory if needed)
	return possiblePaths[0]
}

var scoresFilePath = getScoresFilePath()

// Score represents a single leaderboard submission.
type Score struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Score       int       `json:"score"`
	TimeSeconds int       `json:"timeSeconds"`
	CreatedAt   time.Time `json:"createdAt"`
}

type scoreStore struct {
	mu       sync.RWMutex
	scores   []Score
	nextID   int
	filePath string
}

func newScoreStore(filePath string) (*scoreStore, error) {
	store := &scoreStore{
		nextID:   1,
		filePath: filePath,
	}
	if err := store.loadFromFile(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *scoreStore) add(name string, scoreVal, timeSeconds int) (Score, int, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := Score{
		ID:          s.nextID,
		Name:        name,
		Score:       scoreVal,
		TimeSeconds: timeSeconds,
		CreatedAt:   time.Now().UTC(),
	}
	s.nextID++
	s.scores = append(s.scores, entry)
	if err := s.persistLocked(); err != nil {
		s.scores = s.scores[:len(s.scores)-1]
		s.nextID--
		return Score{}, 0, 0, err
	}

	sorted := s.sortedScoresLocked()
	rank := rankForID(sorted, entry.ID)
	percentile := computePercentile(rank, len(sorted))

	return entry, rank, percentile, nil
}

func (s *scoreStore) loadFromFile() error {
	if s.filePath == "" {
		return nil
	}
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("scores file not found at %s, starting with empty scores", s.filePath)
			return nil
		}
		return err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		log.Printf("scores file at %s is empty, starting with empty scores", s.filePath)
		return nil
	}
	var stored []Score
	if err := json.Unmarshal(data, &stored); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores = append([]Score(nil), stored...)
	maxID := 0
	for _, sc := range stored {
		if sc.ID > maxID {
			maxID = sc.ID
		}
	}
	s.nextID = maxID + 1
	if s.nextID <= 1 {
		s.nextID = 1
	}
	log.Printf("loaded %d scores from %s (next ID: %d)", len(stored), s.filePath, s.nextID)
	return nil
}

func (s *scoreStore) persistLocked() error {
	if s.filePath == "" {
		return nil
	}
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("failed to create directory %s: %v", dir, err)
		return err
	}
	
	// Get absolute path for logging
	absPath, _ := filepath.Abs(s.filePath)
	log.Printf("persisting %d scores to %s (absolute: %s)", len(s.scores), s.filePath, absPath)
	
	tmp, err := os.CreateTemp(dir, "scores-*.tmp")
	if err != nil {
		log.Printf("failed to create temp file in %s: %v", dir, err)
		return err
	}
	tmpPath := tmp.Name()
	encoder := json.NewEncoder(tmp)
	encoder.SetIndent("", "  ")
	records := s.scores
	if records == nil {
		records = []Score{}
	}
	if err := encoder.Encode(records); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		log.Printf("failed to encode scores: %v", err)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		log.Printf("failed to sync temp file: %v", err)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		log.Printf("failed to close temp file: %v", err)
		return err
	}
	if err := os.Rename(tmpPath, s.filePath); err != nil {
		os.Remove(tmpPath)
		log.Printf("failed to rename temp file to %s: %v", s.filePath, err)
		return err
	}
	log.Printf("successfully persisted scores to %s", s.filePath)
	return nil
}

func (s *scoreStore) sortedScoresLocked() []Score {
	c := make([]Score, len(s.scores))
	copy(c, s.scores)
	sort.Slice(c, func(i, j int) bool {
		if c[i].Score == c[j].Score {
			return c[i].CreatedAt.Before(c[j].CreatedAt)
		}
		return c[i].Score > c[j].Score
	})
	return c
}

type scoreListItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	TimeSeconds int    `json:"timeSeconds"`
	Rank        int    `json:"rank"`
}

func (s *scoreStore) page(page, size int) ([]scoreListItem, int, int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if size <= 0 {
		size = 5
	}
	if page <= 0 {
		page = 1
	}

	sorted := s.sortedScoresLocked()
	totalItems := len(sorted)

	totalPages := 1
	if totalItems > 0 {
		totalPages = (totalItems + size - 1) / size
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * size
	if start > totalItems {
		start = totalItems
	}

	end := start + size
	if end > totalItems {
		end = totalItems
	}

	items := make([]scoreListItem, 0, end-start)
	for i := start; i < end; i++ {
		entry := sorted[i]
		items = append(items, scoreListItem{
			ID:          entry.ID,
			Name:        entry.Name,
			Score:       entry.Score,
			TimeSeconds: entry.TimeSeconds,
			Rank:        i + 1,
		})
	}

	return items, totalItems, totalPages, page
}

func rankForID(scores []Score, id int) int {
	for i, s := range scores {
		if s.ID == id {
			return i + 1
		}
	}
	return len(scores)
}

func computePercentile(rank, total int) int {
	if total <= 0 || rank <= 0 {
		return 0
	}
	return ((rank - 1) * 100) / total
}

type scoreHandler struct {
	store *scoreStore
}

type postScoreRequest struct {
	Name        string `json:"name"`
	Score       int    `json:"score"`
	TimeSeconds int    `json:"timeSeconds"`
}

type postScoreResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	TimeSeconds int    `json:"timeSeconds"`
	Rank        int    `json:"rank"`
	Percentile  int    `json:"percentile"`
}

type scoresResponse struct {
	Items      []scoreListItem `json:"items"`
	Page       int             `json:"page"`
	Size       int             `json:"size"`
	TotalItems int             `json:"totalItems"`
	TotalPages int             `json:"totalPages"`
}

func (h *scoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *scoreHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, 1<<20)
	defer body.Close()

	var req postScoreRequest
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	req.Name = sanitizeName(req.Name)
	if req.Score < 0 || req.TimeSeconds < 0 {
		http.Error(w, "score and timeSeconds must be non-negative", http.StatusBadRequest)
		return
	}

	entry, rank, percentile, err := h.store.add(req.Name, req.Score, req.TimeSeconds)
	if err != nil {
		log.Printf("failed to persist score: %v", err)
		http.Error(w, "failed to save score", http.StatusInternalServerError)
		return
	}
	log.Printf("saved score: name=%s, score=%d, timeSeconds=%d, id=%d, rank=%d", entry.Name, entry.Score, entry.TimeSeconds, entry.ID, rank)

	response := postScoreResponse{
		ID:          entry.ID,
		Name:        entry.Name,
		Score:       entry.Score,
		TimeSeconds: entry.TimeSeconds,
		Rank:        rank,
		Percentile:  percentile,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *scoreHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	page, err := parseIntDefault(r.URL.Query().Get("page"), 1)
	if err != nil {
		http.Error(w, "invalid page parameter", http.StatusBadRequest)
		return
	}

	size, err := parseIntDefault(r.URL.Query().Get("size"), 5)
	if err != nil {
		http.Error(w, "invalid size parameter", http.StatusBadRequest)
		return
	}

	items, totalItems, totalPages, resolvedPage := h.store.page(page, size)
	resp := scoresResponse{
		Items:      items,
		Page:       resolvedPage,
		Size:       size,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	writeJSON(w, http.StatusOK, resp)
}

func parseIntDefault(value string, def int) (int, error) {
	if strings.TrimSpace(value) == "" {
		return def, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func sanitizeName(raw string) string {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "Anon"
	}
	if len(name) > 32 {
		return name[:32]
	}
	return name
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	// Allow requests from common localhost ports
	allowedOrigins := []string{
		"http://localhost:8080",
		"http://localhost:8000",
		"http://127.0.0.1:8080",
		"http://127.0.0.1:8000",
	}
	
	// Check if the origin is in the allowed list
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Vary", "Origin")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("error writing response: %v", err)
	}
}

func main() {
	log.Printf("initializing score store with file path: %s", scoresFilePath)
	store, err := newScoreStore(scoresFilePath)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/scores", &scoreHandler{store: store})

	server := &http.Server{
		Addr:              ":8090",
		Handler:           loggingMiddleware(mux),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Scoreboard API listening on :8090")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
