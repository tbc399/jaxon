package categories

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
	//"jaxon.app/jaxon/internal/templates"
)

const (
	IncomeCategoryType  = "income"
	ExpenseCategoryType = "expense"
)

type Category struct {
	// templates.Selectable
	Id        string
	Name      string
	Type      string    // from CategorType constants
	UserId    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (self Category) GetId() string {
	return self.Id
}

func (self Category) GetName() string {
	return self.Name
}

func (self Category) Save(db *sqlx.DB) error {
	modelType := reflect.TypeOf(Category{})

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

	sqls := fmt.Sprintf(
		"INSERT INTO categories (%s) VALUES (:%s)",
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

func NewCategory(name, catType, userId string) *Category {
	now := time.Now().UTC()
	return &Category{
		Id:        shortuuid.New(),
		Name:      name,
		Type:      catType,
		UserId:    userId,
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

func Fetch(id string, db *sqlx.DB) (*Category, error) {
	sqls := "SELECT * FROM categories WHERE id = $1"
	slog.Info("Executing sql", "sql", sqls, "category_id", id)
	category := new(Category)
	err := db.Get(category, sqls, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch category", "id", id, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch category", "id", id, "error", err.Error())
		return nil, err
	}
	return category, nil
}

func FetchByName(name, userId string, db *sqlx.DB) (*Category, error) {
	sqls := "SELECT * FROM categories WHERE name = $1 AND user_id = $2"
	slog.Info("Executing sql", "sql", sqls, "name", name, "user_id", userId)
	category := new(Category)
	err := db.Get(category, sqls, name, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch category", "name", name, "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch category", "name", name, "user", userId, "error", err.Error())
		return nil, err
	}
	return category, nil
}
