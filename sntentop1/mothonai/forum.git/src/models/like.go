package models

import (
	"database/sql"
	"errors"
	"forum/src/utils"
)

func CheckIfUserLikedPost(userId, postId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func CheckIfUserDislikedPost(userId, postId int64) (int64, error) {
	var existingDislikeId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 2
		`, userId, postId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddLikeToPost(userId, postId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, post_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, postId, utils.GetCurrentTimestamp())
	return err
}

func RemoveDislikeFromPost(dislikeId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, dislikeId)
	return err
}

func AddDislikeToPost(userId, postId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, post_id, value, timestamp)
		VALUES (?, ?, 2, ?)
		`, userId, postId, utils.GetCurrentTimestamp())
	return err
}

func RemoveLikeFromPost(userId, postId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND post_id = ? AND value = 1
		`, userId, postId)
	return err
}

func CheckIfUserLikedComment(userId, commentId int64) (int64, error) {
	var existingReactionId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId).Scan(&existingReactionId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingReactionId, nil
}

func CheckIfUserDislikedComment(userId, commentId int64) (int64, error) {
	var existingDislikeId int64
	err := db.QueryRow(`
		SELECT id FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 2
		`, userId, commentId).Scan(&existingDislikeId)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return existingDislikeId, nil
}

func AddLikeToComment(userId, commentId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value, timestamp)
		VALUES (?, ?, 1, ?)
		`, userId, commentId, utils.GetCurrentTimestamp())
	return err
}

func AddDislikeToComment(userId, commentId int64) error {
	_, err := db.Exec(`
		INSERT INTO reactions (user_id, comment_id, value, timestamp)
		VALUES (?, ?, 2, ?)
		`, userId, commentId, utils.GetCurrentTimestamp())
	return err
}

func RemoveLikeFromComment(userId, commentId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE user_id = ? AND comment_id = ? AND value = 1
		`, userId, commentId)
	return err
}

func (user *User) LikeComment(commentId int64) error {
	alreadyLiked, err := HasUserLikedComment(user.Id, commentId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return RemoveLikeFromComment(user.Id, commentId)
	}
	existingDislikeId, err := CheckIfUserDislikedComment(user.Id, commentId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		if err = RemoveReaction(existingDislikeId); err != nil {
			return err
		}
	}
	return AddLikeToComment(user.Id, commentId)
}

func (user *User) LikePost(postId int64) error {
	alreadyLiked, err := HasUserLikedPost(user.Id, postId)
	if err != nil {
		return err
	}
	if alreadyLiked {
		return RemoveLikeFromPost(user.Id, postId)
	}
	existingDislikeId, err := CheckIfUserDislikedPost(user.Id, postId)
	if err != nil {
		return err
	}
	if existingDislikeId != 0 {
		if err = RemoveReaction(existingDislikeId); err != nil {
			return err
		}
	}
	return AddLikeToPost(user.Id, postId)
}
