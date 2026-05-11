// internal/db/sqlite/session_repo.go
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"forum-advanced-features/internal/backend/models"
	"log"
	"os"
)

type SessionRepo struct{ DBlogger *DBlogger }

func NewSessionRepo(db *sql.DB, logfile *os.File) *SessionRepo {
	dblogger := &DBlogger{DB: db, logfile: logfile}
	return &SessionRepo{DBlogger: dblogger}
}
func (r *SessionRepo) FindByCookie(ctx context.Context, cookie string) (models.Session, error) {
	const q = `SELECT uuid, cookie, createdAt, expiration, absoluteExpiration FROM sessions WHERE cookie = ?;`
	var s models.Session
	err := r.DBlogger.DB.QueryRowContext(ctx, q, cookie).
		Scan(&s.UUID, &s.CookieValue, &s.CreationDate, &s.Expiration, &s.AbsoluteExpiration)
	if err == sql.ErrNoRows {
		return models.Session{}, errors.New("not found")
	}
	return s, err
}

func (r *SessionRepo) FindByUUID(ctx context.Context, uuid string) (models.Session, error) {
	const q = `SELECT uuid, cookie, createdAt, expiration, absoluteExpiration FROM sessions WHERE uuid = ?;`
	var s models.Session
	err := r.DBlogger.DB.QueryRowContext(ctx, q, uuid).
		Scan(&s.UUID, &s.CookieValue, &s.CreationDate, &s.Expiration, &s.AbsoluteExpiration)
	if err == sql.ErrNoRows {
		return models.Session{}, errors.New("not found")
	}
	return s, err
}

func (r *SessionRepo) Create(ctx context.Context, s models.Session) error {
	const q = `INSERT INTO sessions (uuid, cookie, createdAt, expiration, absoluteExpiration) VALUES(?,?,?,?,?);`
	_, err := r.DBlogger.LogExecContext(ctx, q, s.UUID, s.CookieValue, s.CreationDate, s.Expiration, s.AbsoluteExpiration)
	if err != nil {
		log.Println("error in session_repo Create", err)
	}
	return err
}

func (r *SessionRepo) DeleteByCookie(ctx context.Context, cookie string) error {
	const q = `DELETE FROM sessions WHERE cookie = ?;`
	_, err := r.DBlogger.LogExecContext(ctx, q, cookie)
	return err
}
