package routes

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	budgetmods "jaxon.app/jaxon/internal/budget/models/budgets"
	catmods "jaxon.app/jaxon/internal/budget/models/categories"
	budgettemps "jaxon.app/jaxon/internal/budget/templates"
	"jaxon.app/jaxon/internal/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /budgets", getBudgetsFullPage)
	router.HandleFunc("GET /budgets/partial", getBudgetsPartial)
	router.HandleFunc("GET /budgets/create", getBudgetCreatePage)
	router.HandleFunc("POST /budgets", createBudget)
	router.HandleFunc("GET /budgets/categories", getCategoriesPage)
	router.HandleFunc("POST /budgets/categories", createCategory)
	router.HandleFunc("GET /budgets/{id}", getBudgetDetailPage)
	router.HandleFunc("DELETE /budgets/{id}", removeBudget)
	router.HandleFunc("PUT /budgets/{id}", updateBudget)
}

func getBudgetsFullPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	now := time.Now().UTC()
	budgets, err := budgetmods.FetchBudgetViewsByMonth(userId, now.Year(), now.Month(), db)

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
	budgets, err := budgetmods.FetchBudgetViewsByMonth(userId, now.Year(), now.Month(), db)

	//overview, err := budgetmods.FetchBudgetOverview(userId, now.Year(), now.Month(), db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	budgettemps.Budgets(budgets, "budgets").Render(r.Context(), w)
}

func getBudgetCreatePage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	hxRequest := r.Header.Get("Hx-Request")

	categories, err := catmods.FetchAll(userId, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if hxRequest == "true" {
		budgettemps.BudgetCreate(categories).Render(r.Context(), w)
	} else {
		createPartial := budgettemps.BudgetCreate(categories)
		templates.App("Create Budget", "budgets", createPartial).Render(r.Context(), w)
	}
}

func removeBudget(w http.ResponseWriter, r *http.Request) {
	//db := r.Context().Value("db").(*sqlx.DB)
	//userId := r.Context().Value("userId").(string)
	

}

func getBudgetDetailPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	budgetId := r.PathValue("id")

	budget, err := budgetmods.FetchBudget(budgetId, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hxRequest := r.Header.Get("Hx-Request")

	categories, err := catmods.FetchAll(userId, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if hxRequest == "true" {
		budgettemps.BudgetDetail(budget, categories).Render(r.Context(), w)
	} else {
		updatePartial := budgettemps.BudgetDetail(budget, categories)
		templates.App("Budget Detail", "budgets", updatePartial).Render(r.Context(), w)
	}
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	categoryId := r.PostFormValue("category")
	amount, err := strconv.ParseUint(r.PostFormValue("amount"), 10, 0)
	if err != nil {
		slog.Error("Failed to parse amount", "amount", r.PostFormValue("amount"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category, err := catmods.Fetch(categoryId, db)
	if category == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	period, err := budgetmods.FetchCurrentPeriod(userId, db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if period == nil {
		// create a new period
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Second)
		newPeriod := budgetmods.NewBudgetPeriod(userId, start, end)
		newPeriod.Save(db)
	}

	budget := budgetmods.NewBudget(period.Id, userId, category.Id, amount)
	err = budget.Save(db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/budgets", http.StatusSeeOther)
}

func updateBudget(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	budgetId := r.PathValue("id")
	//categoryId := r.PostFormValue("category")

	amount, err := strconv.ParseUint(r.PostFormValue("amount"), 10, 0)
	if err != nil {
		slog.Error("Failed to parse amount", "amount", r.PostFormValue("amount"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	/*
	category, err := catmods.Fetch(categoryId, db)
	if category == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	*/

	budget, err := budgetmods.FetchBudget(budgetId, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	budget.Amount = amount
	err = budget.Save(db)

	http.Redirect(w, r, "/budgets", http.StatusSeeOther)
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

func createCategory(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	catName := r.PostFormValue("category")

	existingCat, err := catmods.FetchByName(catName, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if existingCat != nil {
		slog.Info("Category already exists", "name", catName, "user", userId)
		w.WriteHeader(http.StatusConflict)
		return
	}

	cat := catmods.NewCategory(catName, catmods.ExpenseCategoryType, userId)
	err = cat.Save(db)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/budgets/categories", http.StatusSeeOther)
}
