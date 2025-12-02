package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

type Account struct {
	Id              string
	Name            string
	Type            string         // pulled from AccountTypeXXX constants
	SubType         string         `db:"sub_type"` // pulled from AccountSubTypeXXX constants
	UserId          string         `db:"user_id"`
	InstitutionName string         `db:"institution_name"`
	Last4           string         `db:"last4"`
	PlaidItemId     sql.NullString `db:"plaid_item_id"`
	LastSync        sql.NullTime   `db:"last_sync"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

func NewAccount(name string, accountType string, subType string, userId string, institutionName string, last4 string, plaidItemId string) *Account {
	now := time.Now().UTC()
	return &Account{
		Id:              shortuuid.New(),
		Name:            name,
		Type:            accountType,
		SubType:         subType,
		UserId:          userId,
		InstitutionName: institutionName,
		Last4:           last4,
		PlaidItemId:     sql.NullString{String: plaidItemId, Valid: true},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func FetchAll(userId string, db *sqlx.DB) ([]Account, error) {
	sqls := "SELECT * FROM accounts WHERE user_id = $1"
	slog.Info("Executing sql", "sql", sqls)
	accounts := []Account{}
	err := db.Select(&accounts, sqls, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch accounts", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch accounts", "user", userId, "error", err.Error())
		return nil, err
	}
	return accounts, nil
}

func (self *Account) Save(db *sqlx.DB) error {
	modelType := reflect.TypeOf(Account{})

	columnNames := []string{}
	for i := range modelType.NumField() {
		field := modelType.Field(i)
		tag, ok := field.Tag.Lookup("db")
		if ok {
			columnNames = append(columnNames, tag)
		} else {
			columnNames = append(columnNames, strings.ToLower(field.Name))
		}
	}

	sqls := fmt.Sprintf(
		"INSERT INTO accounts (%s) VALUES (:%s)",
		strings.Join(columnNames, ", "),
		strings.Join(columnNames, ", :"),
	)
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

func (self *Account) LastSyncDisplay() string {
	timeSinceSync := time.Now().UTC().Sub(self.LastSync.Time)
	if timeSinceSync < time.Duration(5)*time.Minute {
		return "a moment ago"
	} else if timeSinceSync < time.Duration(24)*time.Hour {
		return "today"
	} else if timeSinceSync < time.Duration(5)*24*time.Hour {
		return "a week ago"
	} else {
		return "over a month ago"
	}
}
