package views

import "forum/src/models"

func Categories(data models.ResponseStruct4ViewsIface) {
	data.SetView("categories_view").WriteResponse()
}
