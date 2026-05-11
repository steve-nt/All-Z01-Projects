package handlers

import (
	"encoding/json"
	"fmt"
	"forum-image-upload/internal/backend/models"
	"forum-image-upload/internal/backend/services"
	"forum-image-upload/internal/utils"
	"log"
	"net/http"
)

type SessionHandler struct {
	SessionService *services.SessionService
	UserService    *services.UserService
	Config         *models.Config
}

func (h *SessionHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var userInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decoder.Decode(&userInput); err != nil {
		log.Println(err)
	}

	user, err := h.UserService.ValidateCredentials(r.Context(), userInput.Username, userInput.Password)
	if err != nil {
		log.Printf("Login failed for user %s: %v", userInput.Username, err)
		utils.ErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.Verified {
		log.Printf("Login attempt for unverified user: %s", user.Username)
		resp := map[string]any{
			"success":  false,
			"verified": false,
			"message":  "Email not verified. Please verify your account.",
			"resend":   true,
		}
		log.Printf("Response being sent: %+v", resp)
		utils.JsonResponse(w, resp, http.StatusOK)
		return
	}

	log.Printf("Login successful for user: %s", user.Username)

	cookieValue, err := h.SessionService.CreateOrGet(r.Context(), user.UUID)
	if err != nil {
		log.Printf("Session creation failed for user %s: %v", user.Username, err)
		utils.ErrorResponse(w, "Session creation failed", http.StatusInternalServerError)
		return
	}
	log.Printf("Session created successfully for user %s, cookie: %s", user.Username, cookieValue)

	http.SetCookie(w, &http.Cookie{
		Name:     h.Config.CookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   900,
	})
	utils.JsonResponse(w, map[string]any{
		"success":  true,
		"verified": true,
		"message":  "User Logged In",
	}, http.StatusOK)
}

func (h *SessionHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("example")
	if cookie != nil {
		if err := h.SessionService.DeleteByCookie(r.Context(), cookie.Value); err != nil {
			log.Println(err)
		}
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
}

func (h *SessionHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Βρες τον χρήστη από τη βάση
	user, err := h.UserService.GetByUsername(r.Context(), req.Username)
	if err != nil {
		utils.ErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	// Αν είναι ήδη verified
	if user.Verified {
		utils.JsonResponse(w, map[string]any{
			"success": true,
			"message": "Your account is already verified.",
		}, http.StatusOK)
		return
	}

	// Δημιούργησε νέο verification link
	link := fmt.Sprintf("http://localhost:8080/verify?uuid=%s", user.UUID)

	// Στείλε email
	if err := utils.SendVerificationEmail(user.Mail, link); err != nil {
		log.Println("Failed to resend verification email:", err)
		utils.ErrorResponse(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, map[string]any{
		"success": true,
		"message": "Verification email resent successfully!",
	}, http.StatusOK)
}

// GetAuthStatus returns the authentication status and user info
// func (h *SessionHandler) GetAuthStatus(w http.ResponseWriter, r *http.Request, user models.User) {
// 	var response struct {
// 		IsAuthenticated bool   `json:"isAuthenticated"`
// 		Username        string `json:"username,omitempty"`
// 		Role            string `json:"role,omitempty"`
// 	}
//
// 	// Check if user is authenticated (not empty user struct)
// 	if (user == models.User{}) {
// 		response.IsAuthenticated = false
// 	} else {
// 		response.IsAuthenticated = true
// 		response.Username = user.Username
// 		response.Role = user.Role
// 	}
//
// 	utils.JsonResponse(w, response, 200)
// }
