package controllers

import (
	"fmt"
	"forum/src/models"
	"forum/src/views"
	"forum/src/utils"
	"net/http"
	"strings"
)

func parsePostId(data models.ResponseStruct) (int64, error) {
	postIdStr := data.Request.FormValue("post-id")
	if len(postIdStr) == 0 {
		postIdStr = strings.TrimPrefix(data.Request.RequestURI, "/post/edit/")
	}
	if len(postIdStr) == 0 {
		return 0, models.ErrorPostEmptyId
	}
	postId, err := utils.StringToInt64(postIdStr)
	if err != nil {
		return 0, models.ErrorInvalidPostId
	}
	return postId, nil
}

func markSelectedCategories(categories models.Categories, selected models.Categories) models.Categories {
	selectedIDs := make(map[int64]bool, len(selected))
	for _, category := range selected {
		selectedIDs[category.Id] = true
	}
	for i := range categories {
		categories[i].Selected = selectedIDs[categories[i].Id]
	}
	return categories
}

func getPostDataById(data *models.ResponseStruct) error {
	var err error
	post := &data.Posts[0]
	err = post.GetById()
	if err != nil {
		if err == models.ErrorNoRows {
			err = models.ErrorContentNotFound
			return err
		} else {
			return err
		}
	}
	comments, err := post.GetComments()
	if err != nil {
		if err == models.ErrorNoRows {
			err = models.ErrorContentNotFound
			return err
		} else {
			return err
		}
	}
	for i := range comments {
		err = comments[i].GetReactions()
		if err != nil {
			return err
		}
		err = comments[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			return err
		}
	}
	post.Comments = comments
	categories, err := models.GetCategoriesByPostId(post.Id)
	if err != nil {
		return err
	}
	post.Categories = categories
	err = post.GetReactions()
	if err != nil {
		return err
	}
	err = post.GetReactionsByUserId(data.User.Id)
	if err != nil {
		return err
	}
	return nil
}


func validateViewPostByIdRequest(data models.ResponseStruct) (models.Post, error) {
	var post models.Post
	var err error
	id, ok := strings.CutPrefix(data.Request.RequestURI, "/post/view/")
	if !ok || len(id) == 0 {
		return post, models.ErrorPostEmptyId
	}
	post.Id, err = utils.StringToInt64(id)
	if err != nil {
		return post, models.ErrorInvalidPostId
	}
	return post, nil
}

func showPost(data models.ResponseStruct) {
	var err error
	var post models.Post
	post, err = validateViewPostByIdRequest(data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.Posts = models.Posts{post}
	err = getPostDataById(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.User.MarkAsReadPost(post)
	views.PostView(&data)
}

func showPosts(data models.ResponseStruct) {
	posts, err := models.GetAllPosts()
	if err != nil {
		data.Error.Consume(err)
		(&models.Error{}).Consume(err).RespondError(data.Response, data.User)
		return
	}
	for i := range posts {
		err = posts[i].GetReactions()
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = posts[i].GetReactionsByUserId(data.User.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	data.SetPosts(posts)
	views.PostsView(&data)
}

func createPost(data models.ResponseStruct) {
	post, err := parseCreatePostRequest(data)
	if err != nil {
		data.Error.Consume(err)
		views.PostCreate(&data)
		return
	}
	postId, err := post.Add()
	if err != nil {
		data.Error.Consume(err)
		views.PostCreate(&data)
		return
	}
	redirectURL := fmt.Sprintf("/post/view/%d", postId)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func parseCreatePostRequest(data models.ResponseStruct) (models.Post, error) {
	title := data.Request.FormValue("title")
	body := data.Request.FormValue("body")
	categories, err := models.GetAllCategories()
	if err != nil {
		return models.Post{}, err
	}
	var PostCategories models.Categories
	for _, category := range categories {
		cc := fmt.Sprintf("category-%d", category.Id)
		if data.Request.Form.Has(cc) && data.Request.Form.Get(cc) == "on" {
			PostCategories = append(PostCategories, category)
		}
	}
	imagePath := ""
	imageFile, _, err := data.Request.FormFile("image")
	if err == nil {
		defer imageFile.Close()
		imagePath, err = models.SaveImage(imageFile)
		if err != nil {
			return models.Post{}, err
		}
	}
	return models.Post{
		Title:      title,
		Body:       body,
		UserId:     data.User.Id,
		Categories: PostCategories,
		ImagePath: imagePath,
	}, nil
}

func handlePostCreate(data models.ResponseStruct) {
	switch {
	case strings.Compare(data.Request.RequestURI, "/post/create") == 0:
		switch data.Request.Method {
		case http.MethodGet:
			categories, err := models.GetAllCategories()
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
			data.SetCategories(categories)
			views.PostCreate(&data)
			return
		case http.MethodPost:
			err := data.Request.ParseMultipartForm(models.MaxImageSize)
			if err != nil {
				(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
				return
			}
			createPost(data)
		default:
			(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
		}
	}
}

func handlePostReaction(data models.ResponseStruct) {
	var post models.Post
	var err error
	post.Id, err = parsePostId(data)
	if err != nil {
		err = models.ErrorInvalidPostId
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if data.Request.FormValue("action") == "like" {
		err = data.User.LikePost(post.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = post.CreateReactionNotification(data.User.Id, "like")
		if err != nil {
		(&models.Error{}).Consume(err).LogError()
		}
	}
	if data.Request.FormValue("action") == "dislike" {
		err = data.User.DislikePost(post.Id)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		err = post.CreateReactionNotification(data.User.Id, "dislike")
		if err != nil {
		(&models.Error{}).Consume(err).LogError()
		}		
	}
	redirectURL := fmt.Sprintf("/post/view/%d", post.Id)
	http.Redirect(data.Response, data.Request, redirectURL, http.StatusSeeOther)
}

func handlePostDelete(data models.ResponseStruct) {
	var err error
	var post models.Post
	post.Id, err = parsePostId(data)
	if err != nil {
		(&models.Error{}).Consume(models.ErrorInvalidPostId).LogAndRespondError(data.Response, data.User)
		return
	}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	if post.UserId != data.User.Id {
		(&models.Error{}).Consume(models.ErrorPostPermissionDenied).LogAndRespondError(data.Response, data.User)
		return
	}
	err = post.Delete()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	http.Redirect(data.Response, data.Request, "/posts", http.StatusSeeOther)
}

func handlePostEdit(data models.ResponseStruct){
	var err error
	var post models.Post
	post.Id, err = parsePostId(data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = post.GetById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.Posts = models.Posts{post}
	data.Categories, err = models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = verifyUserPostAssociation(&data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	switch data.Request.Method {
	case http.MethodGet:
		err = showEditPost(&data)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		views.PostCreate(&data)
	case http.MethodPost:
		err = updatePost(&data)
		if err != nil {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
		http.Redirect(data.Response, data.Request, fmt.Sprintf("/post/view/%d", post.Id), http.StatusSeeOther)
	default:
		(&models.Error{}).Consume(models.ErrorMethodNotAllowed).LogAndRespondError(data.Response, data.User)
	}
}

func verifyUserPostAssociation(data *models.ResponseStruct) error {
	post := &data.Posts[0]
	// Check your priviledge
	if post.UserId != data.User.Id {
		return models.ErrorCommentPermissionDenied
	}
	return nil
}

func showEditPost(data *models.ResponseStruct) error {
	var err error
	post := data.Posts[0]
	categories := data.Categories
	post.Categories, err = models.GetCategoriesByPostId(post.Id)
	if err != nil {
		return err
	}
	data.Categories = markSelectedCategories(categories, post.Categories)
	data.EditPost = true
	return nil
}

func updatePost(data *models.ResponseStruct) error {
	post := data.Posts[0]
	err := data.Request.ParseMultipartForm(models.MaxImageSize)
	if err != nil {
		return err
	}
	updatedPost, err := parseCreatePostRequest(*data)
	if err != nil {
		return err
	}
	updatedPost.Id = post.Id
	if updatedPost.ImagePath == "" {
		updatedPost.ImagePath = post.ImagePath
	}
	post = updatedPost
	err = post.Update()
	if err != nil {
		return err
	}
	data.Posts = models.Posts{post}
	data.EditPost = false
	return nil
}
