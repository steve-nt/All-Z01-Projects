// internal/repository/session.go
package repository

import (
	"context"
	"forum-image-upload/internal/backend/models"
)

type SessionRepository interface {
	FindByCookie(ctx context.Context, cookie string) (models.Session, error)
	FindByUUID(ctx context.Context, uuid string) (models.Session, error)
	Create(ctx context.Context, s models.Session) error
	DeleteByCookie(ctx context.Context, cookie string) error
}
