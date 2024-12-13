package services

import (
	"log/slog"

	"jaxon.app/jaxon/internal/budget/models/categories"
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
		cats := append(cats, *categories.NewCategory(name, categories.ExpenseCategoryType, userId))
	}

	for _, name := range incomeCategories {
		cats := append(cats, *categories.NewCategory(name, categories.IncomeCategoryType, userId))
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

/*
async def create_default_categories(message: UserCreated):
    logger.info(f"Creating default categories for user {message.user_id}")
    expense_categories = [
        # Home
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
        # Auto
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
        # Food
        "Food",
        "Groceries",
        "Restaurants",
        "Fast Food",
        "Coffee Shop",
        # Education
        "Education",
        "Tuition",
        "Student Loan",
        "Books & Supplies",
        # Cash & ATM
        "Cash & ATM",
        # Charity
        "Charity & Donations",
        # Entertainment
        "Entertainment",
        "Movies",
        "Family Night",
        "Date Night",
        # Financial
        "Financial",
        "Life Insurance",
        "Retirement Savings",
        "Investments",
        # Fitness
        "Fitness",
        "Gym Membership",
        "Personal Training",
        # Health
        "Health",
        "Suppliments",
        "Doctor",
        "Dentist",
        "Health Insurance",
        "Health Share",
        "Eyecare",
        "Pharmacy",
        # Gifts
        "Gifts",
        # Kids
        "Kids",
        "Child Care",
        "Child Clothing",
        "Babysitter",
        "Diapers",
        "Formula",
        "Toys",
        # Personal Care
        "Personal Care",
        "Salon",
        "Barber",
        "Spa",
        "Laundry",
        # Savings
        "Savings",
        "Emergency Fund",
        "Vacation Fund",
        "Car Fund",
        # Pets
        "Pets",
        "Veterinary",
        "Pet Food",
        "Pet Grooming",
        "Pet Boarding",
        # Shopping
        "Shopping",
        "Electronics",
        "Clothing",
        "Books",
        # Travel
        "Travel",
        "Airfare",
        "Rental Cars",
        "Hotels",
        # Utilities
        "Utilities",
        "Electricity",
        "Water",
        "Gas",
        "Internet & Cable",
        "Phone",
        "Trash",
    ]
    income_categories = [
        "Income",
        "Paycheck",
        "Bonus",
        "Tax Refund",
        "Earned Interest",
        "Dividends",
        "Rental Income",
    ]
    cats = [
        Category(name=name, type=CategoryType.expense, user_id=message.user_id)
        for name in expense_categories
    ]
    cats.extend(
        [
            Category(name=name, type=CategoryType.income, user_id=message.user_id)
            for name in income_categories
        ]
    )
    async with db.pool.acquire() as connection:
        await Category.create_many(cats, connection)
*/
