package scoreboard

import (
	"errors"
	"fmt"
	"strings"
)

const (
	defaultPageSize = 5
	maxPageSize     = 50
)

// score validation, storage, and ranking logic.
type Service struct {
	store *Store
}

// NewService creates a Service.
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// CreateScore validates and registers a new score submission.
func (s *Service) CreateScore(input ScoreInput) (CreateResult, error) {
	name := strings.TrimSpace(input.Name)
	if n := len(name); n == 0 || n > 20 {
		return CreateResult{}, ErrInvalidName
	}
	if input.Score < 0 {
		return CreateResult{}, ErrInvalidScore
	}

	seconds, err := ParseTimeString(input.Time)
	if err != nil {
		return CreateResult{}, err
	}

	record, err := s.store.Add(name, input.Score, seconds)
	if err != nil {
		return CreateResult{}, fmt.Errorf("persist score: %w", err)
	}
	records := s.store.Snapshot()
	ranked := RankScores(records)

	entry, ok := FindRankedScore(ranked, record.ID)
	if !ok {
		return CreateResult{}, errors.New("ranking lookup failed")
	}

	percentRounded := ComputePercentile(entry.Position, len(ranked))

	return CreateResult{
		Record:            record,
		Position:          entry.Position,
		PercentileRounded: percentRounded,
		Message:           buildPlacementMessage(record.Name, percentRounded, entry.Position),
	}, nil
}

// ListScores paginates scores and optionally includes percentile data for a record.
func (s *Service) ListScores(opts ListOptions) (ListResult, error) {
	page, pageSize, includeID := normalizeListOptions(opts)

	records := s.store.Snapshot()
	ranked := RankScores(records)

	totalCount := len(ranked)
	totalPages := computeTotalPages(totalCount, pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	} else if start > totalCount {
		start = totalCount
	}

	end := start + pageSize
	if end > totalCount {
		end = totalCount
	}

	sliced := ranked[start:end]
	items := make([]ListItem, len(sliced))
	for i, entry := range sliced {
		items[i] = ListItem(entry)
	}

	var subject *SubjectInfo
	if includeID > 0 {
		entry, ok := FindRankedScore(ranked, includeID)
		if !ok {
			return ListResult{}, ErrScoreNotFound
		}
		rounded := ComputePercentile(entry.Position, totalCount)
		subject = &SubjectInfo{
			Record:            entry.Record,
			Position:          entry.Position,
			PercentileRounded: rounded,
			Message:           buildPlacementMessage(entry.Record.Name, rounded, entry.Position),
		}
	}

	return ListResult{
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
		Items:      items,
		Subject:    subject,
	}, nil
}

func computeTotalPages(total, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func normalizeListOptions(opts ListOptions) (page, pageSize int, includeID int64) {
	page = opts.Page
	if page < 1 {
		page = 1
	}

	pageSize = opts.PageSize
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > maxPageSize:
		pageSize = maxPageSize
	}

	if opts.IncludeID > 0 {
		includeID = opts.IncludeID
	}
	return
}

func buildPlacementMessage(name string, percentile float64, position int) string {
	return fmt.Sprintf(
		"Congrats %s, you are in the top %.1f%%, on the %s position.",
		name,
		percentile,
		OrdinalSuffix(position),
	)
}
