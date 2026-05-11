package models

func (user *User) DislikePost(postId int64) error {
	dislikeId, err := CheckIfUserDislikedPost(user.Id, postId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return RemoveDislikeFromPost(dislikeId)
	}
	existingLikeId, err := CheckIfUserLikedPost(user.Id, postId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = RemoveLikeFromPost(user.Id, postId)
		if err != nil {
			return err
		}
	}
	return AddDislikeToPost(user.Id, postId)
}

func (user *User) DislikeComment(commentId int64) error {
	dislikeId, err := CheckIfUserDislikedComment(user.Id, commentId)
	if err != nil {
		return err
	}
	if dislikeId != 0 {
		return RemoveReaction(dislikeId)
	}
	existingLikeId, err := CheckIfUserLikedComment(user.Id, commentId)
	if err != nil {
		return err
	}
	if existingLikeId != 0 {
		err = RemoveReaction(existingLikeId)
		if err != nil {
			return err
		}
	}
	return AddDislikeToComment(user.Id, commentId)
}
