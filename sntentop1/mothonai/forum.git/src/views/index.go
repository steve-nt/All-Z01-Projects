package views

import (
	"forum/src/models"
)

func Index(data models.ResponseStruct4ViewsIface) {
	data.SetView("index_view").WriteResponse()
}
