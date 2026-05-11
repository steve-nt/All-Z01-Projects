package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum-authentication/internal/backend/models"
	front_end_repo "forum-authentication/internal/frontend/repositories"
	"forum-authentication/internal/utils"
	"io"
	"log"
	"net/http"
	"strings"
)

type GetHandlers struct {
	FrontEndService *front_end_repo.FrontEndRepo
}

func (h *GetHandlers) GetLanding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	urltocallBackEnd := "/"
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	r.ParseForm()
	category, hascategories := r.Form["categories"]
	if hascategories {
		urltocallBackEnd += "?category=" + strings.Join(category, "_")
	}

	resp, err := h.FrontEndService.Do(r, w, "GET", urltocallBackEnd, nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error reading response:%v", err), 500)
		return
	}

	var home models.PostResponse
	if err := json.Unmarshal(data, &home); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("JSON parse error: %v, received data: %s", err, string(data)), 500)
		return
	}

	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "landing.page.html", map[string]models.PostResponse{"Homeresponse": home}); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}
func (h *GetHandlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	urltocallBackEnd := "/profile"
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	resp, err := h.FrontEndService.Do(r, w, "GET", urltocallBackEnd, nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error reading response:%v", err), 500)
		return
	}

	var profile models.ProfileResponse
	if err := json.Unmarshal(data, &profile); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("JSON parse error: %v, received data: %s", err, string(data)), 500)
		return
	}

	utils.ReversePosts(&profile.LikedPosts)
	utils.ReversePosts(&profile.CreatedPosts)
	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "profile.page.html", map[string]models.ProfileResponse{"profileresponse": profile}); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}
func (h *GetHandlers) GetHome(w http.ResponseWriter, r *http.Request) {
	urltocallBackEnd := "/"
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	r.ParseForm()
	category, hascategories := r.Form["categories"]
	if hascategories {
		urltocallBackEnd += "?category=" + strings.Join(category, "_")
	}

	resp, err := h.FrontEndService.Do(r, w, "GET", urltocallBackEnd, nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error reading response:%v", err), 500)
		return
	}

	var home models.PostResponse
	if err := json.Unmarshal(data, &home); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("JSON parse error: %v, received data: %s", err, string(data)), 500)
		return
	}

	utils.ReversePosts(&home.Posts)
	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "home.page.html", map[string]models.PostResponse{"Homeresponse": home}); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}

func (h *GetHandlers) GetPostbyID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	urltocallBackEnd := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "GET", urltocallBackEnd, nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		http.Redirect(w, r, "/posts", http.StatusFound)
		return
	}

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error reading response:%v", err), 500)
		return
	}

	var postbyid models.PostByIdResponse
	if err := json.Unmarshal(data, &postbyid); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("JSON parse error: %v, received data: %s", err, string(data)), 500)
		return
	}

	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "postbyid.page.html", map[string]models.PostByIdResponse{"Response": postbyid}); err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}
func (h *GetHandlers) GetCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "createpost.page.html", nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}

func (h *GetHandlers) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "login.page.html", nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, err, 500)
	}
}

func (h *GetHandlers) GetSignUpPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "signup.page.html", nil); err != nil {
		log.Println("Error rendering signup page:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *GetHandlers) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	success := r.URL.Query().Get("success")
	msg := r.URL.Query().Get("msg")

	data := struct {
		Success bool
		Message string
	}{
		Success: success == "true",
		Message: msg,
	}

	// χρησιμοποιούμε το FrontEndService.Tmpl (όπως όλα τα άλλα pages)
	if err := h.FrontEndService.Tmpl.ExecuteTemplate(w, "verify.page.html", data); err != nil {
		log.Println("Error rendering verify page:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *GetHandlers) GetSocialSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Missing provider"), 400)
		return
	}

	log.Printf("Social signup initiated for provider: %s", provider)

	backendBase := "http://localhost:8080"
	target := fmt.Sprintf("%s/auth/signup?provider=%s", backendBase, provider)
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *GetHandlers) GetSocialSignupCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	backendURL := "/auth/callback?" + r.URL.RawQuery
	resp, err := h.FrontEndService.Do(r, w, "GET", backendURL, nil)
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		http.Redirect(w, r, "/signup?error=social_auth_failed", http.StatusFound)
	}

}

// func (h *GetHandlers) GetUserHome(w http.ResponseWriter, r *http.Request) {
// }
//
// // GetAuthStatus checks if user is authenticated by calling backend
// func (h *GetHandlers) GetAuthStatus(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	resp, err := h.FrontEndService.Do(r, w, "GET", "/auth-status", nil)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(`{"isAuthenticated": false, "error": "backend unavailable"}`))
// 		return
// 	}
// 	defer resp.Body.Close()
//
// 	// Forward the JSON response from backend
// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(`{"isAuthenticated": false, "error": "read error"}`))
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(resp.StatusCode)
// 	w.Write(data)
// }
