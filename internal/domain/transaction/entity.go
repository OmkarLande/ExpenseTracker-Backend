package domain_transaction

import (
	"time"

	"ExpenseTracker-Backend/internal/utils"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          int64                      `db:"id"`
	UserID      int64                      `db:"user_id"`
	Amount      decimal.Decimal            `db:"amount"`
	Category    utils.TransactionCategory  `db:"category"`
	Type        utils.TransactionType      `db:"type"`
	Description string                     `db:"description"`
	Status      utils.TransactionStatus    `db:"status"`
	Date        time.Time                  `db:"date"`
	CreatedAt   time.Time                  `db:"created_at"`
	UpdatedAt   time.Time                  `db:"updated_at"`
}