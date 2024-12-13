package routes

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	budgettemps "jaxon.app/jaxon/internal/budget/templates"
	"jaxon.app/jaxon/internal/templates"
	budgetmods "jaxon.app/jaxon/internal/budget/models/budgets"
	catmods "jaxon.app/jaxon/internal/budget/models/categories"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /budgets", getBudgetsFullPage)
	router.HandleFunc("GET /budgets/partial", getBudgetsPartial)
	router.HandleFunc("GET /budgets/create", getBudgetCreatePage)
	router.HandleFunc("GET /budgets/categories", getCategoriesPage)
	//router.HandleFunc("GET /budgets/categories/create", getCategoriesCreatePage)
	//router.HandleFunc("POST /budgets", createBudget)
	//router.HandleFunc("GET /budgets/{id}/edit", getOneTimePass)
	//router.HandleFunc("GET /login/magic/{magic_token}", getOtpValidationPage)
	//router.HandleFunc("POST /login/magic/{magic_token}", submitOtpValidation)
}

func getBudgetsFullPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	now := time.Now().UTC()
	budgets, err := budgetmods.FetchAllByMonth(userId, now.Year(), now.Month(), db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	budgetPartial := budgettemps.Budgets(budgets, "budgets")
	templates.App("Budgets", "budgets", budgetPartial).Render(r.Context(), w)

}

func getBudgetsPartial(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	now := time.Now().UTC()
	budgets, err := budgetmods.FetchAllByMonth(userId, now.Year(), now.Month(), db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	budgettemps.Budgets(budgets, "budgets").Render(r.Context(), w)
}

func getBudgetCreatePage(w http.ResponseWriter, r *http.Request) {
	hxRequest := r.Header.Get("Hx-Request")
	categories := []catmods.Category{}
	if hxRequest == "true" {
		budgettemps.BudgetCreate(categories).Render(r.Context(), w)
	} else {
		createPartial := budgettemps.BudgetCreate(categories)
		templates.App("Create Budget", "budgets", createPartial).Render(r.Context(), w)
	}
}

func getCategoriesPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	hxRequest := r.Header.Get("Hx-Request")
	cats, err := catmods.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if hxRequest == "true" {
		budgettemps.Categories(cats).Render(r.Context(), w)
	} else {
		createPartial := budgettemps.Categories(cats)
		templates.App("Categories", "budgets", createPartial).Render(r.Context(), w)
	}
}
