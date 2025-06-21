package plaid

import (
	"log/slog"
	"os"

	"github.com/plaid/plaid-go/v37/plaid"
)

func NewConfiguredPlaidClient() *plaid.APIClient {
	// Configure Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", os.Getenv("PLAID_CLIENT_ID"))
	configuration.AddDefaultHeader("PLAID-SECRET", os.Getenv("PLAID_SECRET"))
	configuration.Debug = true

	// Set the environment
	var env plaid.Environment
	switch os.Getenv("PLAID_ENV") {
	case "sandbox":
		slog.Info("setting plaid env for sandbox")
		env = plaid.Sandbox
	case "production":
		slog.Info("setting plaid env for production")
		env = plaid.Production
	default:
		slog.Info("setting plaid env to default of sandbox")
		env = plaid.Sandbox
	}
	configuration.UseEnvironment(env)

	return plaid.NewAPIClient(configuration)
}
