package utils

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

type TransactionStatus int

const (
	Active   TransactionStatus = 1
	Deleted  TransactionStatus = 2
	Archived TransactionStatus = 3
)
