package routes

import (
	"net/http"

	dashTemplates "jaxon.app/jaxon/internal/dashboard/templates"
	"jaxon.app/jaxon/internal/templates"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", getDashboard)
	return router
}

func getDashboard(w http.ResponseWriter, r *http.Request) {
	dashboard := dashTemplates.Dashboard()
	templates.App("Dashboard", "dashboard", dashboard).Render(r.Context(), w)
}
