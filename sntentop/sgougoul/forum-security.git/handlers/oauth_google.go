package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"forum/db"
	"forum/sessions"
)

// oauthStateCookie stores a short-lived CSRF protection token for OAuth.
// We set it before redirecting to Google and verify it on callback.
const oauthStateCookie = "oauth_state"

// OAuthGoogleStart begins the Google OAuth login flow.
// It generates a random state token (CSRF protection), stores it in an HttpOnly cookie,
// and redirects the user to Google's authorization endpoint.
func OAuthGoogleStart(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	_ = clientSecret // checked in callback too
	if clientID == "" {
		RenderError(w, r, http.StatusInternalServerError, "Google OAuth not configured.")
		return
	}

	// Create CSRF state token and store it in a cookie for validation in callback.
	state, err := newStateToken()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not start OAuth.")
		return
	}
	setStateCookie(w, r, state)

	// Callback URL must match what is configured in Google Cloud Console.
	redirectURI := oauthRedirectBase(r) + "/auth/google/callback"

	// Build Google authorization URL.
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)

	// access_type=online is enough because we do not need refresh tokens.
	q.Set("access_type", "online")

	// prompt=select_account forces Google to show the account chooser (useful for switching accounts).
	q.Set("prompt", "select_account")

	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// OAuthGoogleCallback finishes the Google OAuth flow.
// It verifies the OAuth state (CSRF protection), exchanges the authorization code for an access token,
// fetches the user's profile, links/creates a local forum user, and finally creates a normal forum session.
func OAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	if clientID == "" || clientSecret == "" {
		RenderError(w, r, http.StatusInternalServerError, "Google OAuth not configured.")
		return
	}

	// Validate state from query matches the stored cookie to prevent CSRF.
	if err := verifyState(r); err != nil {
		RenderError(w, r, http.StatusBadRequest, "Invalid OAuth state.")
		return
	}

	// "code" is the short-lived authorization code returned by Google.
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		RenderError(w, r, http.StatusBadRequest, "Missing OAuth code.")
		return
	}

	redirectURI := oauthRedirectBase(r) + "/auth/google/callback"

	// Exchange authorization code for an access token.
	accessToken, err := googleExchangeCodeForToken(code, clientID, clientSecret, redirectURI)
	if err != nil {
		RenderError(w, r, http.StatusUnauthorized, "Google authentication failed.")
		return
	}

	// Fetch OpenID Connect user info (sub/email/name).
	profile, err := googleFetchUserInfo(accessToken)
	if err != nil {
		RenderError(w, r, http.StatusUnauthorized, "Google authentication failed.")
		return
	}

	// Create or link local user, then create our normal forum session.
	preferredUsername := profile.Name
	if preferredUsername == "" && profile.Email != "" {
		preferredUsername = strings.Split(profile.Email, "@")[0]
	}

	// provider="google" and provider_id=profile.Sub uniquely identify the OAuth identity.
	userID, err := db.FindOrCreateOAuthUser("google", profile.Sub, profile.Email, preferredUsername)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create user.")
		return
	}

	// State cookie is only for this flow; remove it after success.
	clearStateCookie(w)

	// Create the forum's own session cookie (separate from Google).
	if err := sessions.CreateSession(w, r, userID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create session.")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// googleTokenResp is the minimal token response we need from Google.
type googleTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// googleExchangeCodeForToken exchanges an authorization code for an access token.
func googleExchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (string, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("redirect_uri", redirectURI)
	form.Set("grant_type", "authorization_code")

	req, _ := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("token exchange failed")
	}

	var tr googleTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	return tr.AccessToken, nil
}

// googleUserInfo is the minimal OpenID user profile we use for creating/linking local users.
type googleUserInfo struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// googleFetchUserInfo calls Google's OpenID userinfo endpoint using the access token.
func googleFetchUserInfo(accessToken string) (googleUserInfo, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://openidconnect.googleapis.com/v1/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return googleUserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return googleUserInfo{}, errors.New("userinfo failed")
	}

	var ui googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return googleUserInfo{}, err
	}

	// Normalize/validate required fields.
	ui.Sub = strings.TrimSpace(ui.Sub)
	ui.Email = strings.TrimSpace(strings.ToLower(ui.Email))
	ui.Name = strings.TrimSpace(ui.Name)

	if ui.Sub == "" || ui.Email == "" {
		return googleUserInfo{}, errors.New("missing profile fields")
	}
	return ui, nil
}

// ---------- shared helpers ----------

// oauthRedirectBase determines the base URL used in redirect_uri.
// In production you should set OAUTH_REDIRECT_BASE explicitly to avoid host/proxy issues.
func oauthRedirectBase(r *http.Request) string {
	if v := strings.TrimSpace(os.Getenv("OAUTH_REDIRECT_BASE")); v != "" {
		return strings.TrimRight(v, "/")
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

// newStateToken creates a random string for OAuth state (CSRF protection).
func newStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// setStateCookie stores the OAuth state token in a short-lived HttpOnly cookie.
func setStateCookie(w http.ResponseWriter, r *http.Request, state string) {
	secure := r != nil && r.TLS != nil

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/",
		MaxAge:   10 * 60, // 10 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

// verifyState compares the query "state" to the cookie value to prevent CSRF.
func verifyState(r *http.Request) error {
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if state == "" {
		return errors.New("missing state")
	}

	c, err := r.Cookie(oauthStateCookie)
	if err != nil || c.Value == "" {
		return errors.New("missing state cookie")
	}

	if subtleConstantTimeEqual(c.Value, state) == false {
		return errors.New("state mismatch")
	}
	return nil
}

// clearStateCookie removes the OAuth state cookie after the flow completes.
func clearStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func subtleConstantTimeEqual(a, b string) bool {
	// Tiny constant-time compare without importing crypto/subtle (still standard though).
	if len(a) != len(b) {
		return false
	}
	var v byte
	ab := []byte(a)
	bb := []byte(b)
	for i := 0; i < len(ab); i++ {
		v |= ab[i] ^ bb[i]
	}
	return v == 0
}

// Optional helper if you ever need to read response bodies for debugging without breaking flow.
func readAll(r *http.Response) []byte {
	if r == nil || r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r.Body)
	return buf.Bytes()
}