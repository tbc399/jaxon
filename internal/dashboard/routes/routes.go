package routes

import (
	"net/http"

	dashTemplates "jaxon.app/jaxon/internal/dashboard/templates"
	"jaxon.app/jaxon/internal/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /dashboard", getDashboard)
	router.HandleFunc("GET /dashboard/", getDashboard)
}

func getDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := dashTemplates.Dashboard()
	templates.App("Dashboard", "dashboard", dashboard, "your@email.com").Render(r.Context(), w)
}
