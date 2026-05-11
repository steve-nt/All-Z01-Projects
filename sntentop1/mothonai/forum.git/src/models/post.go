package models

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

type Post struct {
	Id              int64
	Title           string
	Body            string
	ImagePath       string
	UserId          int64
	User            User
	Timestamp       int64
	TimestampString string
	Likes           int
	Liked           bool
	Dislikes        int
	Disliked        bool
	Category        Category
	Categories      Categories
	Comments        Comments
}

func (p *Post) ValidatePost() error {
	if len(p.Title) == 0 {
		return ErrorPostTitleEmpty
	}
	if len(p.Body) == 0 {
		return ErrorPostBodyEmpty
	}
	if p.Categories.IsEmpty() {
		return ErrorPostHasNoCategory
	}
	return nil
}

// Adds a Post in the database. Returns its id or error
func (p *Post) Add() (int64, error) {
	err := p.ValidatePost()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	stmt, err := db.Prepare("INSERT INTO posts (title, body, image_path, user_id, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	res, err := stmt.Exec(p.Title, p.Body, p.ImagePath, p.UserId, utils.GetCurrentTimestamp())
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	postId, err := res.LastInsertId()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	p.Id = postId
	for _, category := range p.Categories {
		err = p.AddCategory(category)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return 0, err
		}
	}
	return postId, nil
}

func (p *Post) AddCategory(category Category) error {
	stmt, err := db.Prepare("INSERT INTO posts_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = stmt.Exec((*p).Id, category.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func (p *Post) GetById() error {
	var ts string
	err := db.QueryRow(`SELECT title, body, image_path, timestamp, user_id FROM posts WHERE id = ?`, p.Id).Scan(&p.Title, &p.Body, &p.ImagePath, &ts, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrorNoRows
			return err
		} else {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
	}
	t, err := utils.ConvertStringToTime(ts)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	p.TimestampString = utils.ConvertTimeToString(t)
	p.User, err = getUserById(p.UserId)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}

func (p *Post) GetReactions() error {
	var err error
	(*p).Likes, err = getLikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	(*p).Dislikes, err = getDislikesCountByPostId((*p).Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetReactionsByUserId(user_id int64) error {
	var err error
	(*p).Liked, err = HasUserLikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	(*p).Disliked, err = HasUserDislikedPost(user_id, (*p).Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetComments() (Comments, error) {
	rows, err := db.Query(`
	SELECT
	c.id,
	c.post_id,
	c.user_id,
	c.body,
	c.timestamp,
	u.username
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.timestamp ASC`, p.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrorNoRows
			return Comments{}, err
		} else {
			err = errors.Join(utils.GetFunctionName(), err)
			return Comments{}, err
		}
	}
	defer rows.Close()
	var comments Comments
	for rows.Next() {
		var comment Comment
		var ts string
		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.UserId,
			&comment.Body,
			&ts,
			&comment.Username,
		)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Comments{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Comments{}, err
		}
		comment.TimestampString = utils.ConvertTimeToString(t)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (p *Post) Delete() error {
	tx, err := db.Begin()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM posts_categories WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE post_id = ?", p.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return tx.Commit()
}

func (p *Post) Update() error {
	err := p.ValidatePost()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = db.Exec("UPDATE posts SET title = ?, body = ?, image_path = ? WHERE id = ?", p.Title, p.Body, p.ImagePath, p.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = db.Exec("DELETE FROM posts_categories WHERE post_id = ?", p.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	for _, category := range p.Categories {
		err = p.AddCategory(category)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return err
		}
	}
	return nil
}
