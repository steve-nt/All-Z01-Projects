package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var errNoEmail = errors.New("oauth: missing email")

func GetUserByProvider(provider, providerID string) (int, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	providerID = strings.TrimSpace(providerID)
	if provider == "" || providerID == "" {
		return 0, sql.ErrNoRows
	}

	var id int
	err := DB.QueryRow(
		`SELECT id FROM users WHERE provider = ? AND provider_id = ?`,
		provider, providerID,
	).Scan(&id)
	return id, err
}

func CreateOAuthUser(provider, providerID, email, username string) (int, error) {
	provider = strings.TrimSpace(strings.ToLower(provider))
	providerID = strings.TrimSpace(providerID)
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	if provider == "" || providerID == "" {
		return 0, fmt.Errorf("oauth: provider/provider_id required")
	}
	if email == "" {
		// We keep users.email as NOT NULL. For GitHub we fetch email via user:email scope.
		return 0, errNoEmail
	}
	if username == "" {
		username = "user"
	}

	username = MakeUniqueUsername(username)

	// Password is unused for OAuth users, but users.password is NOT NULL.
	passwordPlaceholder := ""

	// The first ever account may be created via OAuth.
	// In that case it becomes admin so the forum still has an administrator.
	role := nextAssignedRole()

	res, err := DB.Exec(
		`INSERT INTO users (email, username, password, provider, provider_id, role)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		email, username, passwordPlaceholder, provider, providerID, role,
	)
	if err != nil {
		return 0, err
	}

	id64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id64), nil
}

func FindOrCreateOAuthUser(provider, providerID, email, preferredUsername string) (int, error) {
	// 1) Existing OAuth identity?
	if id, err := GetUserByProvider(provider, providerID); err == nil {
		return id, nil
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	// 2) If email already exists (maybe user registered by password), we link OAuth to that user.
	// This keeps a single account per email.
	emailNorm := strings.TrimSpace(strings.ToLower(email))
	if emailNorm == "" {
		return 0, errNoEmail
	}

	var existingID int
	err := DB.QueryRow(`SELECT id FROM users WHERE email = ?`, emailNorm).Scan(&existingID)
	if err == nil {
		// Link provider fields to existing user row.
		_, uerr := DB.Exec(
			`UPDATE users SET provider = ?, provider_id = ? WHERE id = ?`,
			strings.ToLower(strings.TrimSpace(provider)),
			strings.TrimSpace(providerID),
			existingID,
		)
		if uerr != nil {
			return 0, uerr
		}
		return existingID, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// 3) Create a new OAuth user.
	return CreateOAuthUser(provider, providerID, email, preferredUsername)
}

func MakeUniqueUsername(base string) string {
	base = strings.TrimSpace(base)
	base = sanitizeUsername(base)
	if base == "" {
		base = "user"
	}

	// Try base first
	if _, err := GetUserByUsername(base); err == sql.ErrNoRows {
		return base
	}

	// Add small suffix until available
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 30; i++ {
		try := fmt.Sprintf("%s-%04d", base, rand.Intn(10000))
		if _, err := GetUserByUsername(try); err == sql.ErrNoRows {
			return try
		}
	}

	// Fallback
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}

func sanitizeUsername(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")

	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-' || r == '.':
			b.WriteRune(r)
		}
	}

	out := b.String()
	out = strings.Trim(out, "._-")
	if len(out) > 24 {
		out = out[:24]
	}
	return out
}