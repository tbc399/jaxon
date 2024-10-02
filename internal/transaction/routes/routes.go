package routes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/templates"
	"jaxon.app/jaxon/internal/transaction/models"
	transactionTemplates "jaxon.app/jaxon/internal/transactions/templates"
)

func Router() * http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", getTransactions)
	return router
}

func getTransactions(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	transactions := models.FetchMany(userId, db)

	transactionsPartial := transactionTemplates.Transactions(&transactions, "transactions")
	templates.App("Transactions", "transactions", transactionsPartial).Render(r.Context(), w)

}
