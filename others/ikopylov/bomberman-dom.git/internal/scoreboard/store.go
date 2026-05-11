package scoreboard

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// Store provides a thread-safe, file-backed store.
type Store struct {
	mu     sync.RWMutex
	scores []ScoreRecord
	nextID int64
	path   string
}

type persistedScore struct {
	ScoreRecord
	Position int `json:"position"`
}

// NewStore loads an existing JSON file (if present) or creates an empty store at the provided path.
func NewStore(path string) (*Store, error) {
	store := &Store{
		nextID: 1,
		path:   path,
	}

	if err := store.loadFromDisk(); err != nil {
		return nil, err
	}

	return store, nil
}

// Add appends a new score record, persists it to disk, and returns the stored record.
func (s *Store) Add(name string, score int, timeSeconds int) (ScoreRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := ScoreRecord{
		ID:          s.nextID,
		Name:        name,
		Score:       score,
		TimeSeconds: timeSeconds,
	}

	s.nextID++
	s.scores = append(s.scores, record)

	if err := s.persistLocked(); err != nil {
		return ScoreRecord{}, err
	}

	return record, nil
}

// Snapshot returns a copy of the in-memory slice for safe iteration.
func (s *Store) Snapshot() []ScoreRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make([]ScoreRecord, len(s.scores))
	copy(snapshot, s.scores)
	return snapshot
}

func (s *Store) loadFromDisk() error {
	if s.path == "" {
		return errors.New("store path cannot be empty")
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	var records []ScoreRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}

	s.scores = make([]ScoreRecord, len(records))
	copy(s.scores, records)

	var maxID int64
	for _, record := range s.scores {
		if record.ID > maxID {
			maxID = record.ID
		}
	}

	if maxID == 0 {
		maxID = int64(len(s.scores))
		for i := range s.scores {
			s.scores[i].ID = int64(i + 1)
		}
	}

	s.nextID = maxID + 1

	return nil
}

func (s *Store) persistLocked() error {
	dir := filepath.Dir(s.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	ranked := RankScores(s.scores)
	persisted := make([]persistedScore, len(ranked))
	for i, entry := range ranked {
		persisted[i] = persistedScore{
			ScoreRecord: entry.Record,
			Position:    entry.Position,
		}
	}

	data, err := json.MarshalIndent(persisted, "", "  ")
	if err != nil {
		return err
	}

	tempPath := s.path + ".tmp"
	if err := os.WriteFile(tempPath, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tempPath, s.path)
}
