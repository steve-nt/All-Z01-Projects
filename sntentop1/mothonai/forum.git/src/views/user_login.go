package views

import (
	"forum/src/models"
)

func UserLogin(data models.ResponseStruct4ViewsIface) {
	data.SetView("user_login_view").WriteResponse()
}
