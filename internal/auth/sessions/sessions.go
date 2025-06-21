package sessions

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

type Session struct {
	Id            string
	Expiry        time.Time
	OtpId         string         `db:"otp_id"` // the otp this session was generated from
	UserId        string         `db:"user_id"`
	DeviceId      sql.NullString `db:"device_id"`
	DeviceTrusted sql.NullBool   `db:"device_trusted"`
	Invalidated   bool
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (self *Session) IsExpired() bool {
	return self.Expiry.Before(time.Now().UTC())
}

func New(userId string, otpId string) *Session {
	now := time.Now().UTC()
	return &Session{
		Id:          shortuuid.New(),
		Expiry:      now.Add(time.Duration(30) * time.Hour * 24), // 30 days
		OtpId:       otpId,
		UserId:      userId,
		Invalidated: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (session *Session) Save(db *sqlx.DB) (*Session, error) {
	sql := `INSERT INTO sessions (id, expiry, otp_id, user_id, device_id, device_trusted, invalidated, created_at, updated_at) VALUES (:id, :expiry, :otp_id, :user_id, :device_id, :device_trusted, :invalidated, :created_at, :updated_at)`

	slog.Info("Executing sql", "sql", sql)
	tx := db.MustBegin()
	var err error
	_, err = tx.NamedExec(sql, session)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func Fetch(id string, db *sqlx.DB) (*Session, error) {
	sqls := "SELECT * FROM sessions WHERE id = $1"
	session := Session{}
	err := db.Get(&session, sqls, id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &session, nil
}

func FetchByOtpId(otpId string, db *sqlx.DB) (*Session, error) {
	sqls := "SELECT * FROM sessions WHERE otp_id = $1 LIMIT 1"
	session := Session{}
	err := db.Get(&session, sqls, otpId)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &session, nil
}
