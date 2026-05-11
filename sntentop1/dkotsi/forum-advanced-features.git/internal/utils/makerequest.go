package utils

import (
	"bytes"
	"log"
	"net/http"
)

const (
	cookieName = "example"
)

func MakeRequestwithCookies(r *http.Request, w http.ResponseWriter, method string, url string, body *bytes.Buffer) (response *http.Response, cookie *http.Cookie) {

	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println(err)
	}

	client := http.DefaultClient

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
	}

	if cookie != nil {
		request.AddCookie(cookie)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	allcookies := resp.Cookies()
	for _, c := range allcookies {
		if c.Name == cookieName {
			http.SetCookie(w, c)
			return resp, c
		}
	}

	return resp, cookie
}
