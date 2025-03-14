package budgets

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
	//"jaxon.app/jaxon/internal/auth/users"
)

type Budget struct {
	Id         string
	PeriodId   string `db:"period_id"`
	UserId     string `db:"user_id"`
	CategoryId string `db:"category_id"`
	Amount     int64
	Rollover   bool
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func NewBudget(periodId, userId, categoryId string, amount int64) *Budget {
	now := time.Now().UTC()
	return &Budget{
		Id:         shortuuid.New(),
		PeriodId:   periodId,
		UserId:     userId,
		CategoryId: categoryId,
		Amount:     amount,
		Rollover:   false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (self *Budget) Save(db *sqlx.DB) error {
	modelType := reflect.TypeOf(Budget{})

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
		"INSERT INTO budgets (%s) VALUES (:%s)",
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

func (self *Budget) RolloverNew(period *BudgetPeriod) *Budget {
	now := time.Now().UTC()
	return &Budget{
		Id:         shortuuid.New(),
		PeriodId:   period.Id,
		UserId:     self.UserId,
		CategoryId: self.CategoryId,
		Amount:     self.Amount,
		Rollover:   self.Rollover,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func SaveMany(budgets []Budget, db *sqlx.DB) error {
	if len(budgets) == 0 {
		return nil
	}

	modelType := reflect.TypeOf(Budget{})

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
		"INSERT INTO budgets (%s) VALUES (:%s)",
		strings.Join(columnNames, ", "),
		strings.Join(columnNames, ", :"),
	)
	slog.Info("Executing sql", "sql", sqls)

	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, budgets)
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

func FetchBudget(id, userId string, db *sqlx.DB) (*Budget, error) {
	sqls := `SELECT * FROM budgets WHERE id = $1 AND user_id = $2`
	slog.Info("Executing sql", "sql", sqls)
	budget := new(Budget)
	err := db.Get(budget, sqls, id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch budget", "budget_id", id, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch budgets", "budget_id", id, "error", err.Error())
		return nil, err
	}
	return budget, nil
}

type BudgetView struct {
	Id                string
	UserId            string `db:"user_id"`
	CategoryId        string `db:"category_id"`
	CategoryName      string `db:"category_name"`
	Amount            int64
	TransactionsTotal int `db:"transactions_total"`
}

func FetchBudgetViewsByMonth(userId string, year int, month time.Month, db *sqlx.DB) ([]BudgetView, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	sqls := `SELECT budgets.id, budgets.user_id, budgets.category_id, budgets.amount,
            categories.name as category_name, SUM(COALESCE(transactions.amount, 0)) as transactions_total
            FROM budgets LEFT JOIN budget_periods ON budgets.period_id = budget_periods.id LEFT JOIN categories ON budgets.category_id = categories.id
            LEFT JOIN transactions ON budgets.category_id = transactions.category_id AND (transactions.date BETWEEN $4 AND $5 OR transactions.date is NULL) 
            WHERE budgets.user_id = $1 
                AND EXTRACT(YEAR from budget_periods.start) = $2 
                AND EXTRACT(MONTH from budget_periods.start) = $3 
                GROUP BY budgets.id, categories.id
			`
	slog.Info("Executing sql", "sql", sqls)
	budgets := []BudgetView{}
	err := db.Select(&budgets, sqls, userId, year, month, startOfMonth, endOfMonth)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch budgets", "user", userId, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch budgets", "user", userId, "error", err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Found %d budgets", len(budgets)))
	return budgets, nil
}

type BudgetPeriod struct {
	Id        string
	UserId    string `db:"user_id"`
	Start     time.Time
	End       time.Time
	CreatedAt time.Time `db:"created_at"`
}

func NewBudgetPeriod(userId string, start, end time.Time) *BudgetPeriod {
	now := time.Now().UTC()
	return &BudgetPeriod{
		Id:        shortuuid.New(),
		UserId:    userId,
		Start:     start,
		End:       end,
		CreatedAt: now,
	}
}

func (self *BudgetPeriod) Save(db *sqlx.DB) (*BudgetPeriod, error) {
	sqls := `INSERT INTO budget_periods (id, user_id, start, "end", created_at) VALUES (:id, :user_id, :start, :end, :created_at)`
	slog.Info("Executing sql", "sql", sqls)
	tx := db.MustBegin()
	_, err := tx.NamedExec(sqls, self)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return self, nil
}

func FetchCurrentPeriod(userId string, db *sqlx.DB) (*BudgetPeriod, error) {
	now := time.Now().UTC()
	sqls := `SELECT * FROM budget_periods WHERE user_id = $1 AND budget_periods.end > $2`
	slog.Info("Executing sql", "sql", sqls)
	period := new(BudgetPeriod)
	err := db.Get(period, sqls, userId, now)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch budget period", "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch budget period", "error", err.Error())
		return nil, err
	}
	return period, nil
}

func (self *BudgetPeriod) FetchBudgets(db *sqlx.DB) ([]Budget, error) {
	sqls := `SELECT * FROM budgets WHERE period_id = $1`
	slog.Info("Executing sql", "sql", sqls)
	budgets := []Budget{}
	err := db.Select(&budgets, sqls, self.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch budgets", "user", self.UserId, "period", self.Id, "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch budgets", "user", self.UserId, "period", self.Id, "error", err.Error())
		return nil, err
	}
	slog.Info(fmt.Sprintf("Found %d budgets", len(budgets)))
	return budgets, nil
}

func (self *BudgetPeriod) SumBudgets(userId string, db *sqlx.DB) (int64, error) {
	sqls := `SELECT COALESCE(SUM(amount),0) FROM budgets WHERE period_id = $1 AND user_id = $2`
	// TODO: do I need a new index?
	slog.Info("Executing sql", "sql", sqls)
	var amount int64
	err := db.Get(&amount, sqls, self.Id, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to sum budget amounts", "user", self.UserId, "period", self.Id, "error", err.Error())
			return 0, nil
		}
		slog.Error("Failed to sum budget amounts", "user", self.UserId, "period", self.Id, "error", err.Error())
		return 0, err
	}
	return amount, nil
}

func FetchLatestPeriods(db *sqlx.DB) ([]BudgetPeriod, error) {
	sqls := `SELECT DISTINCT ON (bp.user_id) bp.* FROM budget_periods AS bp LEFT JOIN users AS u ON bp.user_id = u.id WHERE u.active = true ORDER BY bp.user_id, bp.end DESC`
	slog.Info("Executing sql", "sql", sqls)
	rollovers := []BudgetPeriod{}
	err := db.Select(&rollovers, sqls)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("Failed to fetch budget periods", "error", err.Error())
			return nil, nil
		}
		slog.Error("Failed to fetch budget periods", "error", err.Error())
		return nil, err
	}
	return rollovers, nil
}


