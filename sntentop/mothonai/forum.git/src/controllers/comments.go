package controllers

import (
	"fmt"
	"forum/src/models"
	"forum/src/utils"
	"forum/src/views"
	"net/http"
)

func parseCommentId(data models.ResponseStruct) (int64, error) {
	commentIdStr := data.Request.FormValue("comment-id")
	if len(commentIdStr) == 0 {
		return 0, models.ErrorCommentEmptyId
	}
	commentId, err := utils.StringToInt64(commentIdStr)
	if err != nil {
		return 0, models.ErrorInvalidCommentId
	}
	return commentId, nil
}

func handleCommentCreate(data models.ResponseStruct) {
	var comment models.Comment
	var err error
	body := data.Request.FormValue("comment")
	comment.PostId, err = parsePostId(data)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	post := models.Post{Id: comment.PostId}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	comment.Body = body
	comment.UserId = data.User.Id
	comment.Id, err = comment.Add()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = comment.CreateCommentNotification(post)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", post.Id, comment.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentReaction(data models.ResponseStruct) {
	var comment models.Comment
	var err error
	comment.Id, err = parseCommentId(data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = comment.GetCommentById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.FormValue("action") == "like" {
		err = data.User.LikeComment(comment.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = comment.CreateReactionNotification(data.User.Id, "commentLike")
		if err != nil {
			(&models.Error{}).Consume(err).LogError()
		}
	}
	if data.Request.FormValue("action") == "dislike" {
		err = data.User.DislikeComment(comment.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = comment.CreateReactionNotification(data.User.Id, "commentDislike")
		if err != nil {
			(&models.Error{}).Consume(err).LogError()
		}
	}
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", comment.PostId, comment.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentDelete(data models.ResponseStruct) {
	commentId, err := parseCommentId(data)
	if err != nil {
		err = models.ErrorInvalidCommentId
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	comment := models.Comment{Id: commentId}
	err = comment.GetCommentById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if comment.UserId != data.User.Id {
		(&models.Error{}).Consume(models.ErrorCommentPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	err = comment.Delete()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	redirectURL := fmt.Sprintf("/post/view/%d", comment.PostId)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handleCommentEdit(data models.ResponseStruct) {
	var err error
	err = validateFormCommentEdit(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = verifyCommentOwnership(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.FormValue("save-comment") == "1" {
		err = updateCommentFromForm(&data)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	err = showEditComment(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	views.PostView(&data)
}

func validateFormCommentEdit(data *models.ResponseStruct) error {
	var err error
	var post models.Post
	var comment models.Comment
	commentId, err := parseCommentId(*data)
	if err != nil {
		return models.ErrorInvalidCommentId
	}
	post.Id, err = parsePostId(*data)
	if err != nil {
		return models.ErrorInvalidPostId
	}
	comment = models.Comment{Id: commentId}
	post.Comments = models.Comments{comment}
	data.Posts = models.Posts{post}
	return nil
}

func verifyCommentOwnership(data *models.ResponseStruct) error {
	var err error
	comment := &data.Posts[0].Comments[0]
	err = comment.GetCommentById()
	if err != nil {
		return err
	}
	// Check your priviledge
	if comment.UserId != data.User.Id {
		return models.ErrorCommentPermissionDenied
	}
	return nil
}

func showEditComment(data *models.ResponseStruct) error {
	var err error
	comment := data.Posts[0].Comments[0]
	err = getPostDataById(data)
	if err != nil {
		return err
	}
	data.EditCommentId = comment.Id
	return nil
}

func updateCommentFromForm(data *models.ResponseStruct) error {
	var err error
	comment := &data.Posts[0].Comments[0]
	comment.Body = data.Request.FormValue("comment")
	err = comment.Update()
	if err != nil {
		return err
	}
	redirectURL := fmt.Sprintf("/post/view/%d#comment-%d", comment.PostId, comment.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
	return nil
}
