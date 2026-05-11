package scoreboard

import "errors"

var (
	// ErrInvalidJSON indicates the payload was malformed.
	ErrInvalidJSON = errors.New("invalid JSON payload")
	// ErrInvalidName indicates name validation failed.
	ErrInvalidName = errors.New("name must be between 1 and 20 characters")
	// ErrInvalidScore indicates score validation failed.
	ErrInvalidScore = errors.New("score must be a non-negative integer")
	// ErrInvalidTimeFormat indicates the time string is not in mm:ss.
	ErrInvalidTimeFormat = errors.New("time must match mm:ss (e.g. 03:45)")
	// ErrInvalidTimeValue indicates seconds are out of range.
	ErrInvalidTimeValue = errors.New("time values must be non-negative with seconds < 60")
	// ErrScoreNotFound indicates the requested score isn't present.
	ErrScoreNotFound = errors.New("score not found")
)
