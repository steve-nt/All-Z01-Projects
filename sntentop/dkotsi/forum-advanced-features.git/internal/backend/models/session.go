package models

import "time"

type Session struct {
	UUID               string
	CookieValue        string
	CreationDate       time.Time
	Expiration         time.Time // Sliding window expiration (e.g., 30 minutes of inactivity)
	AbsoluteExpiration time.Time // Hard limit expiration (e.g., 24 hours from creation)
}
