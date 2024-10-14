package models

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

const (
	AccountTypeCash       = "cash"
	AccountTypeCredit     = "credit"
	AccountTypeInvestment = "investment"
	AccountTypeOther      = "other"
	AccountTypeManual     = "manual"
)

const (
	AccountSubTypeChecking     = "checking"
	AccountSubTypeSavings      = "savings"
	AccountSubTypeCreditCard   = "credit_card"
	AccountSubTypeLineOfCredit = "line_of_credit"
	AccountSubTypeMortgage     = "mortgage"
	AccountSubTypeOther        = "other"
)

type Account struct {
	Id              string
	Name            string
	Type            string         // pulled from AccountTypeXXX constants
	SubType         string         `db:"sub_type"` // pulled from AccountSubTypeXXX constants
	UserId          string         `db:"user_id"`
	InstitutionName string         `db:"institution_name"`
	Last4           string         `db:"last4"`
	StripeId        sql.NullString `db:"stripe_id"`
	LastSync        sql.NullTime   `db:"last_sync"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

func NewAccount(name string, accountType string, subType string, userId string, insitutionName string, last4 string, stripeId string) *Account {
	now := time.Now().UTC()
	return &Account{
		Id:              shortuuid.New(),
		Name:            name,
		Type:            accountType,
		SubType:         subType,
		UserId:          userId,
		InstitutionName: insitutionName,
		Last4:           last4,
		StripeId:        sql.NullString{String: stripeId, Valid: true},
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

func (self *Account) Save(db *sqlx.DB) (*Account, error) {
	sqls := `INSERT INTO accounts 
		(id, name, type, sub_type, user_id, institution_name, last4, stripe_id, last_sync, created_at, updated_at)
		VALUES 
		(:id, :name, :type, :sub_type, :user_id, :institution_name, :last4, :stripe_id, :last_sync, :created_at, :updated_at)`
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
