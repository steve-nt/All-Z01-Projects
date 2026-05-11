package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum-image-upload/internal/frontend/repositories"
	"forum-image-upload/internal/utils"
	"io"
	"log"
	"net/http"
	"strings"
)

type PostHandlers struct {
	FrontEndService *front_end_repo.FrontEndRepo
}

func (h *PostHandlers) PostStorePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	data, err, status := h.FrontEndService.GetFormData(r)
	if err != nil {
		log.Println(err)
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error Getting Post Data from request :%v", err), status)
		return
	}
	resp, err := h.FrontEndService.Do(r, w, "POST", "/store-post", bytes.NewBuffer(data))
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	defer resp.Body.Close()
	http.Redirect(w, r, "/posts", http.StatusFound)
}

func (h *PostHandlers) PostCreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	content := r.FormValue("comment-content")
	if content == "" {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Comment is empty"), 400)
		return
	}

	data, _ := json.Marshal(content)

	r.Header.Set("Content-Type", "application/json")

	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, bytes.NewBuffer(data))
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	originurl := "/postbyid" + strings.TrimPrefix(urltocallbackend, "/create-comment")
	http.Redirect(w, r, originurl, http.StatusFound)
}

func (h *PostHandlers) PostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
		})
		return
	}
	defer r.Body.Close()

	log.Println("Frontend received login request, proxying to backend...")

	// Proxy to backend
	resp, err := h.FrontEndService.Do(r, w, "POST", "/login", bytes.NewReader(body))
	if err != nil {
		log.Println("Backend connection error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
		})
		return
	}
	defer resp.Body.Close()

	log.Printf("Backend responded with status: %d\n", resp.StatusCode)

	// Read backend response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading backend response:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
		})
		return
	}

	log.Printf("Backend response body: %s\n", string(responseBody))

	// Forward backend response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func (h *PostHandlers) PostResendVerification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
			"resend":   true,
		})
		return
	}
	defer r.Body.Close()

	resp, err := h.FrontEndService.Do(r, w, "POST", "/resend-verification", bytes.NewReader(body))
	if err != nil {
		log.Println("Backend connection error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
			"resend":   true,
		})
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading backend response:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"success":  false,
			"verified": false,
			"resend":   true,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func (h *PostHandlers) PostSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  []string{"Failed to read request"},
		})
		return
	}
	defer r.Body.Close()

	log.Println("Frontend received signup request, proxying to backend...")

	// Proxy to backend
	resp, err := h.FrontEndService.Do(r, w, "POST", "/signup", bytes.NewReader(body))
	if err != nil {
		log.Println("Backend connection error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  []string{"Backend server unavailable"},
		})
		return
	}
	defer resp.Body.Close()

	log.Printf("Backend responded with status: %d\n", resp.StatusCode)

	// Read backend response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading backend response:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"errors":  []string{"Failed to process response"},
		})
		return
	}

	log.Printf("Backend response body: %s\n", string(responseBody))

	// Forward backend response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

func (h *PostHandlers) PostLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	resp, err := h.FrontEndService.Do(r, w, "POST", "/logout", &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *PostHandlers) PostLikePostOrComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	utils.SuccessResponse(w, "Successful reaction from user", 200)
}

func (h *PostHandlers) PostDislikePostOrComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	utils.SuccessResponse(w, "Successful reaction from user", 200)
}

func (h *PostHandlers) PostSeeNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	utils.SuccessResponse(w, "Successfully seen notification", 200)
}

func (h *PostHandlers) PostRemovePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	utils.SuccessResponse(w, "User Successfully deleted post", 200)
}

func (h *PostHandlers) PostRemoveComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	urltocallbackend := r.URL.Path
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, &bytes.Buffer{})
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}
	utils.SuccessResponse(w, "User Successfully deleted comment", 200)
}

func (h *PostHandlers) PostEditPost(w http.ResponseWriter, r *http.Request) {
	urltocallbackend := r.URL.Path
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}

	data, err, status := h.FrontEndService.GetFormData(r)
	if err != nil {
		log.Println(err)
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Error Getting Data from Request :%v", err), status)
		return
	}
	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, bytes.NewBuffer(data))
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	defer resp.Body.Close()
	utils.SuccessResponse(w, "Post was successfully editted", 200)
}

func (h *PostHandlers) PostEditComment(w http.ResponseWriter, r *http.Request) {
	urltocallbackend := r.URL.Path
	if r.Method != http.MethodPost {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, errors.New("Method Not Allowed"), 405)
		return
	}
	content := r.FormValue("comment-content")

	data, _ := json.Marshal(content)

	resp, err := h.FrontEndService.Do(r, w, "POST", urltocallbackend, bytes.NewBuffer(data))
	if err != nil {
		h.FrontEndService.FrontEndServerErrorwithHTML(w, fmt.Errorf("Backend connection error:%v", err), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		h.FrontEndService.ErrorFromBackEndHtml(resp, w)
		return
	}

	utils.SuccessResponse(w, "Comment was successfully editted", 200)
}
