package services

import (
	"encoding/csv"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/transaction/models"
)

func parseCents(amount string) (int, error) {
	p, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return 0, err
	}

	return int(p * 100), nil
}
func UploadTransactions(file io.Reader, userId string, db *sqlx.DB) {

	reader := csv.NewReader(file)

	records := [][]string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Failed to read line")
			continue
		}
		records = append(records, record)
	}

	transactions := []models.Transaction{}
	header := map[string]int{
		"date": 0,
		"description": 0,
		"amount": 0,
	}

	for i, record := range records {
		if i == 0 {
			// establish header column placement
			for j, token := range record {
				header[strings.ToLower(token)] = j
			}
			continue
		}

		date, err := time.Parse("01/02/2006", record[header["date"]])

		if err != nil {
			slog.Error("Failed to parse date", "row", i)
			return
		}

		cents, err := parseCents(record[header["amount"]])

		if err != nil {
			slog.Error("Failed to parse amount", "row", i)
			return
		}

		transactions = append(transactions,
			*models.NewTransaction(
				userId,
				"",
				record[header["description"]],
				date,
				cents))
	}

	models.CreateMany(transactions, db)

}
