package services

import (
	"log/slog"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/account/models/accounts"
	"jaxon.app/jaxon/internal/auth/users"
	"jaxon.app/jaxon/internal/budget/models/budgets"
)

func rollover(db *sqlx.DB) {

	for {
		time.Sleep(time.Second * 30)
		now := time.Now().UTC()
		activeUsers, err := users.FetchMany(db)
		if err != nil {
			continue
		}


		rollovers, err := budgets.FetchRolloversByMonth(now.Year(), now.Month(), db)
		if err != nil {
			continue
		}

		for _, user := range(activeUsers) {
			if slices.Contains(rollovers, user) {

			}
			
		}

	}

}
