package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/budget/models/budgets"
)

// setup to run in parallel
func Rollover(ctx context.Context, db *sqlx.DB) {
	slog.Info("Starting rollover job")
	for {
		time.Sleep(time.Second * 60)

		periods, err := budgets.FetchLatestPeriods(db)
		if err != nil {
			continue
		}

		periodUserIds := []string{}

		for _, r := range periods {
			periodUserIds = append(periodUserIds, r.UserId)
		}

		now := time.Now().UTC()
		for _, period := range periods {
			if period.End.Before(now) {
				slog.Info("Starting a budget period rollover", "user", period.UserId, "period", period.Id)
				start := period.Start.AddDate(0, 1, 0)
				end := period.End.AddDate(0, 1, 0)
				newPeriod := budgets.NewBudgetPeriod(period.UserId, start, end)
				previousBudgets, err := period.FetchBudgets(db)
				if err != nil {
					continue
				}

				newBudgets := []budgets.Budget{}
				for _, budget := range previousBudgets {
					newBudgets = append(newBudgets, *budget.RolloverNew(newPeriod))
				}

				// TODO do I need to wrap these in a transactions?
				newPeriod.Save(db)
				budgets.SaveMany(newBudgets, db)

			}
		}

	}
}
