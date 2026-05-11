package scoreboard

// ScoreInput represents the payload submitted by the client.
type ScoreInput struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Time  string `json:"time"`
}

// ScoreRecord stores the canonical, server-side representation of a score.
type ScoreRecord struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	TimeSeconds int    `json:"timeSeconds"`
}

// RankedScore embeds the score together with its sequential rank position.
type RankedScore struct {
	Record   ScoreRecord
	Position int
}

// CreateResult is returned after inserting a score.
type CreateResult struct {
	Record            ScoreRecord
	Position          int
	PercentileRounded float64
	Message           string
}

// ListOptions controls pagination and percentile context.
type ListOptions struct {
	Page      int
	PageSize  int
	IncludeID int64
}

// ListItem is returned during pagination.
type ListItem struct {
	Record   ScoreRecord
	Position int
}

// SubjectInfo holds percentile and rank data for a specific score.
type SubjectInfo struct {
	Record            ScoreRecord
	Position          int
	PercentileRounded float64
	Message           string
}

// ListResult contains paginated leaderboard output.
type ListResult struct {
	TotalCount int
	TotalPages int
	Page       int
	PageSize   int
	Items      []ListItem
	Subject    *SubjectInfo
}
