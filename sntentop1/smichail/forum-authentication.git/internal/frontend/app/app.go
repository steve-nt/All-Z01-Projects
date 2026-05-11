package app

import (
	"forum-authentication/internal/backend/models"
	"forum-authentication/internal/frontend/handlers"
	"forum-authentication/internal/frontend/middleware"
	front_end_repo "forum-authentication/internal/frontend/repositories"
	"forum-authentication/internal/frontend/routes"
	"forum-authentication/internal/utils"
	"html/template"
	"net/http"
	"os"
)

func New() (http.Handler, *models.Config) {
	CONFIG := utils.DecodeConf()

	// 🔥 Override από Docker ENV variable (αν υπάρχει)
	if v := os.Getenv("BACKEND_URL"); v != "" {
		CONFIG.Api.Api_base_url = v
	}

	tmpl := template.Must(template.ParseGlob("../../assets/templates/*.page.html"))

	frontendservice := front_end_repo.NewFrontEndRepo(tmpl, &CONFIG)
	posthandlers := &handlers.PostHandlers{FrontEndService: frontendservice}
	gethandlers := &handlers.GetHandlers{FrontEndService: frontendservice}
	rlMW := middleware.InitializeIPRateLimiter(&CONFIG)
	router := &routes.Router{
		PostHandlers:           posthandlers,
		GetHandlers:            gethandlers,
		RateLimitingMiddleware: rlMW,
	}
	return router.Setup(), &CONFIG
}
