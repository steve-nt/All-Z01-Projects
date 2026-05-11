package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"web/internal/scoreboard"
)

// Handler wires HTTP routes to the scoreboard service.
type Handler struct {
	svc   *scoreboard.Service
	lobby *lobby
}

// NewHandler creates a Handler instance.
func NewHandler(svc *scoreboard.Service) *Handler {
	return &Handler{svc: svc, lobby: newLobby()}
}

// RegisterRoutes attaches endpoints to the mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /scores", h.postScore)
	mux.HandleFunc("GET /scores", h.getScores)
	mux.HandleFunc("GET /ws", h.lobby.handleWebSocket)
	go h.lobby.monitorStart()
}

func (h *Handler) postScore(w http.ResponseWriter, r *http.Request) {
	var input scoreboard.ScoreInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, scoreboard.ErrInvalidJSON.Error())
		return
	}

	result, err := h.svc.CreateScore(input)
	if err != nil {
		status, message := mapCreateScoreError(err)
		writeJSONError(w, status, message)
		return
	}

	writeJSON(w, http.StatusCreated, formatCreatePayload(result))
}

func (h *Handler) getScores(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := parsePositiveInt(query.Get("page"), 1)
	pageSize := parsePositiveInt(query.Get("pageSize"), 5)

	sortParam := strings.TrimSpace(query.Get("sort"))
	if sortParam != "" && !strings.EqualFold(sortParam, "desc") {
		writeJSONError(w, http.StatusBadRequest, "only sort=desc is supported")
		return
	}

	includeIDStr := strings.TrimSpace(query.Get("includePercentileForId"))
	var includeID int64
	if includeIDStr != "" {
		id, ok := parsePositiveInt64(includeIDStr)
		if !ok {
			writeJSONError(w, http.StatusBadRequest, "includePercentileForId must be a positive integer")
			return
		}
		includeID = id
	}

	result, err := h.svc.ListScores(scoreboard.ListOptions{
		Page:      page,
		PageSize:  pageSize,
		IncludeID: includeID,
	})
	if err != nil {
		if errors.Is(err, scoreboard.ErrScoreNotFound) {
			writeJSONError(w, http.StatusNotFound, "score not found for includePercentileForId")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "unable to list scores")
		return
	}

	response := map[string]any{
		"totalCount": result.TotalCount,
		"totalPages": result.TotalPages,
		"page":       result.Page,
		"pageSize":   result.PageSize,
		"items":      formatListItems(result.Items),
	}

	if result.Subject != nil {
		response["subject"] = formatSubjectPayload(*result.Subject)
	}

	writeJSON(w, http.StatusOK, response)
}

func parsePositiveInt(value string, defaultVal int) int {
	if value == "" {
		return defaultVal
	}
	num, err := strconv.Atoi(value)
	if err != nil || num <= 0 {
		return defaultVal
	}
	return num
}

func parsePositiveInt64(value string) (int64, bool) {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
	if status >= 500 {
		fmt.Printf("server error (%d): %s\n", status, message)
	}
}

func mapCreateScoreError(err error) (int, string) {
	switch {
	case errors.Is(err, scoreboard.ErrInvalidName):
		return http.StatusBadRequest, scoreboard.ErrInvalidName.Error()
	case errors.Is(err, scoreboard.ErrInvalidScore):
		return http.StatusBadRequest, scoreboard.ErrInvalidScore.Error()
	case errors.Is(err, scoreboard.ErrInvalidTimeFormat):
		return http.StatusBadRequest, scoreboard.ErrInvalidTimeFormat.Error()
	case errors.Is(err, scoreboard.ErrInvalidTimeValue):
		return http.StatusBadRequest, scoreboard.ErrInvalidTimeValue.Error()
	default:
		return http.StatusInternalServerError, "unable to store score"
	}
}

func formatCreatePayload(result scoreboard.CreateResult) map[string]any {
	record := result.Record
	return map[string]any{
		"id":                record.ID,
		"name":              record.Name,
		"score":             record.Score,
		"time":              scoreboard.FormatSeconds(record.TimeSeconds),
		"position":          scoreboard.OrdinalSuffix(result.Position),
		"percentileRounded": result.PercentileRounded,
		"message":           result.Message,
	}
}

func formatListItems(items []scoreboard.ListItem) []map[string]any {
	formatted := make([]map[string]any, len(items))
	for i, item := range items {
		formatted[i] = map[string]any{
			"id":       item.Record.ID,
			"name":     item.Record.Name,
			"score":    item.Record.Score,
			"time":     scoreboard.FormatSeconds(item.Record.TimeSeconds),
			"position": scoreboard.OrdinalSuffix(item.Position),
		}
	}
	return formatted
}

func formatSubjectPayload(subject scoreboard.SubjectInfo) map[string]any {
	return map[string]any{
		"id":                subject.Record.ID,
		"position":          scoreboard.OrdinalSuffix(subject.Position),
		"percentileRounded": subject.PercentileRounded,
		"message":           subject.Message,
	}
}
