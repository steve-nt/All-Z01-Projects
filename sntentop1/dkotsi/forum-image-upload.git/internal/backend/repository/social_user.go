package repository

import (
	"context"
	"forum-image-upload/internal/backend/models"
)

type SocialUserRepository interface {
	LinkSocialUser(ctx context.Context, su models.SocialUser) error
	FindByProvider(ctx context.Context, provider string, providerUserID string) (models.SocialUser, error)
}
