package views

import "forum/src/models"

func UserActivity(data models.ResponseStruct4ViewsIface) {
	data.SetView("user_activities").WriteResponse()
}
