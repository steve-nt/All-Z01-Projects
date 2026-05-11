// internal/services/session_service.go
package services

import (
	"context"
	"errors"
	"time"

	"forum/internal/backend/models"
	"forum/internal/backend/repository"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	sessions    repository.SessionRepository
	users       repository.UserRepository
	IdleTTL     time.Duration // e.g., 30 * time.Minute - sliding window expiration
	AbsoluteTTL time.Duration // e.g., 24 * time.Hour - hard limit expiration
}

func NewSessionService(sr repository.SessionRepository, ur repository.UserRepository, idleTTL, absoluteTTL time.Duration) *SessionService {
	return &SessionService{sessions: sr, users: ur, IdleTTL: idleTTL, AbsoluteTTL: absoluteTTL}
}

func (s *SessionService) CreateOrGet(ctx context.Context, userUUID string) (string, error) {
	if sess, err := s.sessions.FindByUUID(ctx, userUUID); err == nil && sess.UUID != "" {
		//IMPORTANT:
		//the below lines of code where commented out in order to sustain only one active session when logging in with a user that is already logged in, if some issue comes up regarding the sessin management these line of code should be reviewed

		// // Return existing cookie value if still valid (both idle and absolute expiration)
		// now := time.Now()
		// if now.Before(sess.Expiration) && now.Before(sess.AbsoluteExpiration) {
		// 	return sess.CookieValue, nil
		// }
		// Expired: delete & recreate
		_ = s.sessions.DeleteByCookie(ctx, sess.CookieValue)
	}
	id, _ := uuid.NewV4()
	cookieValue := id.String()
	now := time.Now()
	newSession := models.Session{
		UUID:               userUUID,
		CookieValue:        cookieValue,
		CreationDate:       now,
		Expiration:         now.Add(s.IdleTTL),
		AbsoluteExpiration: now.Add(s.AbsoluteTTL),
	}
	if err := s.sessions.Create(ctx, newSession); err != nil {
		return "", err
	}
	return cookieValue, nil
}

func (s *SessionService) ValidateAndMaybeRefresh(ctx context.Context, cookie string) (models.User, string, error) {
	sess, err := s.sessions.FindByCookie(ctx, cookie)
	if err != nil {
		return models.User{}, "", err
	}

	now := time.Now()

	// Check absolute expiration first - cannot be extended
	if now.After(sess.AbsoluteExpiration) {
		_ = s.sessions.DeleteByCookie(ctx, cookie)
		return models.User{}, "", errors.New("session has reached absolute expiration limit")
	}

	// Sliding window: if expiring within window, rotate cookie + extend idle TTL
	// but preserve the original absolute expiration
	if now.Add(s.IdleTTL / 2).After(sess.Expiration) {
		_ = s.sessions.DeleteByCookie(ctx, cookie)
		id, _ := uuid.NewV4()
		newCookie := id.String()
		sess.CookieValue = newCookie
		sess.CreationDate = now
		sess.Expiration = now.Add(s.IdleTTL)
		// Keep the original absolute expiration - this is the key security feature
		// sess.AbsoluteExpiration stays the same
		if err := s.sessions.Create(ctx, sess); err != nil {
			return models.User{}, "", err
		}
		u, err := s.users.FindByUUID(ctx, sess.UUID)
		return u, newCookie, err
	}

	// Still valid, return current user
	u, err := s.users.FindByUUID(ctx, sess.UUID)
	return u, "", err
}

func (s *SessionService) DeleteByCookie(ctx context.Context, cookie string) error {
	return s.sessions.DeleteByCookie(ctx, cookie)
}
