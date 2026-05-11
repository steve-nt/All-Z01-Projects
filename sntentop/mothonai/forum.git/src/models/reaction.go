package models

type Reaction struct {
	Id int64
	UserId int64
	PostId int64
	CommentId int64
	Timestamp int64
	TimestampString string
	Post Post
	Comment Comment
}

func RemoveReaction(reactionId int64) error {
	_, err := db.Exec(`
		DELETE FROM reactions
		WHERE id = ?
		`, reactionId)
	return err
}
