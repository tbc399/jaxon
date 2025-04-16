package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	accounts "jaxon.app/jaxon/internal/account/routes"
	"jaxon.app/jaxon/internal/auth"
	budgets "jaxon.app/jaxon/internal/budget/routes"
	budgetServices "jaxon.app/jaxon/internal/budget/services"
	dashboard "jaxon.app/jaxon/internal/dashboard/routes"
	"jaxon.app/jaxon/internal/middleware"
	"jaxon.app/jaxon/internal/plaid"
	profile "jaxon.app/jaxon/internal/profile/routes"
	transactions "jaxon.app/jaxon/internal/transaction/routes"
)

func main() {

	// TODO: Need to setup a connection pool
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var router = http.NewServeMux()

	// static files
	var staticServer = http.FileServer(http.Dir("web/"))
	router.Handle("/static/", http.StripPrefix("/static", staticServer))
	var scriptServer = http.FileServer(http.Dir("node_modules/"))
	router.Handle("/modules/", http.StripPrefix("/modules", scriptServer))

	middlewares := middleware.Chain(middleware.LogRequest)

	authRouter := auth.Router()
	//budgetRouter := middleware.EnsureAuth(budgets.Router())
	//dashboardRouter := middleware.EnsureAuth(dashboard.Router())
	//transactionRouter := middleware.EnsureAuth(transactions.Router())

	appRouter := http.NewServeMux()

	//auth.AddRoutes(router)
	budgets.AddRoutes(appRouter)
	dashboard.AddRoutes(appRouter)
	transactions.AddRoutes(appRouter)
	accounts.AddRoutes(appRouter)
	profile.AddRoutes(appRouter)

	// Auth protected routes
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))
	router.Handle("/", middleware.EnsureAuth(appRouter))
	//router.Handle("/dashboard/", http.StripPrefix("", dashboardRouter))
	//router.Handle("/budgets/", http.StripPrefix("", budgetRouter))
	//router.Handle("/transactions/", http.StripPrefix("", transactionRouter))

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "plaidClient", plaid.NewConfiguredPlaidClient())
	defer cancel()

	var server = http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
		BaseContext: func(listner net.Listener) context.Context {
			return context.WithValue(ctx, "db", db)
		},
	}

	go budgetServices.Rollover(ctx, db)

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
