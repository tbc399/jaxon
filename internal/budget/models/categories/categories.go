package categories

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	"jaxon.app/jaxon/internal/templates"
)

const (
	IncomeCategoryType = "income"
	ExpenseCategoryType = "expense"
)

type Category struct {
	templates.Selectable
    Id string
    Name string
    Type string
	UserId string `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

func (self Category) GetId() string {
	return self.Id
}

func (self Category) GetName() string {
	return self.Name
}

func NewCategory(name, catType, userId string) *Category {
	now := time.Now().UTC()
	return &Category{
		Id: shortuuid.New(),
		Name: name,
		Type: catType,
		UserId: userId,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func FetchAll(userId string, db *sqlx.DB) ([]Category, error) {
	sqls := "SELECT * FROM categories WHERE user_id = $1"
	slog.Info("Executing sql", "sql", sqls)
	categories := []Category{}
	err := db.Select(&categories, sqls, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch categories", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch categories", "user", userId, "error", err.Error())
		return nil, err
	}
	return categories, nil
}
