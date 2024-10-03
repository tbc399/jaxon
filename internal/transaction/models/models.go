package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Id string
	UserId string `db:"user_id"`
	AccountId string`db:"account_id"`
	StripeId string`db:"stripe_id"`
	Completed bool
	CreatedAt time.Time`db:"created_at"`
	UpdatedAt time.Time`db:"updated_at"`
}

func FetchMany(userId string, db *sqlx.DB) (*[]Transaction, error) {
	sqls := "SELECT * FROM transactions WHERE user_id = $1"
	transactions := []Transaction{}
	err := db.Select(&transactions, sqls, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch transactions", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch transactions", "user", userId, "error", err.Error())
		return nil, err
	}
	return &transactions, nil
}
