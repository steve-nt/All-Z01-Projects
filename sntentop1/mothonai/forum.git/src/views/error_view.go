package views

import (
	"forum/src/models"
)

func ErrorView(data models.ResponseStruct4ViewsIface) {
	data.SetView("error_view").WriteResponse()
}
