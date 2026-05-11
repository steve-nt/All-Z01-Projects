package views

import "forum/src/models"

func Category(data models.ResponseStruct4ViewsIface) {
	data.SetView("category_view").WriteResponse()
}
