package domain_transaction

import (
	"context"
	"strings"
	"time"

	"ExpenseTracker-Backend/internal/utils"
)

type Service interface {
	Create(ctx context.Context, userID int64, req CreateTransactionRequest) (*TransactionResponse, error)
	GetByID(ctx context.Context, id, userID int64) (*TransactionResponse, error)
	Update(ctx context.Context, id, userID int64, req UpdateTransactionRequest) (*TransactionResponse, error)
	SoftDelete(ctx context.Context, id, userID int64) error
	List(ctx context.Context, userID int64, req TransactionListRequest) (*TransactionListResponse, error)
	GetInfo(ctx context.Context, userID int64, req TransactionInfoRequest) (*TransactionInfoResponse, error)
	GetAnalytics(ctx context.Context, userID int64) (*TransactionAnalyticsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID int64, req CreateTransactionRequest) (*TransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	description := strings.TrimSpace(req.Description)

	tx := &Transaction{
		UserID:      userID,
		Amount:      req.Amount,
		Category:    utils.TransactionCategory(req.Category),
		Type:        utils.TransactionType(req.Type),
		Description: description,
		Status:      utils.Active,
		Date:        req.Date,
	}

	err := s.repo.Create(ctx, tx)
	if err != nil {
		return nil, err
	}

	return toResponse(tx), nil
}

func (s *service) GetByID(ctx context.Context, id, userID int64) (*TransactionResponse, error) {
	tx, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toResponse(tx), nil
}

func (s *service) Update(ctx context.Context, id, userID int64, req UpdateTransactionRequest) (*TransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	tx, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if req.Amount != nil {
		tx.Amount = *req.Amount
	}
	if req.Category != nil {
		tx.Category = utils.TransactionCategory(*req.Category)
	}
	if req.Type != nil {
		tx.Type = utils.TransactionType(*req.Type)
	}
	if req.Description != nil {
		tx.Description = strings.TrimSpace(*req.Description)
	}
	if req.Date != nil {
		tx.Date = *req.Date
	}

	err = s.repo.Update(ctx, tx)
	if err != nil {
		return nil, err
	}

	return toResponse(tx), nil
}

func (s *service) SoftDelete(ctx context.Context, id, userID int64) error {
	_, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, id, userID, int(utils.Deleted))
}

func (s *service) List(ctx context.Context, userID int64, req TransactionListRequest) (*TransactionListResponse, error) {
	transactions, err := s.repo.List(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	var items []TransactionResponse
	for _, tx := range transactions {
		items = append(items, *toResponse(&tx))
	}

	var nextCursor int64 = 0
	hasMore := false

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	if len(items) > 0 {
		if len(items) >= limit {
			nextCursor = items[len(items)-1].ID
			hasMore = true
		}
	}

	return &TransactionListResponse{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	}, nil
}

func (s *service) GetInfo(ctx context.Context, userID int64, req TransactionInfoRequest) (*TransactionInfoResponse, error) {
	return s.repo.GetInfo(ctx, userID, req)
}

func (s *service) GetAnalytics(ctx context.Context, userID int64) (*TransactionAnalyticsResponse, error) {
	return s.repo.GetAnalytics(ctx, userID)
}

func toResponse(tx *Transaction) *TransactionResponse {
	amountFloat, _ := tx.Amount.Float64()
	return &TransactionResponse{
		ID:          tx.ID,
		Amount:      amountFloat,
		Category:    int(tx.Category),
		Type:        int(tx.Type),
		Description: tx.Description,
		Date:        tx.Date.Format(time.RFC3339),
		Status:      int(tx.Status),
	}
}
