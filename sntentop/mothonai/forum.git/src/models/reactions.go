package models

import (
	"errors"
	"forum/src/utils"
)

type Reactions []Reaction

func GetPostLikesByUserId(id int64) (Reactions, error) {
	var reactions Reactions
	rows, err := db.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1 AND post_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Reactions{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction Reaction
		var ts string
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetPostDislikesByUserId(id int64) (Reactions, error) {
	var reactions Reactions
	rows, err := db.Query(`
	SELECT id, post_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=2 AND post_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Reactions{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction Reaction
		var ts string
		err = rows.Scan(&reaction.Id, &reaction.PostId, &reaction.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentLikesByUserId(id int64) (Reactions, error) {
	var reactions Reactions
	rows, err := db.Query(`
	SELECT id, comment_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=1 AND comment_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Reactions{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction Reaction
		var ts string
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func GetCommentDisikesByUserId(id int64) (Reactions, error) {
	var reactions Reactions
	rows, err := db.Query(`
	SELECT id, comment_id, user_id, timestamp
	FROM reactions
	WHERE user_id = ? AND value=2 AND comment_id IS NOT NULL
	`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Reactions{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction Reaction
		var ts string
		err = rows.Scan(&reaction.Id, &reaction.CommentId, &reaction.UserId, &ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		t, err := utils.ConvertStringToTime(ts)
		if err != nil {
			err = errors.Join(utils.GetFunctionName(), err)
			return Reactions{}, err
		}
		reaction.TimestampString = utils.ConvertTimeToString(t)
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
