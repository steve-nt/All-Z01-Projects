package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"forum/db"
	"forum/sessions"
)

// OAuthGitHubStart starts the GitHub OAuth flow.
// It redirects the browser to GitHub's authorization endpoint.
// We include a CSRF state token (stored in a short-lived cookie).
func OAuthGitHubStart(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID"))
	if clientID == "" {
		RenderError(w, r, http.StatusInternalServerError, "GitHub OAuth not configured.")
		return
	}

	// State token protects against CSRF during OAuth redirects.
	state, err := newStateToken()
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not start OAuth.")
		return
	}
	setStateCookie(w, r, state)

	// Callback must match the GitHub OAuth App settings.
	redirectURI := oauthRedirectBase(r) + "/auth/github/callback"

	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", "read:user user:email")
	q.Set("state", state)

	// Force GitHub to show a login prompt instead of silently reusing an existing session.
	// This ensures that after a forum logout, clicking "Login with GitHub" does not auto-log
	// the user back in with the same GitHub account.
	q.Set("prompt", "login")

	authURL := "https://github.com/login/oauth/authorize?" + q.Encode()
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// OAuthGitHubCallback handles GitHub's redirect back to our server.
// It exchanges the authorization code for an access token, fetches user identity,
// then creates/links a local forum user and starts a forum session.
func OAuthGitHubCallback(w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("GITHUB_CLIENT_SECRET"))
	if clientID == "" || clientSecret == "" {
		RenderError(w, r, http.StatusInternalServerError, "GitHub OAuth not configured.")
		return
	}

	// Validate state to prevent CSRF.
	if err := verifyState(r); err != nil {
		RenderError(w, r, http.StatusBadRequest, "Invalid OAuth state.")
		return
	}

	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if code == "" {
		RenderError(w, r, http.StatusBadRequest, "Missing OAuth code.")
		return
	}

	redirectURI := oauthRedirectBase(r) + "/auth/github/callback"

	// Exchange code for token.
	accessToken, err := githubExchangeCodeForToken(code, clientID, clientSecret, redirectURI)
	if err != nil {
		RenderError(w, r, http.StatusUnauthorized, "GitHub authentication failed.")
		return
	}

	// Fetch GitHub user identity (id + login).
	user, err := githubFetchUser(accessToken)
	if err != nil {
		RenderError(w, r, http.StatusUnauthorized, "GitHub authentication failed.")
		return
	}

	// GitHub may not return email in /user; use /user/emails with user:email scope.
	email, err := githubFetchPrimaryEmail(accessToken)
	if err != nil || email == "" {
		RenderError(w, r, http.StatusUnauthorized, "GitHub email permission is required.")
		return
	}

	// Link provider identity to a local forum user (create if needed).
	userID, err := db.FindOrCreateOAuthUser("github", user.IDStr(), email, user.Login)
	if err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create user.")
		return
	}

	// Remove state cookie once callback completed successfully.
	clearStateCookie(w)

	// Create the forum session cookie + DB session row.
	if err := sessions.CreateSession(w, r, userID); err != nil {
		RenderError(w, r, http.StatusInternalServerError, "Could not create session.")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type githubTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// githubExchangeCodeForToken exchanges the OAuth "code" for an access token.
func githubExchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("token exchange failed")
	}

	var tr githubTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}
	if tr.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	return tr.AccessToken, nil
}

type githubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

func (u githubUser) IDStr() string {
	return strconv.FormatInt(u.ID, 10)
}

// githubFetchUser fetches GitHub user identity.
// We only need stable identity fields (id + login).
func githubFetchUser(accessToken string) (githubUser, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "forum-app")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return githubUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return githubUser{}, errors.New("user fetch failed")
	}

	var u githubUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return githubUser{}, err
	}

	u.Login = strings.TrimSpace(u.Login)
	if u.ID == 0 || u.Login == "" {
		return githubUser{}, errors.New("missing github fields")
	}
	return u, nil
}

type githubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// githubFetchPrimaryEmail returns a verified email address for the GitHub account.
// Prefer primary+verified, else any verified email.
func githubFetchPrimaryEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "forum-app")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("emails fetch failed")
	}

	var emails []githubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified && strings.TrimSpace(e.Email) != "" {
			return strings.ToLower(strings.TrimSpace(e.Email)), nil
		}
	}

	for _, e := range emails {
		if e.Verified && strings.TrimSpace(e.Email) != "" {
			return strings.ToLower(strings.TrimSpace(e.Email)), nil
		}
	}

	return "", errors.New("no usable email")
}