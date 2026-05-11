package views

import (
	"forum/src/models"
)

func PostCreate(data models.ResponseStruct4ViewsIface) {
	data.SetView("post_create_view").WriteResponse()
}
