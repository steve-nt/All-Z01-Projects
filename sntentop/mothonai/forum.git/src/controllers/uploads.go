package controllers

import (
	"forum/src/models"
	"os"
	"errors"
	"path/filepath"
	"strings"
	"net/http"
)

func handleImages(data models.ResponseStruct) {
	if strings.HasSuffix(data.Request.URL.Path, "/") {
		(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
		return
	}
	imgURL := filepath.Base(data.Request.URL.Path)
	_, err := os.Stat("./uploads/images/" + imgURL)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			(&models.Error{}).Consume(models.ErrorNotFound).LogAndRespondError(data.Response, data.User)
			return
		} else {
			(&models.Error{}).Consume(err).LogAndRespondError(data.Response, data.User)
			return
		}
	}
	http.ServeFile(data.Response, data.Request, "./uploads/images/"+imgURL)
}
