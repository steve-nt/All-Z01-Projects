// internal/repository/post.go
package repository

import (
	"context"
	"forum/internal/backend/models"
	"time"
)

type PostRepository interface {
	ListByCategory(ctx context.Context, postCategories []string) ([]models.Post, error)
	FindPostInfo(ctx context.Context, p *models.Post) error
	FindPostbyID(ctx context.Context, post_id string) (models.Post, error)
	CreatePost(ctx context.Context, postinfo *models.PostInfo, user_uuid string, createdAt int64) error
	EditPost(ctx context.Context, postinfo *models.PostInfo, post_id string, user_uuid string, createdAt int64) error
	EditComment(ctx context.Context, content string, comment_id string, user_uuid string, createdAt int64) error
	CreateComment(ctx context.Context, commentContent string, user_uuid string, post_id string, createdAt time.Time) error
	LikeButton(ctx context.Context, post_id string, user models.User, onwhat string) error
	DislikeButton(ctx context.Context, post_id string, user models.User, onwhat string) error
	RemovePost(ctx context.Context, user_uuid string, post_id string) error
	RemoveComment(ctx context.Context, user_uuid string, comment_id string) error
}
