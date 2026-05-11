// internal/db/sqlite/user_repo.go
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/backend/models"
	"forum/internal/utils"
	"os"
	"time"
)

type UserRepo struct {
	DBlogger *DBlogger
	PostRepo *PostRepo
}

func NewUserRepo(db *sql.DB, logfile *os.File, post *PostRepo) *UserRepo {
	dblogger := &DBlogger{DB: db, logfile: logfile}
	return &UserRepo{DBlogger: dblogger, PostRepo: post}
}
func (r *UserRepo) Create(ctx context.Context, u models.User) error {
	const q = `INSERT INTO users(uuid, mail, username, password, type, createdAt, verified) VALUES(?,?,?,?,?,?,?);`
	_, err := r.DBlogger.LogExecContext(ctx, q, u.UUID, u.Mail, u.Username, u.Password, u.Role, u.CreationDate, u.Verified)
	return err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (models.User, error) {
	const q = `SELECT uuid, mail, username, password, type, createdAt, verified FROM users WHERE username = ?;`
	var u models.User
	err := r.DBlogger.DB.QueryRowContext(ctx, q, username).
		Scan(&u.UUID, &u.Mail, &u.Username, &u.Password, &u.Role, &u.CreationDate, &u.Verified)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("not found")
	}
	return u, err
}

// Provide FindByUUID for session-service
func (r *UserRepo) FindByUUID(ctx context.Context, uuid string) (models.User, error) {
	const q = `SELECT uuid, mail, username, password, type, createdAt, verified FROM users WHERE uuid = ?;`
	var u models.User
	err := r.DBlogger.DB.QueryRowContext(ctx, q, uuid).
		Scan(&u.UUID, &u.Mail, &u.Username, &u.Password, &u.Role, &u.CreationDate, &u.Verified)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("not found")
	}
	return u, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	const q = `SELECT uuid, mail, username, password, type, createdAt, verified FROM users WHERE mail = ?;`
	var u models.User
	err := r.DBlogger.DB.QueryRowContext(ctx, q, email).
		Scan(&u.UUID, &u.Mail, &u.Username, &u.Password, &u.Role, &u.CreationDate, &u.Verified)
	if err == sql.ErrNoRows {
		return models.User{}, errors.New("not found")
	}

	//Could not manage to initialize the struct properly because of nested slices of structs
	//therefore the GetProfile method is used only inside the handlers in which case
	//the user struct is casted from an interface from the context with value which is created from the middleware
	//In this way there is no problem with the initialization
	//if we manage to properly initialize it later that would be nice
	//if we had that in the middleware(right here)
	// err = r.GetProfile(ctx, &u)
	// if err != nil {
	// 	return models.User{}, errors.New("not found")
	// }
	//
	return u, err
}

func (r *UserRepo) DeleteByUUID(ctx context.Context, uuid string) error {
	const q = `DELETE FROM users WHERE uuid = ?;`
	_, err := r.DBlogger.LogExecContext(ctx, q, uuid)
	return err
}

func (r *UserRepo) Update(ctx context.Context, u models.User) error {
	const q = `UPDATE users SET mail = ?, username = ?, password = ?, type = ?, createdAt = ?, verified = ? WHERE uuid = ?;`
	_, err := r.DBlogger.LogExecContext(ctx, q, u.Mail, u.Username, u.Password, u.Role, u.CreationDate, u.Verified, u.UUID)
	return err
}

func (s *UserRepo) VerifyEmail(ctx context.Context, uuid string) error {
	user, err := s.FindByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.Verified {
		return fmt.Errorf("user already verified")
	}

	user.Verified = true
	return s.Update(ctx, user)
}

func (r *UserRepo) UpdateVerification(ctx context.Context, uuid string, verified bool) error {
	_, err := r.DBlogger.DB.ExecContext(ctx, `
        UPDATE users SET verified = ? WHERE uuid = ?
    `, verified, uuid)
	return err
}

func (r *UserRepo) SeeNotification(ctx context.Context, notId, userId string) error {
	query := "update notifications set seen=true where id==? "
	_, err := r.DBlogger.LogExecContext(ctx, query, notId, userId)
	return err
}

func (r *UserRepo) GetProfile(ctx context.Context, user *models.User) error {
	err := r.FindCreatedPosts(ctx, user)
	if err != nil {
		return err
	}
	err = r.FindLikedPosts(ctx, user)
	if err != nil {
		return err
	}
	err = r.FindUserNotifications(ctx, user)
	if err != nil {
		return err
	}
	err = r.FindUserActivities(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) FindCreatedPosts(ctx context.Context, user *models.User) error {
	user.CreatedPosts = make([]models.Post, 0)
	query := "SELECT * from posts where user_uuid==?"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var p models.Post
		var postime time.Time
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.User_UUID, &p.ImagePath, &postime); err != nil {
			return err
		}
		p.CreationDate = postime.Format("2006-01-02 15:04:05")
		err = r.PostRepo.FindPostInfo(ctx, &p)
		if err != nil {
			return err
		}
		user.CreatedPosts = append(user.CreatedPosts, p)
	}
	return nil
}

func (r *UserRepo) FindLikedPosts(ctx context.Context, user *models.User) error {
	user.LikedPosts = make([]models.Post, 0)
	query := "select distinct posts.id,posts.title,posts.content,posts.user_uuid,posts.image,posts.created_at from posts join reactions on posts.id==reactions.post_id  where reactions.user_uuid==? and like==true;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var p models.Post
		var postime time.Time
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.User_UUID, &p.ImagePath, &postime); err != nil {
			return err
		}

		p.CreationDate = postime.Format("2006-01-02 15:04:05")
		err = r.PostRepo.FindPostInfo(ctx, &p)
		if err != nil {
			return err
		}
		user.LikedPosts = append(user.LikedPosts, p)
	}
	return nil
}

func (r *UserRepo) FindUserNotifications(ctx context.Context, user *models.User) error {
	user.Notifications = make([]models.Notification, 0)
	err := r.FindNotificationsForPosts(ctx, user)
	if err != nil {
		return err
	}
	err = r.FindNotificationsForComments(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) FindNotificationsForPosts(ctx context.Context, user *models.User) error {
	query := "select notifications.id,users.username,action,post_id,seen,title,content,posts.user_uuid,image,posts.created_at from notifications join posts on notifications.post_id==posts.id join users on notifications.user_uuid==users.uuid where posts.user_uuid==? ;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var not models.Notification
		var postime time.Time
		if err := rows.Scan(&not.ID, &not.Username, &not.Action, &not.Post.ID, &not.Seen, &not.Post.Title, &not.Post.Content, &not.Post.User_UUID, &not.Post.ImagePath, &not.Post.CreationDate); err != nil {
			return err
		}

		not.Post.CreationDate = postime.Format("2006-01-02 15:04:05")
		err = r.PostRepo.FindPostInfo(ctx, &not.Post)
		if err != nil {
			return err
		}
		user.Notifications = append(user.Notifications, not)
	}
	return nil
}

// THIS IS NOT READY
func (r *UserRepo) FindNotificationsForComments(ctx context.Context, user *models.User) error {
	query := "select notifications.id,users.username,action,notifications.comment_id,seen,comments.content,comments.user_uuid,comments.created_at,comments.post_id,posts.title,posts.image from notifications join comments on notifications.comment_id==comments.id join posts on comments.post_id==posts.id join users on notifications.user_uuid==users.uuid where comments.user_uuid==? ;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var not models.Notification
		var postime time.Time
		if err := rows.Scan(&not.ID, &not.Username, &not.Action, &not.Comment.ID, &not.Seen, &not.Comment.Content, &not.Comment.User_UUID, &not.Comment.CreationDate, &not.Post.ID, &not.Post.Title, &not.Post.ImagePath); err != nil {
			return err
		}
		not.Comment.CreationDate = postime.Format("2006-01-02 15:04:05")
		err = r.PostRepo.FindCommentInfo(ctx, &not.Comment)
		if err != nil {
			return err
		}
		user.Notifications = append(user.Notifications, not)
	}
	return nil
}

// to do
func (r *UserRepo) FindUserActivities(ctx context.Context, user *models.User) error {
	user.Activities = make([]models.Activity, 0)
	err := r.FindActivitiesForPosts(ctx, user)
	if err != nil {
		return err
	}
	err = r.FindActivitiesForComments(ctx, user)
	if err != nil {
		return err
	}
	user.Activities = utils.SortActivities(user.Activities)
	return nil
}

func (r *UserRepo) FindActivitiesForPosts(ctx context.Context, user *models.User) error {
	query := "select action,posts.id,posts.title,posts.content,posts.user_uuid,posts.image,activities.created_at from activities join posts on activities.post_id==posts.id where activities.user_uuid==?;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var act models.Activity
		var postime time.Time
		if err := rows.Scan(&act.Action, &act.Post.ID, &act.Post.Title, &act.Post.Content, &act.Post.User_UUID, &act.Post.ImagePath, &postime); err != nil {
			return err
		}
		err = r.PostRepo.FindPostInfo(ctx, &act.Post)
		if err != nil {
			return err
		}
		act.CreationDate = postime.Format("2006-01-02 15:04:05")
		user.Activities = append(user.Activities, act)
	}
	return nil
}

func (r *UserRepo) FindActivitiesForComments(ctx context.Context, user *models.User) error {
	query := "select action,comments.id,comments.content,comments.user_uuid,posts.id,posts.title,posts.content,posts.user_uuid,posts.image,activities.created_at from activities join comments on activities.comment_id==comments.id join posts on comments.post_id==posts.id where activities.user_uuid==?"
	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var act models.Activity
		var postime time.Time
		if err := rows.Scan(&act.Action, &act.Comment.ID, &act.Comment.Content, &act.Comment.User_UUID, &act.Post.ID, &act.Post.Title, &act.Post.Content, &act.Post.User_UUID, &act.Post.ImagePath, &postime); err != nil {
			return err
		}
		err = r.PostRepo.FindPostInfo(ctx, &act.Post)
		if err != nil {
			return err
		}
		err = r.PostRepo.FindCommentInfo(ctx, &act.Comment)
		if err != nil {
			return err
		}
		act.CreationDate = postime.Format("2006-01-02 15:04:05")
		user.Activities = append(user.Activities, act)
	}
	return nil
}
