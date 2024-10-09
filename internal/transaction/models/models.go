package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/auth/users"
)

type Transaction struct {
	Id string
	SourceHash sql.NullInt64 `db:"source_hash"`
	UserId string `db:"user_id"`
	User users.User
	AccountId string`db:"account_id"`
	Description string
	Amount int // transaction in cents
	CategoryId sql.NullString `db:"category_id"`
	CategoryName string `db:"category_name"`
	Date time.Time
	Notes sql.NullString
	Hidden bool
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
