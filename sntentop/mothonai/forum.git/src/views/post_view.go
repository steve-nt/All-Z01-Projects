package views

import "forum/src/models"

func PostView(data models.ResponseStruct4ViewsIface) {
	data.SetView("post_view").WriteResponse()
}
