package domain_transaction

import (
	"net/http"
	"strconv"

	"ExpenseTracker-Backend/internal/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req CreateTransactionRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	tx, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, tx)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid transaction id"})
		return
	}

	tx, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, tx)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid transaction id"})
		return
	}

	var req UpdateTransactionRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	tx, err := h.service.Update(r.Context(), id, userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, tx)
}

func (h *Handler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid transaction id"})
		return
	}

	var req UpdateTransactionStatusRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Soft delete matches enum status Deleted=2
	if req.Status != 2 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid soft delete status operation"})
		return
	}

	err = h.service.SoftDelete(r.Context(), id, userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "transaction soft deleted successfully"})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	// Parse query params
	q := r.URL.Query()

	var req TransactionListRequest

	if catStr := q.Get("category"); catStr != "" {
		if val, err := strconv.Atoi(catStr); err == nil {
			req.Category = &val
		}
	}

	if typeStr := q.Get("type"); typeStr != "" {
		if val, err := strconv.Atoi(typeStr); err == nil {
			req.Type = &val
		}
	}

	if statusStr := q.Get("status"); statusStr != "" {
		if val, err := strconv.Atoi(statusStr); err == nil {
			req.Status = &val
		}
	}

	req.FromDate = q.Get("from_date")
	req.ToDate = q.Get("to_date")

	if afterStr := q.Get("after_id"); afterStr != "" {
		if val, err := strconv.ParseInt(afterStr, 10, 64); err == nil {
			req.AfterID = val
		}
	}

	if limitStr := q.Get("limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			req.Limit = val
		}
	}

	req.Sort = q.Get("sort")

	resp, err := h.service.List(r.Context(), userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	resp, err := h.service.GetInfo(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}
