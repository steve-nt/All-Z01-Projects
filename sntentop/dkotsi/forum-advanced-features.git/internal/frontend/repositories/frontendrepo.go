package front_end_repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum-advanced-features/internal/backend/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type FrontEndRepo struct {
	Client *http.Client
	Tmpl   *template.Template
	Config *models.Config
}

func NewFrontEndRepo(tmpl *template.Template, conf *models.Config) *FrontEndRepo {
	return &FrontEndRepo{
		Client: &http.Client{
			Timeout: time.Duration(conf.Durations.ClientTimeOut) * time.Duration(time.Second),
		},
		Tmpl:   tmpl,
		Config: conf,
	}
}

func (a *FrontEndRepo) api(path string) string {
	return a.Config.Api.Api_base_url + path
}

func (a *FrontEndRepo) GetFormData(r *http.Request) (data []byte, err error, status int) {
	if err := r.ParseMultipartForm(20000000); err != nil {
		return nil, err, 500
	}
	var postFormValues struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
		RemoveImg  string   `json:"remove-img"`
		ImagePath  string   `json:"imagepath"`
	}
	postFormValues.Title = r.Form["title"][0]
	if postFormValues.Title == "" {
		return nil, errors.New("Post To be Created has no Title "), 400
	}
	postFormValues.Content = r.Form["content"][0]
	if postFormValues.Content == "" {
		return nil, errors.New("Post To be Created has no Content"), 400
	}
	postFormValues.Categories = r.Form["create-categories"]
	postFormValues.RemoveImg = r.FormValue("remove-img")
	file, fileheader, err := r.FormFile("image")
	if err != nil {
		postFormValues.ImagePath = "null"
	} else {

		size := fileheader.Size
		if size > 20000000 {
			return nil, errors.New("File Larger than 20mb"), 400
		}
		filetype := fileheader.Filename
		idx := strings.LastIndex(filetype, ".")
		if idx == -1 {
			return nil, errors.New("Invalid filetype given"), 400
		}
		filetype = filetype[idx:]
		filename, _ := uuid.NewV4()
		imagepath := "/postimages/" + filename.String() + filetype
		postFormValues.ImagePath = imagepath
		imagedata, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("could not read file :%v", err), 500
		}

		err = os.WriteFile("../../static"+imagepath, imagedata, 0644)
		if err != nil {
			return nil, err, 500
		}
	}

	data, _ = json.Marshal(postFormValues)
	return data, nil, 0
}
func (a *FrontEndRepo) FrontEndServerErrorwithHTML(w http.ResponseWriter, err error, status int) error {
	log.Println(err)
	w.WriteHeader(status)
	_ = a.Tmpl.ExecuteTemplate(w, "error.page.html", err)
	return nil
}
func (a *FrontEndRepo) ErrorFromBackEndHtml(resp *http.Response, w http.ResponseWriter) error {
	log.Printf("Backend returned error status: %d", resp.StatusCode)
	w.WriteHeader(resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		a.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error reading response:%v", err), 500)
		return err
	}
	var newerror string
	if err := json.Unmarshal(data, &newerror); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		a.Tmpl.ExecuteTemplate(w, "error.page.html", err)
		return err
	}
	log.Println("Backend returned error:", newerror)
	err = a.Tmpl.ExecuteTemplate(w, "error.page.html", newerror)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (a *FrontEndRepo) Do(r *http.Request, w http.ResponseWriter, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, a.api(path), body)
	if err != nil {
		return nil, err
	}

	if v := r.Header.Get("Content-Type"); v != "" {
		req.Header.Set("Content-Type", v)
	}
	if v := r.Header.Get("Accept"); v != "" {
		req.Header.Set("Accept", v)
	}

	for _, c := range r.Cookies() {
		req.AddCookie(c)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	for _, v := range resp.Header.Values("Set-Cookie") {
		w.Header().Add("Set-Cookie", v)
	}

	return resp, nil
}
