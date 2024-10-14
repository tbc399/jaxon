package routes

import (
	//"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/templates"
	"jaxon.app/jaxon/internal/transaction/models"
	transactionTemplates "jaxon.app/jaxon/internal/transaction/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /transactions", getTransactions)
	router.HandleFunc("GET /transactions/partial", getTransactionsPartial)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactions, err := models.FetchMany(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactionsPartial := transactionTemplates.Transactions(transactions, "transactions")
	templates.App("Transactions", "transactions", transactionsPartial, "your@email.com").Render(r.Context(), w)

}

func getTransactionsPartial(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactions, err := models.FetchMany(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactionTemplates.Transactions(transactions, "transactions").Render(r.Context(), w)

}
