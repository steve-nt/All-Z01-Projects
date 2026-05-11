package views

import (
	"forum/src/models"
)

func UserLogout(data models.ResponseStruct4ViewsIface) {
	data.SetView("user_logout_view").WriteResponse()
}
