package views

import "forum/src/models"

func PostsView(data models.ResponseStruct4ViewsIface) {
	data.SetView("posts_view").WriteResponse()
}
