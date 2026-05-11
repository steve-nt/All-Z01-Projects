package controllers

import (
	"forum/src/models"
	"forum/src/utils"
	"forum/src/views"
	"strings"
)

func parseCategoryId(data models.ResponseStruct) (int64, error) {
	id, ok := strings.CutPrefix(data.Request.RequestURI, "/category/view/")
	if !ok || len(id) == 0 {
		return 0, models.ErrorCategoryEmptyId
	}
	categoryId, err := utils.StringToInt64(id)
	if err != nil {
		return 0, models.ErrorInvalidCategoryId
	}
	return categoryId, nil
}

func showCategories(data models.ResponseStruct) {
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.SetCategories(categories)
	views.Categories(&data)
}

func showCategory(data models.ResponseStruct) {
	var category models.Category
	var err error
	category.Id, err = parseCategoryId(data)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	err = category.GetCategoryById()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.Categories = models.Categories{category}
	posts, err := models.GetPostsByCategoryId(category.Id)
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
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
	data.Posts = posts
	data.SetPosts(posts)
	views.Category(&data)
}
