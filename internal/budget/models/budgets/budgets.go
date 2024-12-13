package budgets

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	//"jaxon.app/jaxon/internal/auth/users"
)

type Budget struct {
	Id string
	UserId string `db:"user_id"`
	//User users.User
	CategoryId string `db:"category_id"`
	//Category 
	Amount int
	Rollover bool
	Frequency string
	Year uint
	Month uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BudgetView struct {
	Id string
	UserId string `db:"user_id"`
	//User users.User
	CategoryId string `db:"category_id"`
	CategoryName string `db:"category_name"`
	//Category 
	Amount int
	Frequency string
	Year uint
	Month uint
	TransactionsTotal int `db:"transactions_total"`
}



func FetchAllByMonth(userId string, year int, month time.Month, db *sqlx.DB) ([]BudgetView, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	sqls := `SELECT budgets.id, budgets.user_id, budgets.category_id, budgets.amount, budgets.frequency, budgets.year, budgets.month,
            categories.name as category_name, SUM(COALESCE(transactions.amount, 0)) as transactions_total
            FROM budgets LEFT JOIN categories ON budgets.category_id = categories.id
            LEFT JOIN transactions ON budgets.category_id = transactions.category_id AND (transactions.date BETWEEN $4 AND $5 OR transactions.date is NULL) 
            WHERE budgets.user_id = $1 
                AND budgets.year = $2 
                AND budgets.month = $3 
                GROUP BY budgets.id, categories.id
			`
	slog.Info("Executing sql", "sql", sqls)
	budgets := []BudgetView{}
	fmt.Println(userId)
	fmt.Println(year)
	fmt.Println(month)
	fmt.Println(startOfMonth)
	fmt.Println(endOfMonth)
	err := db.Select(&budgets, sqls, userId, year, month, startOfMonth, endOfMonth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch transactions", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch transactions", "user", userId, "error", err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Found %d budgets", len(budgets)))
	return budgets, nil
}

type BudgetRollover struct {
	Id string
	UserId string `db:"user_id"`
	Year int
	Month int
	CreatedAt time.Time `db:"created_at"`
}

func (self *BudgetRollover) Save(db *sqlx.DB) (*BudgetRollover, error) {
	sqls := `INSERT INTO budget_rollovers (id, user_id, year, month, created_at) VALUES (:id, :user_id, :year, :month, :created_at)`
	slog.Info("Executing sql", "sql", sqls)
	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, self)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return self, nil
}

func FetchRolloversByMonth(year int, month time.Month, db *sqlx.DB) ([]BudgetRollover, error) {
	sqls := `SELECT * from budget_rollovers WHERE year = $1 AND month = $2`
	slog.Info("Executing sql", "sql", sqls)
	rollovers := []BudgetRollover{}
	err := db.Select(&rollovers, sqls, year, month)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch rollovers", "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch rollovers", "error", err.Error())
		return nil, err
	}
	return rollovers, nil
}



