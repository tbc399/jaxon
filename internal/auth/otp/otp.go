package otp

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	"log/slog"
	"math/rand/v2"
	"time"
)

type OneTimePass struct {
	Id          string
	MagicToken  string `db:"magic_token"`
	Code        uint32
	Email       string
	Expiry      time.Time
	CreatedAt   time.Time `db:"created_at"`
	Invalidated bool
}

func New(email string) *OneTimePass {
	now := time.Now().UTC()
	return &OneTimePass{
		Id:          shortuuid.New(),
		MagicToken:  shortuuid.New(),
		Code:        rand.Uint32N(9999-1000) + 1000,
		Email:       email,
		Expiry:      now.Add(time.Duration(2) * time.Minute),
		CreatedAt:   now,
		Invalidated: false,
	}
}

func (otp *OneTimePass) IsExpired() bool {
	if otp.Invalidated {
		return true
	}
	return otp.Expiry.Before(time.Now().UTC())
}

func Create(email string, db *sqlx.DB) (*OneTimePass, error) {
	otpass := New(email)
	slog.Info("Creating one time passcode")
	sql := "INSERT INTO otp (id, magic_token, code, email, expiry, created_at, invalidated) VALUES (:id, :magic_token, :code, :email, :expiry, :created_at, :invalidated)"
	tx := db.MustBegin()
	//var err error
	_, err := tx.NamedExec(sql, otpass)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return otpass, nil
}

func Fetch(otpId string, db *sqlx.DB) (*OneTimePass, error) {
	sqls := "SELECT * FROM otp WHERE id = $1"
	otpass := OneTimePass{}
	err := db.Get(&otpass, sqls, otpId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &otpass, nil
}

func FetchByMagicToken(magicToken string, db *sqlx.DB) (*OneTimePass, error) {
	sqls := "SELECT * FROM otp WHERE magic_token = $1"
	otpass := OneTimePass{}
	err := db.Get(&otpass, sqls, magicToken)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &otpass, nil
}

func SendEmail(email string, otpass *OneTimePass, db *sqlx.DB) {
	slog.Info("Sending email")
	magicUrl := fmt.Sprintf("http://localhost:8080/login/magic/%s", otpass.MagicToken)
	slog.Info("magic url", "url", magicUrl)
	slog.Info("sending otp email", "otp", otpass.Id, "email", email)
}
