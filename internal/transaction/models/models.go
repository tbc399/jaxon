package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	"jaxon.app/jaxon/internal/budget/models/budgets"
	//"jaxon.app/jaxon/internal/auth/users"
	//"jaxon.app/jaxon/internal/auth/users"
)

type Transaction struct {
	Id string
	// SourceHash   sql.NullInt64 `db:"source_hash"`
	UserId string `db:"user_id"`
	// User         users.User
	AccountId       sql.NullString `db:"account_id"`
	OrigDescription string         `db:"orig_description"`
	Description     string
	OrigAmount      int            `db:"orig_amount"`
	Amount          int            // transaction in cents
	CategoryId      sql.NullString `db:"category_id"`
	// CategoryName string         `db:"category_name"`
	OrigDate  time.Time `db:"orig_date"`
	Date      time.Time
	Notes     sql.NullString
	Hidden    bool
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (self *Transaction) Save(db *sqlx.DB) error {
	modelType := reflect.TypeOf(Transaction{})

	columnNames := []string{}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag, ok := field.Tag.Lookup("db")
		if ok {
			columnNames = append(columnNames, tag)
		} else {
			columnNames = append(columnNames, strings.ToLower(field.Name))
		}
	}

	var sqls string
	if self.Id == "" {
		self.Id = shortuuid.New()
		sqls = fmt.Sprintf(
			"INSERT INTO transactions (%s) VALUES (:%s)",
			strings.Join(columnNames, ", "),
			strings.Join(columnNames, ", :"),
		)
	} else {
		sqls = fmt.Sprintf(
			"UPDATE transactions SET (%s) = (:%s) WHERE id = '%s'",
			strings.Join(columnNames, ", "),
			strings.Join(columnNames, ", :"),
			self.Id,
		)
	}

	slog.Info("Executing sql", "sql", sqls)
	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, self)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func NewTransaction(userId, accountId, description string, date time.Time, amount int) *Transaction {
	now := time.Now().UTC()
	return &Transaction{
		Id:     "",
		UserId: userId,
		// User: users.User{},
		AccountId:       sql.NullString{String: accountId, Valid: false},
		OrigDescription: description,
		Description:     description,
		OrigAmount:      amount,
		Amount:          amount,
		CategoryId:      sql.NullString{Valid: false},
		// CategoryName: "",
		OrigDate:  date,
		Date:      date,
		Notes:     sql.NullString{Valid: false},
		Hidden:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type TransactionView struct {
	Id           string
	Description  string
	Amount       int            // transaction in cents
	CategoryId   sql.NullString `db:"category_id"`
	CategoryName sql.NullString `db:"category_name"`
	Date         time.Time
	Notes        sql.NullString
	Hidden       bool
}

func FetchMany(userId string, db *sqlx.DB) ([]TransactionView, error) {
	sqls := `SELECT t.id, t.description, t.amount, t.category_id, c.name as category_name, t.date, t.notes, t.hidden FROM transactions AS t LEFT JOIN categories AS c ON t.category_id = c.id 
		WHERE t.user_id = $1 ORDER BY t.date DESC`
	slog.Info("Executing sql", "sql", sqls)
	transactions := []TransactionView{}
	err := db.Select(&transactions, sqls, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch transactions", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch transactions", "user", userId, "error", err.Error())
		return nil, err
	}
	return transactions, nil
}

func Fetch(id string, db *sqlx.DB) (*Transaction, error) {
	sqls := "SELECT * FROM transactions WHERE id = $1"
	transaction := new(Transaction)
	err := db.Get(transaction, sqls, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch transaction", "id", id, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch transaction", "id", id, "error", err.Error())
		return nil, err
	}
	return transaction, nil
}

func SumInPeriod(userId string, period *budgets.BudgetPeriod, db *sqlx.DB) (int64, error) {
	sqls := `SELECT COALESCE(SUM(amount), 0) FROM transactions AS t LEFT JOIN categories AS c ON t.category_id = c.id WHERE t.user_id = $1 AND t.date BETWEEN $2 AND $3 AND c.type = 'expense'`
	slog.Info("Executing sql", "sql", sqls)
	var amount int64
	err := db.Get(&amount, sqls, userId, period.Start, period.End)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to sum transactions for period", "user", userId, "period", period, "error", err.Error())
			return 0, nil
		}
		slog.Error("Failed to sum transactions for period", "user", userId, "period", period, "error", err.Error())
		return 0, err
	}
	slog.Info(strconv.FormatInt(amount, 10))
	return amount, nil
}

func CreateMany(transactions []Transaction, db *sqlx.DB) error {
	if len(transactions) == 0 {
		return nil
	}

	modelType := reflect.TypeOf(Transaction{})

	// names to exclude from writing to the db
	exclude := []string{"User", "CategoryName", "Account"}

	columnNames := []string{}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if slices.Contains(exclude, field.Name) {
			continue
		}
		tag, ok := field.Tag.Lookup("db")
		if ok {
			columnNames = append(columnNames, tag)
		} else {
			columnNames = append(columnNames, strings.ToLower(field.Name))
		}
	}

	sqls := fmt.Sprintf(
		"INSERT INTO transactions (%s) VALUES (:%s)",
		strings.Join(columnNames, ", "),
		strings.Join(columnNames, ", :"),
	)
	slog.Info("Executing sql", "sql", sqls)

	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, transactions)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
