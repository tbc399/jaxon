package models

import (
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Profile struct {
	Id        string
	First     sql.NullString
	Last      sql.NullString
	UserId    string
	CreatedAt string
	UdpatedAt string
}

func NewProfile() *Profile {
	return &Profile{}
}

func Fetch(userId string, db *sqlx.DB) (*Profile, error) {
	sqls := "SELECT * FROM profiles WHERE id = $1"
	slog.Info("Executing sql", "sql", sqls)
	profile := Profile{}
	err := db.Get(&profile, sqls, userId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &profile, nil
}
