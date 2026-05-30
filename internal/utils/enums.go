package utils

import "strings"

type TransactionType int

const (
	Expense TransactionType = 1
	Income  TransactionType = 2
)

type TransactionCategory int

const (
	Entertainment  TransactionCategory = 1
	Health         TransactionCategory = 2
	Clothing       TransactionCategory = 3
	Dining         TransactionCategory = 4
	Groceries      TransactionCategory = 5
	Utilities      TransactionCategory = 6
	Rent           TransactionCategory = 7
	Transportation TransactionCategory = 8
	Education      TransactionCategory = 9
	Gifts          TransactionCategory = 10
	Travel         TransactionCategory = 11
	Food           TransactionCategory = 12
	Bills          TransactionCategory = 13
	Miscellaneous  TransactionCategory = 14
	Salary         TransactionCategory = 15
	Freelancing    TransactionCategory = 16
	SideHustle     TransactionCategory = 17
	Investments    TransactionCategory = 18
	Refunds        TransactionCategory = 19
	Cashback       TransactionCategory = 20
	Others         TransactionCategory = 21
)

var CategoryNames = map[TransactionCategory]string{
	Entertainment:  "Entertainment",
	Health:         "Health",
	Clothing:       "Clothing",
	Dining:         "Dining",
	Groceries:      "Groceries",
	Utilities:      "Utilities",
	Rent:           "Rent",
	Transportation: "Transportation",
	Education:      "Education",
	Gifts:          "Gifts",
	Travel:         "Travel",
	Food:           "Food",
	Bills:          "Bills",
	Miscellaneous:  "Miscellaneous",
	Salary:         "Salary",
	Freelancing:    "Freelancing",
	SideHustle:     "SideHustle",
	Investments:    "Investments",
	Refunds:        "Refunds",
	Cashback:       "Cashback",
	Others:         "Others",
}

// GetCategoryIDsByNameSearch searches category names (case-insensitive substring match)
// and returns their corresponding IDs.
func GetCategoryIDsByNameSearch(q string) []int {
	if q == "" {
		return nil
	}
	q = strings.ToLower(q)
	var ids []int
	for cat, name := range CategoryNames {
		if strings.Contains(strings.ToLower(name), q) {
			ids = append(ids, int(cat))
		}
	}
	return ids
}

type TransactionStatus int

const (
	Active   TransactionStatus = 1
	Deleted  TransactionStatus = 2
	Archived TransactionStatus = 3
)

