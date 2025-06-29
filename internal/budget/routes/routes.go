package routes

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	budgetmods "jaxon.app/jaxon/internal/budget/models/budgets"
	catmods "jaxon.app/jaxon/internal/budget/models/categories"
	"jaxon.app/jaxon/internal/budget/services"
	budgettemps "jaxon.app/jaxon/internal/budget/templates"
	"jaxon.app/jaxon/internal/templates"
	"jaxon.app/jaxon/internal/transaction/models"
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
	router.HandleFunc("GET /landing", getLandingPage)
	// router.HandleFunc("GET /budgets/income/{id}", getIncomeDetailPage)
}

func getLandingPage(w http.ResponseWriter, r *http.Request) {
	templates.LandingPage().Render(r.Context(), w)
	// templates.("Budgets", "budgets", budgetPartial).Render(r.Context(), w)
}

func getBudgets(w http.ResponseWriter, r *http.Request) (*services.BudgetOverview, []budgetmods.BudgetView, error) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	now := time.Now().UTC()
	budgets, err := budgetmods.FetchBudgetViewsByMonth(userId, now.Year(), now.Month(), db)
	if err != nil {
		return nil, nil, errors.New("Failed to get budgets")
	}

	period, err := budgetmods.FetchCurrentPeriod(userId, db)
	if err != nil {
		return nil, nil, errors.New("Failed to get current period")
	}

	overview, err := services.GetBudgetOverview(userId, period, db)
	if err != nil {
		return nil, nil, errors.New("Failed to get overview")
	}

	return overview, budgets, nil
}

func getBudgetsFullPage(w http.ResponseWriter, r *http.Request) {
	overview, budgets, err := getBudgets(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	budgetPartial := budgettemps.Budgets(overview, budgets, "budgets")
	templates.App("Budgets", "budgets", budgetPartial).Render(r.Context(), w)
}

func getBudgetsPartial(w http.ResponseWriter, r *http.Request) {
	overview, budgets, err := getBudgets(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	budgettemps.Budgets(overview, budgets, "budgets").Render(r.Context(), w)
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
	// db := r.Context().Value("db").(*sqlx.DB)
	// userId := r.Context().Value("userId").(string)
}

func getBudgetDetailPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	budgetId := r.PathValue("id")

	budget, err := budgetmods.FetchBudgetView(budgetId, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hxRequest := r.Header.Get("Hx-Request")

	// transactions, err := models.FetchForCategoryInRange()
	transactions := []models.TransactionView{}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if hxRequest == "true" {
		budgettemps.BudgetDetail(budget, transactions).Render(r.Context(), w)
	} else {
		updatePartial := budgettemps.BudgetDetail(budget, transactions)
		templates.App("Budget Detail", "budgets", updatePartial).Render(r.Context(), w)
	}
}

/*
func getIncomeDetailPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	budgetId := r.PathValue("id")

	budget, err := budgetmods.FetchIncomeBudgetView(budgetId, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hxRequest := r.Header.Get("Hx-Request")

	//period, err := budgetmods.FetchPeriod(budget.PeriodId, userId, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactions, err := models.FetchForCategoryInRange(userId, budget.CategoryId, period.Start, period.End, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if hxRequest == "true" {
		budgettemps.BudgetDetail(budget, transactions).Render(r.Context(), w)
	} else {
		updatePartial := budgettemps.BudgetDetail(budget, transactions)
		templates.App("Budget Detail", "budgets", updatePartial).Render(r.Context(), w)
	}
}
*/

func createBudget(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	categoryId := r.PostFormValue("category")
	amount, err := strconv.ParseInt(r.PostFormValue("amount"), 10, 0)
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
	// categoryId := r.PostFormValue("category")

	amount, err := strconv.ParseInt(r.PostFormValue("amount"), 10, 0)
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
