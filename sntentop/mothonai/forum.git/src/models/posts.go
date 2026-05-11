package models

import (
	"errors"
	"forum/src/utils"
)

type Posts []Post

func GetAllPosts() (Posts, error) {
	rows, err := db.Query(`SELECT id, title, body, timestamp, image_path FROM posts`)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	var posts Posts
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts, &post.ImagePath)
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

func GetPostsByCategoryId(id int64) (Posts, error) {
	var posts Posts
	rows, err := db.Query(`
	SELECT posts.id, posts.title, posts.body, posts.timestamp, posts.image_path
	FROM posts
	JOIN posts_categories pc ON posts.id = pc.post_id
	JOIN categories ON pc.category_id = categories.id
	WHERE pc.category_id = ?`, id)
	if err != nil {
		err = errors.Join(utils.GetFunctionName(), err)
		return Posts{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		var ts string
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &ts, &post.ImagePath)
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
