package routes

import (
	"forum/internal/frontend/handlers"
	"forum/internal/frontend/middleware"
	"net/http"
)

type Router struct {
	PostHandlers           *handlers.PostHandlers
	GetHandlers            *handlers.GetHandlers
	RateLimitingMiddleware *middleware.IPRateLimiter
}

func (router *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	//to serve the static
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../../static/"))))

	//home
	mux.HandleFunc("GET /", router.GetHandlers.GetLanding)

	//user
	mux.HandleFunc("GET /signup", router.GetHandlers.GetSignUpPage)
	mux.HandleFunc("POST /signup", router.PostHandlers.PostSignUp)
	mux.HandleFunc("GET /login", router.GetHandlers.GetLoginPage)
	mux.HandleFunc("POST /login", router.PostHandlers.PostLogin)
	mux.HandleFunc("POST /logout", router.PostHandlers.PostLogout)
	mux.HandleFunc("GET /profile", router.GetHandlers.GetProfile)
	mux.HandleFunc("POST /see-notification/", router.PostHandlers.PostSeeNotification)

	// posts and comments pages
	mux.HandleFunc("GET /posts", router.GetHandlers.GetHome)
	mux.HandleFunc("GET /postbyid/", router.GetHandlers.GetPostbyID)

	//actions
	mux.HandleFunc("POST /create-comment/", router.PostHandlers.PostCreateComment)
	mux.HandleFunc("POST /like-post/", router.PostHandlers.PostLikePostOrComment)
	mux.HandleFunc("POST /dislike-post/", router.PostHandlers.PostDislikePostOrComment)
	mux.HandleFunc("POST /like-comment/", router.PostHandlers.PostLikePostOrComment)
	mux.HandleFunc("POST /dislike-comment/", router.PostHandlers.PostDislikePostOrComment)
	mux.HandleFunc("POST /resend-verification", router.PostHandlers.PostResendVerification)
	mux.HandleFunc("GET /auth/signup", router.GetHandlers.GetSocialSignup)
	mux.HandleFunc("GET /auth/callback", router.GetHandlers.GetSocialSignupCallback)
	mux.HandleFunc("GET /verify", router.GetHandlers.VerifyEmail)
	// mux.HandleFunc("GET /userHome", router.GetHandlers.GetUserHome)
	mux.HandleFunc("GET /create-post", router.GetHandlers.GetCreatePost)
	mux.HandleFunc("POST /create-post", router.PostHandlers.PostStorePost)
	mux.HandleFunc("POST /remove-post/", router.PostHandlers.PostRemovePost)
	mux.HandleFunc("POST /remove-comment/", router.PostHandlers.PostRemoveComment)
	mux.HandleFunc("POST /edit-post/", router.PostHandlers.PostEditPost)
	mux.HandleFunc("POST /edit-comment/", router.PostHandlers.PostEditComment)

	return router.RateLimitingMiddleware.Handler(mux)
}
