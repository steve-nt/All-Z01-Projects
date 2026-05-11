package models

import (
	"database/sql"
	"errors"
	"forum/src/utils"
	"net/mail"
	"regexp"
	"slices"
	"sort"
)

type User struct {
	Id                       int64
	Username                 string
	Hash                     string
	Email                    string
	LoggedIn                 bool
	OAuthProvider            string
	Notifications            Notifications
	UnreadNotificationsCount int
	Activities               Activities
}

func GetGuestUser() User {
	return User{
		Username: "guest",
		LoggedIn: false,
	}
}

// Returns ONLY the `User.Hash` field for comparison against the given password
func GetUserPasswordByEmail(email string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT hash FROM users WHERE email = ?`, email).Scan(&user.Hash)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id, email, username FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

func GetUserBySession(sessionValue string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id, email, username FROM users WHERE session_key = ?`, sessionValue).Scan(&user.Id, &user.Email, &user.Username)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	return user, nil
}

func GetUserByOAuthProviderAndEmail(provider, email string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id, email, username, oauth_provider FROM users WHERE oauth_provider = ? AND email = ?`, provider, email).Scan(&user.Id, &user.Email, &user.Username, &user.OAuthProvider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrorNoRows
		}
		return User{}, err
	}
	return user, nil
}

func getUserById(id int64) (User, error) {
	var user User
	err := db.QueryRow(`SELECT username FROM users WHERE id = ?`, id).Scan(&user.Username)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return User{}, err
	}
	user.Id = id
	return user, nil
}

func (u *User) ValidateUsername() error {
	unameMask := regexp.MustCompile(`^[a-zA-Z0-9_]{4,50}$`)
	if !unameMask.MatchString((*u).Username) {
		return ErrorInvalidUsername
	}
	return nil
}

func (u *User) ValidateEmail() error {
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ValidateUser() error {
	var err error
	if err = u.ValidateUsername(); err != nil {
		return err
	}
	if err = u.ValidateEmail(); err != nil {
		return err
	}
	return nil
}

func (u *User) Add() error {
	err := u.ValidateUser()
	if err != nil {
		return err
	}
	if IsEmailRegistered(u.Email) {
		return ErrorEmailIsRegistered
	}
	stmt, err := db.Prepare("INSERT INTO users (username, email, hash) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Username, u.Email, u.Hash)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) AddOAuth() error {
	if err := u.ValidateUser(); err != nil {
		return err
	}
	if IsEmailRegistered(u.Email) {
		return ErrorEmailIsRegistered
	}
	if !IsUniqueUsername(u.Username) {
		return ErrorUsernameTaken
	}
	stmt, err := db.Prepare("INSERT INTO users (username, email, oauth_provider) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.Username, u.Email, u.OAuthProvider)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetPosts() (Posts, error) {
	var posts Posts
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	WHERE user_id = ?`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func (u *User) GetLikedPosts() (Posts, error) {
	var posts Posts
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp
	FROM posts
	JOIN reactions r ON posts.id = r.post_id
	WHERE r.user_id = ? AND r.value = 1
	`, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Posts{}, err
		}
		post.TimestampString = utils.ConvertTimeToString(t)
		posts = append(posts, post)
	}
	return posts, nil
}

func (u *User) SetUserSession(session_key string) error {
	stmt, err := db.Prepare("UPDATE users SET session_key = ? WHERE id = ?")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec(session_key, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func IsUniqueUsername(username string) bool {
	usernames, err := GetAllUsernames()
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(usernames, username)
}

func IsUniqueEmail(email string) bool {
	emails, err := GetAllUserEmails()
	if err != nil {
		(&Error{}).Consume(err).LogError()
		return false
	}
	return !slices.Contains(emails, email)
}

func IsEmailRegistered(email string) bool {
	return !IsUniqueEmail(email)
}

// Check if user already liked this post
func HasUserLikedPost(userId, postId int64) (bool, error) {
	reactionId, err := CheckIfUserLikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this post
func HasUserDislikedPost(userId, postId int64) (bool, error) {
	reactionId, err := CheckIfUserDislikedPost(userId, postId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already liked this comment
func HasUserLikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := CheckIfUserLikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

// Check if user already disliked this comment
func HasUserDislikedComment(userId, commentId int64) (bool, error) {
	reactionId, err := CheckIfUserDislikedComment(userId, commentId)
	if err != nil {
		return false, err
	}
	return reactionId != 0, nil
}

func (u *User) GetNotifications() error {
	notifications, err := GetNotificationsByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	u.Notifications = notifications
	u.CountUnreadNotifications()
	return nil
}

func (u *User) MarkNotificationAsRead(notificationId int64) error {
	stmt, err := db.Prepare("UPDATE notifications SET read = 1 WHERE id = ? AND user_id = ?")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec(notificationId, (*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func (u *User) MarkAllNotificationsAsRead() error {
	stmt, err := db.Prepare("UPDATE notifications SET read = 1 WHERE user_id = ?")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec((*u).Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func (u *User) GetActivity() error {
	err := u.GetPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetLikedPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetDislikedPostsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetLikedCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	err = u.GetDislikedCommentsActivity()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	sort.Slice(u.Activities, func(i, j int) bool {
		return u.Activities[i].TimestampString > u.Activities[j].TimestampString
	})
	return nil
}

func (u *User) GetPostsActivity() error {
	posts, err := u.GetPosts()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
	}
	for _, post := range posts {
		err := post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
		}
		var activity Activity
		activity.TimestampString = post.TimestampString
		activity.Post = post
		activity.Type = "post"
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetCommentsActivity() error {
	comments, err := GetCommentsByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, comment := range comments {
		var activity Activity
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "comment"
		activity.Comment = comment
		activity.TimestampString = comment.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetLikedPostsActivity() error {
	reactions, err := GetPostLikesByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
	}
	for _, reaction := range reactions {
		var activity Activity
		var post = Post{Id: reaction.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "postLike"
		activity.TimestampString = reaction.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetDislikedPostsActivity() error {
	reactions, err := GetPostDislikesByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, reaction := range reactions {
		var activity Activity
		var post = Post{Id: reaction.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = post.GetReactionsByUserId((*u).Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "postDislike"
		activity.TimestampString = reaction.TimestampString
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetLikedCommentsActivity() error {
	reactions, err := GetCommentLikesByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, reaction := range reactions {
		var activity Activity
		var comment = Comment{Id: reaction.CommentId}
		err = comment.GetCommentById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "commentLike"
		activity.TimestampString = reaction.TimestampString
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) GetDislikedCommentsActivity() error {
	reactions, err := GetCommentDisikesByUserId(u.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, reaction := range reactions {
		var activity Activity
		var comment = Comment{Id: reaction.CommentId}
		err = comment.GetCommentById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactions()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		err = comment.GetReactionsByUserId(u.Id)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		var post = Post{Id: comment.PostId}
		err = post.GetById()
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
		activity.Type = "commentDislike"
		activity.TimestampString = reaction.TimestampString
		activity.Comment = comment
		activity.Post = post
		u.Activities = append(u.Activities, activity)
	}
	return nil
}

func (u *User) CountUnreadNotifications() {
	for _, notification := range u.Notifications {
		if !notification.Read {
			u.UnreadNotificationsCount++
		}
	}
}
