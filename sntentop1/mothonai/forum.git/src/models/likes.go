package models

import (
	"errors"
	utils "forum/src/utils"
)

func getLikesCountByPostId(postId int64) (int, error) {
	var likes int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE post_id = ? AND value = 1
    `, postId).Scan(&likes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return likes, nil
}

func getLikesCountByCommentId(commentId int64) (int64, error) {
	var likes int64
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM reactions
        WHERE comment_id = ? AND value = 1
    `, commentId).Scan(&likes)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return 0, err
	}
	return likes, nil
}
