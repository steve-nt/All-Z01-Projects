package middleware

import (
	"context"
	"forum-authentication/internal/backend/models"
	"forum-authentication/internal/backend/services"
	"log"
	"net/http"
)

type AuthMiddleware struct {
	SessionService *services.SessionService
	Config         *models.Config
}

// Middleware wrapper
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/resend-verification" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, _ := r.Cookie(m.Config.CookieName)

		var role models.User
		refreshCookie := ""

		if cookie == nil {
			// guest user is represented as an empty user struct in the handlers
			role = models.User{}
		} else {
			user, newCookie, err := m.SessionService.ValidateAndMaybeRefresh(r.Context(), cookie.Value)
			refreshCookie = newCookie
			if err != nil {
				log.Println(err)
				role = models.User{}
			} else {
				role = user
			}
		}

		ctx := context.WithValue(r.Context(), "user", role)

		if refreshCookie != "" {
			cookie.Value = refreshCookie
			cookie.MaxAge = 3600
			http.SetCookie(w, cookie)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
