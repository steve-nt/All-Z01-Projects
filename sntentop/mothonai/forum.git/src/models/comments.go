package models

import (
	"errors"
	"forum/src/utils"
)

type Comments []Comment

func GetCommentsByUserId(id int64) (Comments, error) {
	var comments Comments
	rows, err := db.Query(`
	SELECT id, post_id, body, timestamp, user_id
	FROM comments
	WHERE user_id = ?`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Comments{}, err
	}
	for rows.Next() {
		var comment Comment
		var ts string
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Body, &ts, &comment.UserId)
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
