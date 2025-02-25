package users

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

type User struct {
	Id               string
	Email            string
	StripeCustomerId sql.NullString `db:"stripe_customer_id"`
	Active           bool
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

func New(email string) *User {
	now := time.Now().UTC()
	return &User{
		Id:        shortuuid.New(),
		Email:     email,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (user *User) Save(db *sqlx.DB) (*User, error) {
	sqls := `INSERT INTO users (id, email, stripe_customer_id, active, created_at, updated_at) VALUES (:id, :email, :stripe_customer_id, :active, :created_at, :updated_at)`
	slog.Info("Executing sql", "sql", sqls)
	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, user)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FetchByEmail(email string, db *sqlx.DB) (*User, error) {
	sqls := "SELECT * FROM users WHERE email = $1"
	slog.Info("Executing sql", "sql", sqls)
	user := User{}
	err := db.Get(&user, sqls, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch user", "email", email, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch user", "email", email, "error", err.Error())
		return nil, err
	}
	return &user, nil
}

func FetchMany(db *sqlx.DB) ([]User, error) {
	sqls := "SELECT * FROM users WHERE active = true"
	slog.Info("Executing sql", "sql", sqls)
	users := []User{}
	err := db.Get(&users, sqls)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch users", "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch user", "error", err.Error())
		return nil, err
	}
	return users, nil
}
