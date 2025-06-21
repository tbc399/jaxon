package routes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/plaid/plaid-go/v37/plaid"
	accountModels "jaxon.app/jaxon/internal/account/models/accounts"
	assetModels "jaxon.app/jaxon/internal/account/models/assets"
	accountsTemplates "jaxon.app/jaxon/internal/account/templates"
	plaidModels "jaxon.app/jaxon/internal/plaid/models"
	"jaxon.app/jaxon/internal/templates"
	// plaidModels "jaxon.app/jaxon/internal/plaid/models"
)

func AddRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /accounts", getAccountsFullPage)
	router.HandleFunc("GET /accounts/partial", getAccountsPartialPage)
	router.HandleFunc("GET /accounts/accounts-tab", getAccountsTab)
	router.HandleFunc("GET /accounts/assets", getAssetsFullPage)
	router.HandleFunc("GET /accounts/assets-tab", getAssetsTab)
	router.HandleFunc("POST /accounts/create-link", createLinkToken)
	router.HandleFunc("POST /accounts/exchange-token", exchangePublicToken)
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

type PlaidLinkResponse struct {
	LinkToken string `json:"link_token"`
}

func createLinkToken(w http.ResponseWriter, r *http.Request) {
	plaidClient := r.Context().Value("plaidClient").(*plaid.APIClient)
	userId := r.Context().Value("userId").(string)

	// Create link token request
	request := plaid.NewLinkTokenCreateRequest(
		"jaxon",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		*plaid.NewLinkTokenCreateRequestUser(userId),
	)

	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})
	request.SetOptionalProducts([]plaid.Products{plaid.PRODUCTS_TRANSACTIONS, plaid.PRODUCTS_INVESTMENTS, plaid.PRODUCTS_LIABILITIES})
	request.SetLinkCustomizationName("default")
	// request.SetWebhook()

	// Create the link token
	resp, _, err := plaidClient.PlaidApi.LinkTokenCreate(r.Context()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		slog.Error("Error creating plaid link", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the link token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PlaidLinkResponse{
		LinkToken: resp.GetLinkToken(),
	})
}

type AccessTokenRequest struct {
	PublicToken string `json:"public_token"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func exchangePublicToken(w http.ResponseWriter, r *http.Request) {
	plaidClient := r.Context().Value("plaidClient").(*plaid.APIClient)
	db := r.Context().Value("db").(*sqlx.DB)
	userId := r.Context().Value("userId").(string)

	var req AccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Exchange public token for access token
	exchangeRequest := plaid.NewItemPublicTokenExchangeRequest(req.PublicToken)
	exchangeResp, _, err := plaidClient.PlaidApi.ItemPublicTokenExchange(r.Context()).ItemPublicTokenExchangeRequest(*exchangeRequest).Execute()
	if err != nil {
		slog.Info("Error exchanging public token", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: do not return the access token to the browser
	accessToken := exchangeResp.GetAccessToken()

	slog.Info("Successfully received access token for account link")

	// TODO: need to eventually do this in the background

	accountsRequest := plaid.NewAccountsGetRequest(accessToken)
	resp, _, err := plaidClient.PlaidApi.AccountsGet(r.Context()).AccountsGetRequest(*accountsRequest).Execute()

	slog.Info("", "account_response", resp.GetAccounts()[0])

	item := plaidModels.NewItem(userId, accessToken, *resp.GetItem().InstitutionId.Get(), resp.GetItem().ItemId)
	err = item.Save(db)
	if err != nil {
		slog.Error("Failed to save new Plaid item")
	}

	for _, acct := range resp.GetAccounts() {
		slog.Info("Creating new account")
		account := accountModels.NewAccount(
			acct.GetName(),
			string(acct.GetType()),
			string(acct.GetSubtype()),
			userId,
			*resp.GetItem().InstitutionName.Get(),
			acct.GetMask(),
			item.Id,
		)
		err = account.Save(db)
		if err != nil {
			slog.Error("Failed to save account")
		}
	}

	w.WriteHeader(http.StatusCreated)
}
