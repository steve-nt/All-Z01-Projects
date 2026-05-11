package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum-authentication/internal/backend/models"
	"forum-authentication/internal/utils"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type PostRepo struct{ DBlogger *DBlogger }

func NewPostRepo(db *sql.DB, logfile *os.File) *PostRepo {
	dblogger := &DBlogger{DB: db, logfile: logfile}
	return &PostRepo{DBlogger: dblogger}
}

func (r *PostRepo) ListByCategory(ctx context.Context, postCategories []string) ([]models.Post, error) {
	q := `SELECT id, title, content, user_uuid, image, created_at FROM posts;` // avoid SELECT *

	if postCategories != nil {
		q = "select distinct posts.* from posts join post_categories on posts.id==post_categories.post_id join categories on category_id==categories.id where "

		q = utils.DesignQueryBasedOnCategories(q, postCategories)
	}
	rows, err := r.DBlogger.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Post
	for rows.Next() {
		var p models.Post
		var postime time.Time
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.User_UUID, &p.ImagePath, &postime); err != nil {
			return nil, err
		}
		p.CreationDate = postime.Format("2006-01-02 15:04:05")
		if err := r.FindPostInfo(ctx, &p); err != nil {
			log.Println(err)
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PostRepo) CreatePost(ctx context.Context, postinfo *models.PostInfo, user_uuid string, createdAt int64) error {
	post_id, _ := uuid.NewV4()

	const q = `INSERT INTO posts(id,title, content, user_uuid,image,created_at) VALUES(?,?,?,?,?, datetime(?, 'unixepoch'));`

	_, err := r.DBlogger.LogExecContext(ctx, q, post_id, postinfo.Title, postinfo.Content, user_uuid, postinfo.ImagePath, createdAt)
	if err != nil {
		return err
	}

	err = r.InsertPostToPost_CategoriesTable(ctx, post_id.String(), postinfo.Categories)
	if err != nil {
		return err
	}

	err = r.AddActivity(ctx, user_uuid, "created", post_id.String(), "post")
	if err != nil {
		return err
	}

	return err
}

func (r *PostRepo) CompareUserToCreator(ctx context.Context, user_uuid string, post_id string, onwhat string) (IstheSame bool, err error) {
	query := fmt.Sprintf("select user_uuid from %vs where id==?", onwhat)
	rows, err := r.DBlogger.DB.QueryContext(ctx, query, post_id)
	if err != nil {
		return true, err
	}

	var userid string
	counter := 0
	for rows.Next() {
		if err := rows.Scan(&userid); err != nil {
			return true, err
		}
		if counter >= 1 {
			return true, errors.New("More than one times the same post or comment id")
		}
		counter++
	}

	if userid == user_uuid {
		return true, nil
	}
	return false, nil
}

// Notifications are added here and each user will have his own notifications which will be seen or not
// To-do dont do the notification if the post is liked from the creator
func (r *PostRepo) AddNotification(ctx context.Context, user_uuid string, action string, post_id string, onwhat string) error {
	UserIsCreator, err := r.CompareUserToCreator(ctx, user_uuid, post_id, onwhat)
	if err != nil {
		return err
	}
	if UserIsCreator {
		return nil
	}
	query := fmt.Sprintf("insert into notifications(user_uuid,action,%v_id,seen) values(?,?,?,false);", onwhat)

	_, err = r.DBlogger.LogExecContext(ctx, query, user_uuid, action, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) RemoveNotification(ctx context.Context, user_uuid string, action string, post_id string, onwhat string) error {
	UserIsCreator, err := r.CompareUserToCreator(ctx, user_uuid, post_id, onwhat)
	if err != nil {
		return err
	}
	if UserIsCreator {
		return nil
	}
	query := fmt.Sprintf("delete from notifications where user_uuid==? and action==? and %v_id==?;", onwhat)

	_, err = r.DBlogger.LogExecContext(ctx, query, user_uuid, action, post_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) CreateComment(ctx context.Context, commentContent string, user_uuid string, post_id string, createdAt time.Time) error {
	q := "insert into comments(id,content,user_uuid,post_id,created_at) values(?,?,?,?,?);"

	comment_id, _ := uuid.NewV4()
	_, err := r.DBlogger.LogExecContext(ctx, q, comment_id, commentContent, user_uuid, post_id, createdAt)
	if err != nil {
		return err
	}

	err = r.AddNotification(ctx, user_uuid, "commented on your post", post_id, "post")
	if err != nil {
		return err
	}

	err = r.AddActivity(ctx, user_uuid, "commented", comment_id.String(), "comment")
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) InsertPostToPost_CategoriesTable(ctx context.Context, post_id string, categories []string) error {
	category_ids, err := r.ReturnCategoriesId(ctx, categories)
	if err != nil {
		return err
	}
	if category_ids == nil {
		return nil
	}
	q := "insert into post_categories (post_id,category_id) values(?,?);"
	for _, cat_id := range category_ids {
		_, err := r.DBlogger.LogExecContext(ctx, q, post_id, cat_id)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (r *PostRepo) ReturnCategoriesId(ctx context.Context, categories []string) (ids []int, err error) {
	q := "select id from categories where "
	q1 := utils.DesignQueryBasedOnCategories(q, categories)

	if q == q1 {
		return nil, nil
	}
	rows, err := r.DBlogger.DB.QueryContext(ctx, q1)
	if err != nil {
		return ids, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PostRepo) FindPostInfo(ctx context.Context, p *models.Post) error {

	err := r.FindPostNumofComments(ctx, p)
	if err != nil {
		return err
	}
	err = r.FindPostNumofLikes(ctx, p)
	if err != nil {
		return err
	}
	err = r.FindPostNumofDislikes(ctx, p)
	if err != nil {
		return err
	}
	err = r.FindPostCategories(ctx, p)
	if err != nil {
		return err
	}
	err = r.FindPostCreator(ctx, p)
	if err != nil {
		return err
	}
	err = r.FindIfUserReacted(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) FindIfUserReacted(ctx context.Context, p *models.Post) error {
	user := ctx.Value("user").(models.User)
	if user.Username == "" {
		return nil
	}
	query := "select like,dislike from reactions where user_uuid==? and post_id==?"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var liked bool
	var disliked bool

	counter := 0
	for rows.Next() {
		if counter >= 1 {
			return errors.New("Unexpected query result: User has more than one reactions to a post")
		}
		if err := rows.Scan(&liked, &disliked); err != nil {
			return err
		}
		counter++
	}

	if !liked && !disliked {
		return nil
	}
	switch liked {
	case true:
		p.Liked = true
		return nil
	case false:
		p.Disliked = true
		return nil
	default:
		return nil
	}
}

func (r *PostRepo) FindPostNumofComments(ctx context.Context, p *models.Post) error {
	query := "select count(*) from comments where post_id==?;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.NumberOfComments); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindPostNumofLikes(ctx context.Context, p *models.Post) error {
	query := "select count(*) from reactions where post_id==? and like==TRUE;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.NumberOfLikes); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindPostNumofDislikes(ctx context.Context, p *models.Post) error {
	query := "select count(*) from reactions where post_id==? and dislike==TRUE;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.NumberOfDislikes); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindPostCategories(ctx context.Context, p *models.Post) error {
	query := "select distinct categories.name from categories join post_categories on categories.id==post_categories.category_id join posts on post_categories.post_id==posts.id where posts.id==?;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var category string
	for rows.Next() {
		if err := rows.Scan(&category); err != nil {
			return err
		}
		p.Categories = append(p.Categories, category)
	}
	return nil
}

func (r *PostRepo) FindPostCreator(ctx context.Context, p *models.Post) error {
	query := "select users.username from users where uuid=?"
	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.User_UUID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.User_UUID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindPostbyID(ctx context.Context, post_id string) (models.Post, error) {
	query := "select * from posts where id=?"
	rows, err := r.DBlogger.DB.QueryContext(ctx, query, post_id)
	if err != nil {
		return models.Post{}, err
	}
	defer rows.Close()

	var post models.Post
	var postime time.Time
	for rows.Next() {
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.User_UUID, &post.ImagePath, &postime); err != nil {
			return models.Post{}, err
		}
		post.CreationDate = postime.Format("2006-01-02 15:04:05")

	}

	if err := r.FindPostInfo(ctx, &post); err != nil {
		log.Println(err)
		return post, err
	}
	if err := r.FindCommentsofPost(ctx, &post); err != nil {
		log.Println(err)
		return post, err
	}
	return post, nil
}

func (r *PostRepo) FindCommentsofPost(ctx context.Context, post *models.Post) error {
	q := "select * from comments where post_id==?;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, q, post.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var comment models.Comment
	for rows.Next() {
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.User_UUID, &comment.Post_id, &comment.CreationDate); err != nil {
			return err
		}
		if err = r.FindCommentInfo(ctx, &comment); err != nil {
			return err
		}

		post.Comments = append(post.Comments, comment)
	}
	return nil
}

func (r *PostRepo) FindCommentInfo(ctx context.Context, comment *models.Comment) error {
	if err := r.FindCommentCreator(ctx, comment); err != nil {
		return err
	}
	if err := r.FindCommentLikes(ctx, comment); err != nil {
		return err
	}
	if err := r.FindCommentDislikes(ctx, comment); err != nil {
		return err
	}
	if err := r.FindIfUserReactedOnComment(ctx, comment); err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) FindIfUserReactedOnComment(ctx context.Context, p *models.Comment) error {
	user := ctx.Value("user").(models.User)
	if user.Username == "" {
		return nil
	}
	query := "select like,dislike from reactions where user_uuid==? and comment_id==?"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, user.UUID, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var liked bool
	var disliked bool

	counter := 0
	for rows.Next() {
		if counter >= 1 {
			return errors.New("Unexpected query result: User has more than one reactions to a post")
		}
		if err := rows.Scan(&liked, &disliked); err != nil {
			return err
		}
		counter++
	}

	p.Liked = liked
	p.Disliked = disliked
	return nil
}

func (r *PostRepo) FindCommentCreator(ctx context.Context, p *models.Comment) error {
	query := "select users.username from users where uuid=?"
	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.User_UUID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.User_UUID); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindCommentLikes(ctx context.Context, p *models.Comment) error {
	query := "select count(*) from reactions where comment_id==? and like==TRUE;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.NumberOfLikes); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) FindCommentDislikes(ctx context.Context, p *models.Comment) error {
	query := "select count(*) from reactions where comment_id==? and dislike==TRUE;"

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&p.NumberOfDislikes); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) ReturnUserLikeStatusOnPost(ctx context.Context, post_id string, user models.User, onwhat string) (likestatus string, err error) {
	var query string
	switch onwhat {
	case "post":
		query = "select like from reactions where post_id==? and user_uuid=?;"
	case "comment":
		query = "select like from reactions where comment_id==? and user_uuid=?;"
	}

	rows, err := r.DBlogger.DB.QueryContext(ctx, query, post_id, user.UUID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var liked bool
	counter := 0
	for rows.Next() {
		if err := rows.Scan(&liked); err != nil {
			return "", err
		}
		counter++
	}
	switch {
	case counter > 1:
		return "", errors.New("more than one reaction of user")
	case counter == 1:
		if liked == true {
			return "like", nil
		} else {
			return "dislike", nil
		}
	default:
		return "no reaction", nil
	}

}
func (r *PostRepo) LikeButton(ctx context.Context, post_id string, user models.User, onwhat string) error {
	likestatus, err := r.ReturnUserLikeStatusOnPost(ctx, post_id, user, onwhat)
	if err != nil {
		return err
	}

	switch likestatus {
	case "like":
		if err := r.RemoveReactionFromPost(ctx, post_id, user, "like", onwhat); err != nil {
			return err
		}
	case "dislike":
		if err := r.RemoveReactionFromPost(ctx, post_id, user, "dislike", onwhat); err != nil {
			return err
		}
		if err := r.AddReactionToPost(ctx, post_id, user, "like", onwhat); err != nil {
			return err
		}
	case "no reaction":
		if err := r.AddReactionToPost(ctx, post_id, user, "like", onwhat); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) RemoveReactionFromPost(ctx context.Context, post_id string, user models.User, typeofreaction string, onwhat string) error {
	var query string
	switch onwhat {
	case "post":
		query = "delete from reactions where post_id==? and user_uuid=?;"
	case "comment":
		query = "delete from reactions where comment_id==? and user_uuid=?;"
	}
	_, err := r.DBlogger.LogExecContext(ctx, query, post_id, user.UUID)
	if err != nil {
		return err
	}
	switch typeofreaction {
	case "like":
		if err := r.RemoveNotification(ctx, user.UUID, fmt.Sprintf("liked your %v", onwhat), post_id, onwhat); err != nil {
			return err
		}
		err = r.RemoveActivity(ctx, user.UUID, "liked", post_id, onwhat)
		if err != nil {
			return err
		}
	case "dislike":
		if err := r.RemoveNotification(ctx, user.UUID, fmt.Sprintf("disliked your %v", onwhat), post_id, onwhat); err != nil {
			return err
		}
		err = r.RemoveActivity(ctx, user.UUID, "disliked", post_id, onwhat)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) AddReactionToPost(ctx context.Context, post_id string, user models.User, reactiontype string, onwhat string) error {

	var query string
	switch onwhat {
	case "post":
		query = "insert into reactions (user_uuid,post_id,like,dislike) values (?,?,?,?);"
	case "comment":
		query = "insert into reactions (user_uuid,comment_id,like,dislike) values (?,?,?,?);"
	}
	switch reactiontype {
	case "like":
		_, err := r.DBlogger.LogExecContext(ctx, query, user.UUID, post_id, true, false)
		if err != nil {
			return err
		}
		err = r.AddNotification(ctx, user.UUID, fmt.Sprintf("liked your %v", onwhat), post_id, onwhat)
		if err != nil {
			return err
		}

		err = r.AddActivity(ctx, user.UUID, "liked", post_id, onwhat)
		if err != nil {
			return err
		}

		return nil
	case "dislike":
		_, err := r.DBlogger.LogExecContext(ctx, query, user.UUID, post_id, false, true)
		if err != nil {
			return err
		}
		err = r.AddNotification(ctx, user.UUID, fmt.Sprintf("disliked your %v", onwhat), post_id, onwhat)
		if err != nil {
			return err
		}

		err = r.AddActivity(ctx, user.UUID, "disliked", post_id, onwhat)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("no valid reaction given to add")
}

func (r *PostRepo) DislikeButton(ctx context.Context, post_id string, user models.User, onwhat string) error {
	likestatus, err := r.ReturnUserLikeStatusOnPost(ctx, post_id, user, onwhat)
	if err != nil {
		return err
	}

	switch likestatus {
	case "dislike":
		if err := r.RemoveReactionFromPost(ctx, post_id, user, "dislike", onwhat); err != nil {
			return err
		}
	case "like":
		if err := r.RemoveReactionFromPost(ctx, post_id, user, "like", onwhat); err != nil {
			return err
		}
		if err := r.AddReactionToPost(ctx, post_id, user, "dislike", onwhat); err != nil {
			return err
		}
	case "no reaction":
		if err := r.AddReactionToPost(ctx, post_id, user, "dislike", onwhat); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) AddActivity(ctx context.Context, user_uuid string, action string, post_id string, onwhat string) error {
	query := fmt.Sprintf("insert into activities(user_uuid,action,%v_id,created_at) values(?,?,?,datetime(?, 'unixepoch'));", onwhat)
	_, err := r.DBlogger.LogExecContext(ctx, query, user_uuid, action, post_id, time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

// to do
func (r *PostRepo) RemoveActivity(ctx context.Context, user_uuid string, action string, post_id string, onwhat string) error {
	query := fmt.Sprintf("delete from activities where user_uuid==? and action==? and %v_id==?", onwhat)
	_, err := r.DBlogger.LogExecContext(ctx, query, user_uuid, action, post_id)
	if err != nil {
		return err
	}
	return nil
}

// to do
func (r *PostRepo) RemovePost(ctx context.Context, user_uuid string, post_id string) error {
	if err := r.RemoveImageFromServerIfRemovedFromPost(ctx, post_id, user_uuid); err != nil {
		return err
	}
	query := "delete from posts where id==? and user_uuid==?"
	_, err := r.DBlogger.LogExecContext(ctx, query, post_id, user_uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) RemoveComment(ctx context.Context, user_uuid string, comment_id string) error {
	query := "delete from comments where id==? and user_uuid==?"
	_, err := r.DBlogger.LogExecContext(ctx, query, comment_id, user_uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) EditComment(ctx context.Context, content string, comment_id string, user_uuid string, createdAt int64) error {
	query := "update comments set content=?, created_at=? where id==? and user_uuid==?"
	_, err := r.DBlogger.LogExecContext(ctx, query, content, createdAt, comment_id, user_uuid)
	if err != nil {
		return err
	}
	err = r.AddActivity(ctx, user_uuid, "editted", comment_id, "comment")
	if err != nil {
		return err
	}

	return err
}

func (r *PostRepo) EditPost(ctx context.Context, postinfo *models.PostInfo, post_id string, user_uuid string, createdAt int64) error {
	switch postinfo.ImagePath {
	case "null":
		query := "update posts set title=?,content=?, created_at=? where id==? and user_uuid==?"
		if postinfo.RemoveImg == "1" {
			query = "update posts set title=?,content=?,image='null', created_at=? where id==? and user_uuid==?"
			err := r.RemoveImageFromServerIfRemovedFromPost(ctx, post_id, user_uuid)
			if err != nil {
				return err
			}
		}
		_, err := r.DBlogger.LogExecContext(ctx, query, postinfo.Title, postinfo.Content, createdAt, post_id, user_uuid)
		if err != nil {
			return err
		}
	default:
		err := r.RemoveImageFromServerIfRemovedFromPost(ctx, post_id, user_uuid)
		if err != nil {
			return err
		}
		query := "update posts set title=?,content=?,image=?, created_at=? where id==? and user_uuid==?"
		_, err = r.DBlogger.LogExecContext(ctx, query, postinfo.Title, postinfo.Content, postinfo.ImagePath, createdAt, post_id, user_uuid)
		if err != nil {
			return err
		}
	}

	err := r.RemovePostFromOldCategories(ctx, post_id)
	if err != nil {
		return err
	}
	err = r.InsertPostToPost_CategoriesTable(ctx, post_id, postinfo.Categories)
	if err != nil {
		return err
	}

	err = r.AddActivity(ctx, user_uuid, "editted", post_id, "post")
	if err != nil {
		return err
	}

	return err
}

func (r *PostRepo) RemovePostFromOldCategories(ctx context.Context, post_id string) error {

	query := "delete from post_categories where post_id==?"

	_, err := r.DBlogger.LogExecContext(ctx, query, post_id)
	if err != nil {
		return err
	}

	return nil
}
func (r *PostRepo) RemoveImageFromServerIfRemovedFromPost(ctx context.Context, post_id string, user_uuid string) error {
	q := "select image from posts where id==? and user_uuid==?"

	rows, err := r.DBlogger.DB.QueryContext(ctx, q, post_id, user_uuid)
	if err != nil {
		return err
	}
	defer rows.Close()

	var image string

	for rows.Next() {
		if err := rows.Scan(&image); err != nil {
			return err
		}
	}

	if err := utils.RemoveImage(image); err != nil {
		return err
	}

	return nil
}
