package models

import (
	"errors"
	"forum/src/utils"
)

type Comment struct {
	Id              int64
	PostId          int64
	UserId          int64
	Body            string
	Timestamp       int64
	TimestampString string
	Likes           int64
	Liked           bool
	Dislikes        int64
	Disliked        bool
	Username        string
}

func (c *Comment) ValidateComment() error {
	if len(c.Body) == 0 {
		return ErrorCommentEmpty
	}
	if len(c.Body) > 1000 {
		return ErrorCommentTooLong
	}
	return nil
}

func (c *Comment) Add() (int64, error) {
	if err := c.ValidateComment(); err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	res, err := db.Exec(
		"INSERT INTO comments (post_id, user_id, body, timestamp) VALUES (?, ?, ?, ?)",
		c.PostId,
		c.UserId,
		c.Body,
		utils.GetCurrentTimestamp(),
	)
	commentId, err := res.LastInsertId()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return commentId, nil
}

func (c *Comment) GetReactions() error {
	var err error
	(*c).Likes, err = getLikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	(*c).Dislikes, err = getDislikesCountByCommentId((*c).Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) GetReactionsByUserId(user_id int64) error {
	var err error
	(*c).Liked, err = HasUserLikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	(*c).Disliked, err = HasUserDislikedComment(user_id, (*c).Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) GetCommentById() error {
	var ts string
	err := db.QueryRow(
		`SELECT id, post_id, user_id, body, timestamp
		FROM comments
		WHERE id = ?`, c.Id).Scan(&c.Id, &c.PostId, &c.UserId, &c.Body, &ts)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	t, err := utils.ConvertStringToTime(ts)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	c.TimestampString = utils.ConvertTimeToString(t)
	return nil
}

func (c *Comment) Delete() error {
	tx, err := db.Begin()
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM reactions WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM comments WHERE id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE comment_id = ?", c.Id)
	if err != nil {
		tx.Rollback()
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return tx.Commit()
}

func (c *Comment) Update() error {
	if err := c.ValidateComment(); err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	_, err := db.Exec("UPDATE comments SET body = ? WHERE id = ?", c.Body, c.Id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return err
	}
	return nil
}
