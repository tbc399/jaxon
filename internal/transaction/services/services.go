package services

import (
	"encoding/csv"
	//	"fmt"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	"jaxon.app/jaxon/internal/transaction/models"
)

func parseCents(amount string) (int, error) {
	p, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return int(p * 100), nil
}

func UploadTransactions(file io.Reader, userId, accountId string, db *sqlx.DB) {
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
		"date":        0,
		"description": 0,
		"amount":      0,
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

		newTransaction := *models.NewTransaction(
			userId,
			accountId,
			record[header["description"]],
			date,
			cents)
		newTransaction.Id = shortuuid.New()
		transactions = append(transactions, newTransaction)

	}

	models.CreateMany(transactions, db)
}

func GroupTransactionsByDate(transactions []models.TransactionView) []interface{} {
	sort.Slice(transactions, func(lhs, rhs int) bool {
		return transactions[lhs].Date.After(transactions[rhs].Date)
	})
	groups := make([]interface{}, 0)

	t := transactions[0].Date
	currentDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	groups = append(groups, currentDate)
	for _, t := range transactions {
		if !(t.Date.Year() == currentDate.Year() && t.Date.Month() == currentDate.Month() && t.Date.Day() == currentDate.Day()) {
			currentDate = time.Date(t.Date.Year(), t.Date.Month(), t.Date.Day(), 0, 0, 0, 0, t.Date.Location())
			groups = append(groups, currentDate)
		}
		groups = append(groups, t)
	}
	return groups
}
