package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum-advanced-features/internal/backend/models"
	"forum-advanced-features/internal/backend/repository"

	"github.com/gofrs/uuid"
)

type SocialUserService struct {
	social_users repository.SocialUserRepository
	users        repository.UserRepository
}

func NewSocialUserService(s_ur repository.SocialUserRepository, u_r repository.UserRepository) *SocialUserService {
	return &SocialUserService{
		social_users: s_ur,
		users:        u_r,
	}
}

func (s *SocialUserService) SocialRegister(ctx context.Context, provider, providerUserID, email, username, role string) (models.User, error) {
	// Έλεγχος αν υπάρχει ήδη χρήστης με αυτόν τον πάροχο
	existingUser, err := s.social_users.FindByProvider(ctx, provider, providerUserID)
	if err == nil && existingUser.UUID != "" {
		return models.User{UUID: existingUser.UUID}, nil // Επιστροφή υπάρχοντος χρήστη
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.User{}, fmt.Errorf("find by provider failed: %w", err)
	}
	// Έλεγχος αν υπάρχει ήδη χρήστης με το ίδιο email
	existingByEmail, err := s.users.FindByEmail(ctx, email)
	if err == nil {
		// Ο χρήστης υπάρχει -> Σύνδεση μόνο στον πίνακα social_users
		su := models.SocialUser{
			UUID:           existingByEmail.UUID,
			Provider:       provider,
			ProviderUserID: providerUserID,
		}
		if err := s.social_users.LinkSocialUser(ctx, su); err != nil {
			return models.User{}, fmt.Errorf("link existing social user failed: %w", err)
		}
		return existingByEmail, nil
	} else {
		uid, _ := uuid.NewV4()
		u := models.User{
			UUID:         uid.String(),
			Mail:         email,
			Username:     username,
			Password:     "", // no password for social users
			Role:         "user",
			CreationDate: time.Now(),
			Verified:     true, // assume verified via provider
		}
		if err := s.users.Create(ctx, u); err != nil {
			return models.User{}, err
		}
		// Σύνδεση social user
		su := models.SocialUser{
			UUID:           uid.String(),
			Provider:       provider,
			ProviderUserID: providerUserID,
		}
		if err := s.social_users.LinkSocialUser(ctx, su); err != nil {
			return models.User{}, fmt.Errorf("link social user failed: %w", err)
		}
		return u, nil
	}
}
