// internal/repository/user.go
package repository

import (
	"context"
	"forum-image-upload/internal/backend/models"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (models.User, error)
	SeeNotification(ctx context.Context, notId, userId string) error
	GetProfile(ctx context.Context, user *models.User) error
	FindByUUID(ctx context.Context, uuid string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)

	Create(ctx context.Context, u models.User) error
	DeleteByUUID(ctx context.Context, uuid string) error
	Update(ctx context.Context, u models.User) error
	VerifyEmail(ctx context.Context, uuid string) error
	UpdateVerification(ctx context.Context, uuid string, verified bool) error
}
