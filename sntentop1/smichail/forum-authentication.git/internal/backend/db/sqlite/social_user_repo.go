package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-authentication/internal/backend/models"
	"os"
)

type SocialUserRepo struct {
	DBlogger *DBlogger
	UserRepo *UserRepo
}

func NewSocialUserRepo(db *sql.DB, logfile *os.File, userRepo *UserRepo) *SocialUserRepo {
	dblogger := &DBlogger{DB: db, logfile: logfile}
	return &SocialUserRepo{DBlogger: dblogger, UserRepo: userRepo}
}

func (r *SocialUserRepo) FindByProvider(ctx context.Context, provider, providerUserID string) (models.SocialUser, error) {
	const q = `SELECT user_uuid FROM social_users WHERE provider = ? AND provider_user_id = ?;`

	var userUUID string
	err := r.DBlogger.DB.QueryRowContext(ctx, q, provider, providerUserID).Scan(&userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SocialUser{}, sql.ErrNoRows // let service handle "not found"
		}
		return models.SocialUser{}, fmt.Errorf("FindByProvider query failed: %w", err)
	}

	return models.SocialUser{
		UUID:           userUUID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	}, nil
}

func (r *SocialUserRepo) LinkSocialUser(ctx context.Context, su models.SocialUser) error {
	const q = `
	INSERT INTO social_users(user_uuid, provider, provider_user_id)
	VALUES(?,?,?)
	ON CONFLICT(provider, provider_user_id) DO NOTHING;
	`

	_, err := r.DBlogger.LogExecContext(ctx, q, su.UUID, su.Provider, su.ProviderUserID)
	if err != nil {
		return fmt.Errorf("LinkSocialUser exec failed: %w", err)
	}
	return nil
}
