package scoreboard

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseTimeString parses "mm:ss" into total seconds.
func ParseTimeString(value string) (int, error) {
	value = strings.TrimSpace(value)
	if len(value) != 5 || value[2] != ':' {
		return 0, ErrInvalidTimeFormat
	}

	minutesPart := value[:2]
	secondsPart := value[3:]

	minutes, err := strconv.Atoi(minutesPart)
	if err != nil || minutes < 0 {
		return 0, ErrInvalidTimeValue
	}

	seconds, err := strconv.Atoi(secondsPart)
	if err != nil || seconds < 0 || seconds > 59 {
		return 0, ErrInvalidTimeValue
	}

	return minutes*60 + seconds, nil
}

// FormatSeconds converts total seconds back to mm:ss.
func FormatSeconds(totalSeconds int) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
