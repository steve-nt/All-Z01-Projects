package controllers

import (
	"forum/src/models"
	"forum/src/views"
)

func Index(data models.ResponseStruct) {
	categories, err := models.GetAllCategories()
	if err != nil {
		(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
		return
	}
	data.SetCategories(categories)
	views.Index(&data)
}
