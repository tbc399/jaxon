package routes

import (
	"net/http"

	//"github.com/jmoiron/sqlx"
	//profilemods "jaxon.app/jaxon/internal/profile/models"
	profiletemps "jaxon.app/jaxon/internal/profile/templates"
	"jaxon.app/jaxon/internal/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /settings/profile", getSettingsFullPage)
	router.HandleFunc("GET /settings/partial", getSettingsPartialPage)
	// router.HandleFunc("GET /settings/accounts-tab", getAccountsTab)
	// router.HandleFunc("GET /settings/assets", getAssetsFullPage)
	// router.HandleFunc("GET /settings/assets-tab", getAssetsTab)
}

func getSettingsFullPage(w http.ResponseWriter, r *http.Request) {
	// db := r.Context().Value("db").(*sqlx.DB)
	// userId := r.Context().Value("userId").(string)

	// profile, err := profilemods.Fetch(userId, db)

	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	profileTab := profiletemps.ProfileTab()
	profilePartial := profiletemps.Profile(profileTab, "profile")
	templates.App("Profile", "settings", profilePartial).Render(r.Context(), w)
}

func getSettingsPartialPage(w http.ResponseWriter, r *http.Request) {
	profileTab := profiletemps.ProfileTab()
	profiletemps.Profile(profileTab, "profile").Render(r.Context(), w)
}
