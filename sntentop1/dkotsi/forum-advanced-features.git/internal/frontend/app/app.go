package app

import (
	"forum-advanced-features/internal/backend/models"
	"forum-advanced-features/internal/frontend/handlers"
	"forum-advanced-features/internal/frontend/middleware"
	front_end_repo "forum-advanced-features/internal/frontend/repositories"
	"forum-advanced-features/internal/frontend/routes"
	"forum-advanced-features/internal/utils"
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
