package domain_transaction

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// Request
type CreateTransactionRequest struct {
	Amount      decimal.Decimal `json:"amount" validate:"required"`
	Category    int             `json:"category" validate:"required"`
	Type        int             `json:"type" validate:"required"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
}

func (r *CreateTransactionRequest) Validate() error {
	if r.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return errors.New("amount must be greater than zero")
	}

	if r.Category < 1 || r.Category > 21 {
		return errors.New("invalid category value")
	}

	if r.Type != 1 && r.Type != 2 {
		return errors.New("invalid transaction type")
	}

	if r.Date.IsZero() {
		r.Date = time.Now()
	}

	return nil
}

type UpdateTransactionRequest struct {
	Amount      *decimal.Decimal `json:"amount"`
	Category    *int             `json:"category"`
	Type        *int             `json:"type"`
	Description *string          `json:"description"`
	Date        *time.Time       `json:"date"`
}

func (r *UpdateTransactionRequest) Validate() error {
	if r.Amount != nil && r.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		return errors.New("amount must be greater than zero")
	}
	if r.Category != nil && (*r.Category < 1 || *r.Category > 21) {
		return errors.New("invalid category value")
	}
	if r.Type != nil && (*r.Type != 1 && *r.Type != 2) {
		return errors.New("invalid transaction type")
	}
	return nil
}

type UpdateTransactionStatusRequest struct {
	Status int `json:"status" validate:"required"`
}

type TransactionResponse struct {
	ID          int64   `json:"id"`
	Amount      float64 `json:"amount"`
	Category    int     `json:"category"`
	Type        int     `json:"type"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Status      int     `json:"status"`
}

// Cursor Pagination
type TransactionListRequest struct {
	Category *int   `json:"category"`
	Type     *int   `json:"type"`
	Status   *int   `json:"status"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
	AfterID  int64  `json:"after_id"`
	Limit    int    `json:"limit"`
	Sort     string `json:"sort"` // newest, amount_asc, amount_desc
}

type TransactionListResponse struct {
	Items      []TransactionResponse `json:"items"`
	NextCursor int64                 `json:"next_cursor"`
	HasMore    bool                  `json:"has_more"`
	Limit      int                   `json:"limit"`
}

type CategorySummary struct {
	Category int     `json:"category"`
	Count    int64   `json:"count"`
	Amount   float64 `json:"amount"`
}

type TransactionInfoResponse struct {
	TotalTransactions int64             `json:"total_transactions"`
	TotalExpense      float64           `json:"total_expense"`
	TotalIncome       float64           `json:"total_income"`
	CategorySummary   []CategorySummary `json:"category_summary"`
}
