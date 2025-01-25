package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	//"jaxon.app/jaxon/internal/auth/users"
)

type Transaction struct {
	Id           string
	//SourceHash   sql.NullInt64 `db:"source_hash"`
	UserId       string        `db:"user_id"`
	//User         users.User
	AccountId    sql.NullString `db:"account_id"`
	OrigDescription string `db:"orig_description"`
	Description  string
	OrigAmount int `db:"orig_amount"`
	Amount       int            // transaction in cents
	CategoryId   sql.NullString `db:"category_id"`
	//CategoryName string         `db:"category_name"`
	OrigDate         time.Time `db:"orig_date"`
	Date         time.Time
	Notes        sql.NullString
	Hidden       bool
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
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
		Id: "",
		UserId: userId,
		//User: users.User{},
		AccountId: sql.NullString{String: accountId, Valid: false},
		OrigDescription: description,
		Description:  description,
		OrigAmount: amount,
		Amount:       amount,
		CategoryId:   sql.NullString{Valid: false},
		//CategoryName: "",
		OrigDate: date,
		Date:         date,
		Notes:        sql.NullString{Valid: false},
		Hidden:       false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func FetchMany(userId string, db *sqlx.DB) ([]Transaction, error) {
	sqls := "SELECT * FROM transactions WHERE user_id = $1 ORDER BY date DESC"
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
