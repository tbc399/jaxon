package plaid

import "github.com/plaid/plaid-go/v23/plaid"
import "os"

func NewConfiguredPlaidClient() *plaid.APIClient {
	// Configure Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", os.Getenv("PLAID_CLIENT_ID"))
	configuration.AddDefaultHeader("PLAID-SECRET", os.Getenv("PLAID_SECRET"))

	// Set the environment
	var env plaid.Environment
	switch os.Getenv("PLAID_ENV") {
	case "sandbox":
		env = plaid.Sandbox
	case "development":
		env = plaid.Development
	case "production":
		env = plaid.Production
	default:
		env = plaid.Sandbox
	}
	configuration.UseEnvironment(env)

	return plaid.NewAPIClient(configuration)
}
