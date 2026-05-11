package views

import (
	"forum/src/models"
)

func UserRegister(data models.ResponseStruct4ViewsIface) {
	data.SetView("user_register_view").WriteResponse()
}
