package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"jaxon.app/jaxon/internal/auth"
	budgets "jaxon.app/jaxon/internal/budget/routes"
	dashboard "jaxon.app/jaxon/internal/dashboard/routes"
	"jaxon.app/jaxon/internal/middleware"
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

	var staticServer = http.FileServer(http.Dir("web/"))
	router.Handle("/static/", http.StripPrefix("/static", staticServer))
	var scriptServer = http.FileServer(http.Dir("node_modules/"))
	router.Handle("/modules/", http.StripPrefix("/modules", scriptServer))

	chain := middleware.Chain(middleware.Logging)

	// Auth routes
	authRouter := auth.Router()
	budgetRouter := budgets.Router()
	dashboardRouter := dashboard.Router()

	router.Handle("/dashboard/", http.StripPrefix("/dashboard", dashboardRouter))
	router.Handle("/login/", http.StripPrefix("/login", authRouter))
	router.Handle("/budgets/", http.StripPrefix("/budgets", budgetRouter))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var server = http.Server{
		Addr:    ":8080",
		Handler: chain(router),
		BaseContext: func(listner net.Listener) context.Context {
			return context.WithValue(ctx, "db", db)
		},
	}

	log.Println("Starting server on port :8080")
	server.ListenAndServe()
}
