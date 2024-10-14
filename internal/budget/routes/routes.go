package routes

import (
	"net/http"

	"jaxon.app/jaxon/internal/budget"
	budgetTemplates "jaxon.app/jaxon/internal/budget/templates"
	"jaxon.app/jaxon/internal/templates"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /budgets", getBudgets)
	//router.HandleFunc("POST /budgets", createBudget)
	//router.HandleFunc("GET /budgets/{id}/edit", getOneTimePass)
	//router.HandleFunc("GET /budgets/partial", getBudgetsPartial)
	//router.HandleFunc("GET /login/magic/{magic_token}", getOtpValidationPage)
	//router.HandleFunc("POST /login/magic/{magic_token}", submitOtpValidation)
	return router
}

func getBudgets(w http.ResponseWriter, r *http.Request) {
	//now := time.Now().UTC()
	//budgets := FetchAllByMonth(user_id, now.year, now.month, db)
	budgets := []budget.Budget{}

	budgetPartial := budgetTemplates.Budgets(&budgets, "budgets")
	templates.App("Budgets", "budgets", budgetPartial, "your@email.com").Render(r.Context(), w)
}
