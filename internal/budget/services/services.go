package services

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"jaxon.app/jaxon/internal/budget/models/budgets"
	"jaxon.app/jaxon/internal/budget/models/categories"
	"jaxon.app/jaxon/internal/transaction/models"
)

func CreateDefaultCategories(userId string) {
	slog.Info("Creating default categories", "user", userId)
	expenseCategories := []string{
		// Home
		"Home",
		"Mortgage",
		"Rent",
		"Home Insurance",
		"Rental Insurance",
		"HOA Dues",
		"Home Supplies",
		"Home Maintenance",
		"Flood Insurance",
		"Furniture",
		// Auto
		"Auto",
		"Car Payment",
		"Car Insurance",
		"Car Maintenance",
		"Gas & Fuel",
		"Car Wash",
		"Toll",
		"Inspection & Registration",
		"Public Transportation",
		"Rideshare",
		"Parking",
		// Food
		"Food",
		"Groceries",
		"Restaurants",
		"Fast Food",
		"Coffee Shop",
		// Education
		"Education",
		"Tuition",
		"Student Loan",
		"Books & Supplies",
		// Cash & ATM
		"Cash & ATM",
		// Charity
		"Charity & Donations",
		// Entertainment
		"Entertainment",
		"Movies",
		"Family Night",
		"Date Night",
		// Financial
		"Financial",
		"Life Insurance",
		"Retirement Savings",
		"Investments",
		// Fitness
		"Fitness",
		"Gym Membership",
		"Personal Training",
		// Health
		"Health",
		"Suppliments",
		"Doctor",
		"Dentist",
		"Health Insurance",
		"Health Share",
		"Eyecare",
		"Pharmacy",
		// Gifts
		"Gifts",
		// Kids
		"Kids",
		"Child Care",
		"Child Clothing",
		"Babysitter",
		"Diapers",
		"Formula",
		"Toys",
		// Personal Care
		"Personal Care",
		"Salon",
		"Barber",
		"Spa",
		"Laundry",
		// Savings
		"Savings",
		"Emergency Fund",
		"Vacation Fund",
		"Car Fund",
		// Pets
		"Pets",
		"Veterinary",
		"Pet Food",
		"Pet Grooming",
		"Pet Boarding",
		// Shopping
		"Shopping",
		"Electronics",
		"Clothing",
		"Books",
		// Travel
		"Travel",
		"Airfare",
		"Rental Cars",
		"Hotels",
		// Utilities
		"Utilities",
		"Electricity",
		"Water",
		"Gas",
		"Internet & Cable",
		"Phone",
		"Trash",
	}
	incomeCategories := []string{
		//"Income",
		"Paycheck",
		"Bonus",
		"Tax Refund",
		"Earned Interest",
		"Dividends",
		"Rental Income",
	}

	cats := []categories.Category{}

	for _, name := range expenseCategories {
		cats = append(cats, *categories.NewCategory(name, categories.ExpenseCategoryType, userId))
	}

	for _, name := range incomeCategories {
		cats = append(cats, *categories.NewCategory(name, categories.IncomeCategoryType, userId))
	}

	//categories.CreateMany(userId, )
	/*
	   cats.extend(
	       [
	           Category(name=name, type=CategoryType.income, user_id=message.user_id)
	           for name in income_categories
	       ]
	   )

	   async with db.pool.acquire() as connection:
	       await Category.create_many(cats, connection)
	*/
}

type BudgetOverview struct {
	ExpectedIncome int64
	CurrentIncome  int64
	ExpectedSpend  int64
	CurrentSpend   int64
}

func GetBudgetOverview(userId string, period *budgets.BudgetPeriod, db *sqlx.DB) (*BudgetOverview, error) {
	var expecedIncome int64 = 9_900
	var currentIncome int64 = 2_000

	currentSpend, err := models.SumInPeriod(userId, period, db)
	if err != nil {
		return nil, err
	}

	currentSpend = currentSpend / -100

	expectedSpend, err := period.SumBudgets(userId, db)
	if err != nil {
		return nil, err
	}

	return &BudgetOverview{
		ExpectedIncome: expecedIncome,
		CurrentIncome:  currentIncome,
		ExpectedSpend:  expectedSpend,
		CurrentSpend:   currentSpend,
	}, nil
}
