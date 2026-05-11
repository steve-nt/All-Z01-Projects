package routes

import (
	"net/http"

	"forum-advanced-features/internal/backend/handlers"
	"forum-advanced-features/internal/backend/middleware"
	"forum-advanced-features/internal/utils"
)

type Router struct {
	PostHandler       *handlers.PostHandler
	UserHandler       *handlers.UserHandler
	SocialAuthHandler *handlers.SocialAuthHandler
	SessionHandler    *handlers.SessionHandler
	AuthMiddleware    *middleware.AuthMiddleware
}

func (r *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	// Posts
	mux.HandleFunc("GET /", r.PostHandler.GetHome)
	mux.HandleFunc("GET /postbyid/", r.PostHandler.GetPostByID)
	mux.HandleFunc("POST /store-post", r.PostHandler.StorePost)
	mux.HandleFunc("POST /edit-post/", r.PostHandler.EditPost)
	mux.HandleFunc("POST /create-comment/", r.PostHandler.CreateComment)
	mux.HandleFunc("POST /edit-comment/", r.PostHandler.EditComment)
	mux.HandleFunc("POST /like-post/", r.PostHandler.LikePost)
	mux.HandleFunc("POST /like-comment/", r.PostHandler.LikeComment)
	mux.HandleFunc("POST /dislike-post/", r.PostHandler.DislikePost)
	mux.HandleFunc("POST /dislike-comment/", r.PostHandler.DislikeComment)
	mux.HandleFunc("POST /remove-post/", r.PostHandler.RemovePost)
	mux.HandleFunc("POST /remove-comment/", r.PostHandler.RemoveComment)

	// Users
	mux.HandleFunc("POST /signup", r.UserHandler.CreateUser)
	mux.HandleFunc("GET /profile", r.UserHandler.GetProfile)
	mux.HandleFunc("POST /see-notification/", r.UserHandler.SeeNotification)
	mux.HandleFunc("GET /verify", r.UserHandler.VerifyEmail)

	// Social signup
	mux.HandleFunc("GET /auth/signup", r.SocialAuthHandler.RedirectToProvider)
	mux.HandleFunc("GET /auth/callback", r.SocialAuthHandler.Callback)
	mux.HandleFunc("GET /debug/auth", func(w http.ResponseWriter, req *http.Request) {
		utils.JsonResponse(w, map[string]interface{}{
			"google_client_id": r.SocialAuthHandler.ClientIDs["google"],
			"github_client_id": r.SocialAuthHandler.ClientIDs["github"],
			"redirect_uris":    r.SocialAuthHandler.RedirectURIs,
		}, http.StatusOK)
	})

	// Sessions
	mux.HandleFunc("POST /login", r.SessionHandler.LoginUser)
	mux.HandleFunc("POST /logout", r.SessionHandler.LogoutUser)
	mux.HandleFunc("POST /resend-verification", r.SessionHandler.ResendVerificationEmail)

	// Wrap with middleware
	return r.AuthMiddleware.Handler(mux)
}
