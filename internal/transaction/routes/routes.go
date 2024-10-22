package routes

import (
	//"log/slog"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	accountmods "jaxon.app/jaxon/internal/account/models/accounts"
	"jaxon.app/jaxon/internal/templates"
	"jaxon.app/jaxon/internal/transaction/models"
	transactiontemps "jaxon.app/jaxon/internal/transaction/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /transactions", getTransactionsFullPage)
	router.HandleFunc("GET /transactions/partial", getTransactionsPartial)
	router.HandleFunc("GET /transactions/upload", getTransactionsUpload)
}

func getTransactionsFullPage(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactions, err := models.FetchMany(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accounts, err := accountmods.FetchAll(userId, db)

	if err != nil {
		slog.Error("Failed to fetch accounts")
	}

	transactionsPartial := transactiontemps.Transactions(transactions, accounts, "transactions")
	templates.App("Transactions", "transactions", transactionsPartial).Render(r.Context(), w)

}

func getTransactionsPartial(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactions, err := models.FetchMany(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accounts, err := accountmods.FetchAll(userId, db)

	if err != nil {
		slog.Error("Failed to fetch accounts")
	}

	transactiontemps.Transactions(transactions, accounts, "transactions").Render(r.Context(), w)

}

func getTransactionsUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)	
	transactiontemps.UploadPage().Render(r.Context(), w)
}
