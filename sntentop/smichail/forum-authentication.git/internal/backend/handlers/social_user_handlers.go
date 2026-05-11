package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum-authentication/internal/backend/models"
	"forum-authentication/internal/backend/services"
	"forum-authentication/internal/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SocialAuthHandler struct {
	SocialService  *services.SocialUserService
	SessionService *services.SessionService
	ClientIDs      map[string]string
	ClientSecrets  map[string]string
	RedirectURIs   map[string]string
	FrontEndBase   string
	Config         *models.Config
}

// /auth/login?provider=google
func (h *SocialAuthHandler) RedirectToProvider(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	provider := r.URL.Query().Get("provider")
	if provider == "" {
		utils.JsonResponse(w, "Missing provider", http.StatusBadRequest)
		return
	}

	clientID, ok := h.ClientIDs[provider]
	redirectURI, ok2 := h.RedirectURIs[provider]
	if !ok || !ok2 {
		utils.JsonResponse(w, "Unsupported provider", http.StatusBadRequest)
		return
	}

	// Generate state (in production, use a secure random string and store it in session)
	state := "static_state_for_now"

	var authURL string
	switch provider {
	case "google":
		authURL = "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
			"client_id":     {clientID},
			"redirect_uri":  {redirectURI},
			"response_type": {"code"},
			"scope":         {"openid email profile"},
			"state":         {state},
			"access_type":   {"offline"},
		}.Encode()
	case "github":
		authURL = "https://github.com/login/oauth/authorize?" + url.Values{
			"client_id":    {clientID},
			"redirect_uri": {redirectURI},
			"scope":        {"user:email"},
			"state":        {state},
		}.Encode()
	case "facebook":
		authURL = "https://www.facebook.com/v20.0/dialog/oauth?" + url.Values{
			"client_id":    {clientID},
			"redirect_uri": {redirectURI},
			"state":        {state},
			"scope":        {"email"},
		}.Encode()
	default:
		utils.JsonResponse(w, "Unsupported provider", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *SocialAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JsonResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	provider := r.URL.Query().Get("provider")
	code := r.URL.Query().Get("code")
	// state := r.URL.Query().Get("state")

	if provider == "" || code == "" {
		utils.JsonResponse(w, "Missing provider or code", http.StatusBadRequest)
		return
	}

	userInfo, err := h.handleOAuth(ctx, provider, code)
	if err != nil {
		log.Printf("OAuth error for provider %s: %v", provider, err)
		utils.JsonResponse(w, "OAuth process failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Register or fetch user
	user, err := h.SocialService.SocialRegister(ctx, userInfo.Provider, userInfo.ID, userInfo.Email, userInfo.Name, "user")
	if err != nil {
		log.Printf("SocialRegister failed for %s: %v", userInfo.Email, err)
		utils.JsonResponse(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cookieValue, err := h.SessionService.CreateOrGet(r.Context(), user.UUID)
	if err != nil {
		log.Printf("Session creation failed for user %s: %v", user.Username, err)
		utils.JsonResponse(w, "Session creation failed", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     h.Config.CookieName,
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   900,
	})

	// redirect browser back to frontend home (or other URL)
	redirectTo := h.FrontEndBase + "/"
	http.Redirect(w, r, redirectTo, http.StatusFound)

	// Return success - frontend will redirect
	utils.JsonResponse(w, map[string]string{
		"message": "Login successful",
		"user_id": user.UUID,
	}, http.StatusOK)
}

func (h *SocialAuthHandler) handleOAuth(ctx context.Context, provider, code string) (SocialUserInfo, error) {
	switch provider {
	case "google":
		return h.googleFlow(ctx, code)
	case "github":
		return h.githubFlow(ctx, code)
	case "facebook":
		return h.facebookFlow(ctx, code)
	default:
		return SocialUserInfo{}, errors.New("unsupported provider")
	}
}

type SocialUserInfo struct {
	ID       string
	Email    string
	Name     string
	Provider string
}

/* -----------------------------
   Provider-specific flows
   ----------------------------- */

// Google
func (h *SocialAuthHandler) googleFlow(ctx context.Context, code string) (SocialUserInfo, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", h.ClientIDs["google"])
	data.Set("client_secret", h.ClientSecrets["google"])
	data.Set("redirect_uri", h.RedirectURIs["google"])
	data.Set("grant_type", "authorization_code")

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return SocialUserInfo{}, err
	}

	req2, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req2.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp2.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&info); err != nil {
		return SocialUserInfo{}, err
	}

	return SocialUserInfo{ID: info.ID, Email: info.Email, Name: info.Name, Provider: "google"}, nil
}

// GitHub
func (h *SocialAuthHandler) githubFlow(ctx context.Context, code string) (SocialUserInfo, error) {
	data := url.Values{}
	data.Set("client_id", h.ClientIDs["github"])
	data.Set("client_secret", h.ClientSecrets["github"])
	data.Set("code", code)
	data.Set("redirect_uri", h.RedirectURIs["github"])

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return SocialUserInfo{}, err
	}

	req2, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	req2.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp2.Body.Close()

	var info struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&info); err != nil {
		return SocialUserInfo{}, err
	}

	if info.Email == "" {
		req3, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
		req3.Header.Set("Authorization", "Bearer "+token.AccessToken)
		resp3, err := http.DefaultClient.Do(req3)
		if err == nil {
			defer resp3.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			_ = json.NewDecoder(resp3.Body).Decode(&emails)
			for _, e := range emails {
				if e.Primary {
					info.Email = e.Email
					break
				}
			}
		}
	}

	name := info.Name
	if name == "" {
		name = info.Login
	}

	return SocialUserInfo{
		ID:       fmt.Sprint(info.ID),
		Email:    info.Email,
		Name:     name,
		Provider: "github",
	}, nil
}

// Facebook
func (h *SocialAuthHandler) facebookFlow(ctx context.Context, code string) (SocialUserInfo, error) {
	tokenURL := fmt.Sprintf(
		"https://graph.facebook.com/v20.0/oauth/access_token?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		url.QueryEscape(h.ClientIDs["facebook"]),
		url.QueryEscape(h.RedirectURIs["facebook"]),
		url.QueryEscape(h.ClientSecrets["facebook"]),
		url.QueryEscape(code),
	)

	resp, err := http.Get(tokenURL)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return SocialUserInfo{}, err
	}

	userURL := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email&access_token=%s", url.QueryEscape(token.AccessToken))
	resp2, err := http.Get(userURL)
	if err != nil {
		return SocialUserInfo{}, err
	}
	defer resp2.Body.Close()

	var info struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&info); err != nil {
		return SocialUserInfo{}, err
	}

	return SocialUserInfo{ID: info.ID, Email: info.Email, Name: info.Name, Provider: "facebook"}, nil
}
