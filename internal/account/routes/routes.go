package routes

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	accountModels "jaxon.app/jaxon/internal/account/models/accounts"
	assetModels "jaxon.app/jaxon/internal/account/models/assets"
	accountsTemplates "jaxon.app/jaxon/internal/account/templates"
	"jaxon.app/jaxon/internal/templates"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /accounts", getAccountsFullPage)
	router.HandleFunc("GET /accounts/partial", getAccountsPartialPage)
	router.HandleFunc("GET /accounts/accounts-tab", getAccountsTab)
	router.HandleFunc("GET /accounts/assets", getAssetsFullPage)
	router.HandleFunc("GET /accounts/assets-tab", getAssetsTab)
}

func groupAccounts(accounts []accountModels.Account) map[string][]accountModels.Account {

	accountMap := make(map[string][]accountModels.Account)

	translate := map[string]string{
		"cash":       "Cash",
		"credit":     "Credit",
		"investment": "Investment",
		"other":      "Other",
		"manual":     "Other",
	}

	for _, account := range accounts {
		accts, ok := accountMap[translate[account.Type]]
		if !ok {
			accts = make([]accountModels.Account, 0)
		}
		accts = append(accts, account)
		accountMap[translate[account.Type]] = accts
	}

	return accountMap
}

func getAccountsFullPage(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	accounts, err := accountModels.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	acctMap := groupAccounts(accounts)

	stripePubKey := os.Getenv("STRIPE_PUB_KEY")

	accountsTab := accountsTemplates.AccountsTab(acctMap, stripePubKey)
	accountsPartial := accountsTemplates.Accounts(accountsTab, "accounts")
	templates.App(
		"Accounts",
		"accounts",
		accountsPartial,
	).Render(r.Context(), w)

}

func getAssetsFullPage(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	assets, err := assetModels.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	assetsTab := accountsTemplates.AssetsTab(assets)
	assetsPartial := accountsTemplates.Accounts(assetsTab, "assets")
	templates.App(
		"Assets",
		"accounts",
		assetsPartial,
	).Render(r.Context(), w)

}

func getAccountsPartialPage(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	accounts, err := accountModels.FetchAll(userId, db)

	acctMap := groupAccounts(accounts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stripePubKey := os.Getenv("STRIPE_PUB_KEY")
	accountsTab := accountsTemplates.AccountsTab(acctMap, stripePubKey)
	accountsTemplates.Accounts(accountsTab, "accounts").Render(r.Context(), w)

}

func getAccountsTab(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	accounts, err := accountModels.FetchAll(userId, db)

	acctMap := groupAccounts(accounts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stripePubKey := os.Getenv("STRIPE_PUB_KEY")

	accountsTemplates.AccountsTab(acctMap, stripePubKey).Render(r.Context(), w)

}

func getAssetsTab(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	assets, err := assetModels.FetchAll(userId, db)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accountsTemplates.AssetsTab(assets).Render(r.Context(), w)

}
