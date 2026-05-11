package scoreboard

import (
	"fmt"
	"math"
	"sort"
)

// RankScores performs a sequential ranking on the given records (score desc).
func RankScores(records []ScoreRecord) []RankedScore {
	if len(records) == 0 {
		return nil
	}

	snapshot := make([]ScoreRecord, len(records))
	copy(snapshot, records)

	sort.SliceStable(snapshot, func(i, j int) bool {
		if snapshot[i].Score == snapshot[j].Score {
			return snapshot[i].ID < snapshot[j].ID
		}
		return snapshot[i].Score > snapshot[j].Score
	})

	ranked := make([]RankedScore, len(snapshot))
	for i, record := range snapshot {
		ranked[i] = RankedScore{
			Record:   record,
			Position: i + 1,
		}
	}

	return ranked
}

// FindRankedScore locates a ranked entry by ID.
func FindRankedScore(ranked []RankedScore, id int64) (RankedScore, bool) {
	for _, entry := range ranked {
		if entry.Record.ID == id {
			return entry, true
		}
	}
	return RankedScore{}, false
}

// ComputePercentile converts a sequential rank into a percentile across the full leaderboard.
// The first-ranked player receives ~0% and the last-ranked approaches 100%.
func ComputePercentile(position, total int) float64 {
	if position <= 0 || total <= 0 {
		return 0
	}
	if position > total {
		position = total
	}

	if total == 1 {
		return 0
	}

	percentage := (float64(position-1) / float64(total-1)) * 100
	if percentage < 0 {
		percentage = 0
	} else if percentage > 100 {
		percentage = 100
	}
	return math.Round(percentage*100) / 100
}

// OrdinalSuffix converts a rank into a human-friendly ordinal string.
func OrdinalSuffix(n int) string {
	if n <= 0 {
		return "0th"
	}

	if n%100 >= 11 && n%100 <= 13 {
		return fmt.Sprintf("%d%s", n, "th")
	}

	switch n % 10 {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}
