package handlers

import (
	"encoding/json"
	"fmt"
	"forum-advanced-features/internal/backend/models"
	"forum-advanced-features/internal/backend/services"
	"forum-advanced-features/internal/utils"
	"log"
	"net/http"
	"net/url"
)

type UserHandler struct {
	UserService *services.UserService
}

type VerifyResult struct {
	Title   string
	Message string
	Success bool
	Color   string
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponseSignup(w, []string{"Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newUser struct {
		Mail           string `json:"mail"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		RepeatPassword string `json:"repeat_password"`
		Role           string `json:"role"`
	}
	if err := decoder.Decode(&newUser); err != nil {
		log.Println("decode signup payload:", err)
		utils.ErrorResponseSignup(w, []string{"Invalid request payload"}, http.StatusBadRequest)
		return
	}

	// Check if passwords match
	if newUser.Password != newUser.RepeatPassword {
		utils.ErrorResponseSignup(w, []string{"Passwords do not match"}, http.StatusBadRequest)
		return
	}

	// Register user
	if err := h.UserService.Register(r.Context(), newUser.Mail, newUser.Username, newUser.Password, newUser.RepeatPassword, newUser.Role); err != nil {
		// Check if it's a validation error
		if validationErr, ok := err.(services.ValidationError); ok {
			utils.ErrorResponseSignup(w, validationErr.Errors, http.StatusBadRequest)
			return
		}
		// Other errors (database, bcrypt, etc.) are internal server errors
		log.Println("register user:", err)
		utils.ErrorResponseSignup(w, []string{"An error occurred during registration"}, http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "User created successfully \n An email was sent to you to confirm your acoount.", http.StatusCreated)
}

func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Redirect(w, r, "https://localhost:3000/verify?success=false&msg=invalid_link", http.StatusSeeOther)
		return
	}

	err := h.UserService.VerifyEmail(r.Context(), uuid)
	if err != nil {
		http.Redirect(w, r,
			fmt.Sprintf("https://localhost:3000/verify?success=false&msg=%s", url.QueryEscape(err.Error())),
			http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "https://localhost:3000/verify?success=true", http.StatusSeeOther)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := h.UserService.GetProfile(r.Context(), r)
	if err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp models.ProfileResponse

	resp.Username = user.Username
	resp.Notifications = user.Notifications
	resp.Activities = user.Activities
	resp.LikedPosts = user.LikedPosts
	resp.CreatedPosts = user.CreatedPosts
	h.UserService.CountUnseenNotifications(&resp)

	utils.JsonResponse(w, resp, 200)
}

func (h *UserHandler) SeeNotification(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.UserService.SeeNotification(r.Context(), r)
	if err != nil {
		log.Println(err)
		utils.JsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, "Notification Seen Successfully", 200)
}
