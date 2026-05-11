// internal/app/app.go
package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"forum-image-upload/internal/backend/db/sqlite"
	"forum-image-upload/internal/backend/handlers"
	"forum-image-upload/internal/backend/middleware"
	"forum-image-upload/internal/backend/routes"
	"forum-image-upload/internal/backend/services"
	"forum-image-upload/internal/utils"
)

func New() (http.Handler, *sql.DB) {
	CONFIG := utils.DecodeConf()
	logfile, err := os.OpenFile("../../back-up/update-up-migration.sql", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	db := sqlite.NewDatabase(&CONFIG, logfile)

	idleTTL := time.Duration(time.Minute) * time.Duration(CONFIG.Durations.IdleTTL)
	absoluteTTL := time.Duration(time.Hour) * time.Duration(CONFIG.Durations.AbsoluteTTL)
	// Repos
	postRepo := sqlite.NewPostRepo(db, logfile)
	userRepo := sqlite.NewUserRepo(db, logfile, postRepo)
	socialUserRepo := sqlite.NewSocialUserRepo(db, logfile, userRepo)
	sessionRepo := sqlite.NewSessionRepo(db, logfile)

	// Services
	postService := services.NewPostService(postRepo)
	userService := services.NewUserService(userRepo, postRepo)
	socialUserService := services.NewSocialUserService(socialUserRepo, userRepo)
	sessionService := services.NewSessionService(sessionRepo, userRepo, idleTTL, absoluteTTL)

	// Handlers
	postHandler := &handlers.PostHandler{PostService: postService}
	userHandler := &handlers.UserHandler{UserService: userService}
	socialAuthHandler := &handlers.SocialAuthHandler{
		SocialService:  socialUserService,
		SessionService: sessionService,
		ClientIDs: map[string]string{
			"google":   os.Getenv("GOOGLE_CLIENT_ID"),
			"github":   os.Getenv("GITHUB_CLIENT_ID"),
			"facebook": os.Getenv("FACEBOOK_CLIENT_ID"),
		},
		ClientSecrets: map[string]string{
			"google":   os.Getenv("GOOGLE_CLIENT_SECRET"),
			"github":   os.Getenv("GITHUB_CLIENT_SECRET"),
			"facebook": os.Getenv("FACEBOOK_CLIENT_SECRET"),
		},
		RedirectURIs: map[string]string{
			"google":   "http://localhost:8080/auth/callback?provider=google",
			"github":   "http://localhost:8080/auth/callback?provider=github",
			"facebook": "http://localhost:8080/auth/callback?provider=facebook",
		},
		FrontEndBase: "https://localhost:3000",
		Config:       &CONFIG,
	}

	sessionHandler := &handlers.SessionHandler{
		SessionService: sessionService,
		UserService:    userService,
		Config:         &CONFIG,
	}

	// Middleware
	authMW := &middleware.AuthMiddleware{SessionService: sessionService, Config: &CONFIG}

	// Router
	router := &routes.Router{
		PostHandler:       postHandler,
		UserHandler:       userHandler,
		SocialAuthHandler: socialAuthHandler,
		SessionHandler:    sessionHandler,
		AuthMiddleware:    authMW,
	}

	return router.Setup(), db
}
