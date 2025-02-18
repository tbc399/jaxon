package assets

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
)

const (
	AssetTypeVehicle  = "vehicle"
	AssetTypeProperty = "property"
	AssetTypeOther    = "other"
)

type Asset struct {
	Id        string
	Name      string
	Type      string    // pulled from AssetTypeXXX constants
	UserId    string    `db:"user_id"`
	LastSync  time.Time `db:"last_sync"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewAsset(name string, assetType string, userId string, lastSync time.Time) *Asset {
	now := time.Now().UTC()
	return &Asset{
		Id:        shortuuid.New(),
		Name:      name,
		Type:      assetType,
		UserId:    userId,
		LastSync:  lastSync,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func FetchAll(userId string, db *sqlx.DB) ([]Asset, error) {
	sqls := "SELECT * FROM assets WHERE user_id = $1"
	assets := []Asset{}
	err := db.Select(&assets, sqls, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch assets", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch assets", "user", userId, "error", err.Error())
		return nil, err
	}
	return assets, nil
}
