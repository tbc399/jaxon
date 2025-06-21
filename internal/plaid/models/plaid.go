package plaid

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

type Item struct {
	Id            string
	UserId        string    `db:"user_id"`
	AccessToken   string    `db:"access_token"`
	InstitutionId string    `db:"institution_id"`
	ItemId        string    `db:"item_id"` // Plaid's id for this Item
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewItem(userId, accessToken, institutionId, itemId string) *Item {
	now := time.Now().UTC()
	return &Item{
		Id:            shortuuid.New(),
		UserId:        userId,
		AccessToken:   accessToken,
		InstitutionId: institutionId,
		ItemId:        itemId,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (self *Item) Save(db *sqlx.DB) error {
	modelType := reflect.TypeOf(Item{})

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
		"INSERT INTO plaid_items (%s) VALUES (:%s)",
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

func Fetch() {
}
