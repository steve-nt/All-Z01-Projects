package repository

import (
	"context"
	"forum-authentication/internal/backend/models"
)

type SocialUserRepository interface {
	LinkSocialUser(ctx context.Context, su models.SocialUser) error
	FindByProvider(ctx context.Context, provider string, providerUserID string) (models.SocialUser, error)
}
