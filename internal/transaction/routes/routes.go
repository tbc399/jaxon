package routes

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	accountmods "jaxon.app/jaxon/internal/account/models/accounts"
	catmods "jaxon.app/jaxon/internal/budget/models/categories"
	"jaxon.app/jaxon/internal/templates"
	"jaxon.app/jaxon/internal/transaction/models"
	"jaxon.app/jaxon/internal/transaction/services"
	transactiontemps "jaxon.app/jaxon/internal/transaction/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /transactions", getTransactionsFullPage)
	router.HandleFunc("GET /transactions/partial", getTransactionsPartial)
	router.HandleFunc("GET /transactions/upload", getTransactionsUpload)
	router.HandleFunc("POST /transactions/upload", uploadTransactions)
	router.HandleFunc("GET /transactions/{transaction_id}/edit", getTransactionEditPartial)
	router.HandleFunc("GET /transactions/{transaction_id}", getTransactionEditPage)
	router.HandleFunc("PUT /transactions/{transaction_id}", updateTransaction)
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

func uploadTransactions(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	file, _, err := r.FormFile("file")

	if err != nil {
		slog.Error("Failed to get the upload file from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// TODO make into a go routine
	services.UploadTransactions(file, userId, db)

	http.Redirect(w, r, "/transactions", http.StatusSeeOther)

}

func getTransactionEditPartial(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactionId := r.PathValue("transaction_id")

	slog.Debug("Getting transaction edit page", "transaction", transactionId)

	transaction, err := models.Fetch(transactionId, db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := catmods.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactiontemps.TransactionEdit(transaction, categories).Render(r.Context(), w)
}

func getTransactionEditPage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	transactionId := r.PathValue("transaction_id")

	slog.Debug("Getting transaction edit page", "transaction", transactionId)

	transaction, err := models.Fetch(transactionId, db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories, err := catmods.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transactiontemps.TransactionEdit(transaction, categories).Render(r.Context(), w)

	transactionEditPartial := transactiontemps.TransactionEdit(transaction, categories)
	templates.App("Transactions", "transactions", transactionEditPartial).Render(r.Context(), w)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)
	
	transactionId := r.PathValue("transaction_id")
	description := r.PostFormValue("description")
	categoryId := r.PostFormValue("category")

	transaction, err := models.Fetch(transactionId, db)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = catmods.Fetch(categoryId, db)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	transaction.Description = description
	transaction.CategoryId = sql.NullString{String: categoryId,Valid: true}
	err = transaction.Save(db)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("HX-Redirect", "/transactions")
	w.WriteHeader(http.StatusOK)
	
}
