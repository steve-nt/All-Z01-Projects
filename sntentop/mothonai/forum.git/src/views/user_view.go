package views

import "forum/src/models"

func UserView(data models.ResponseStruct4ViewsIface) {
	data.SetView("user_view").WriteResponse()
}
