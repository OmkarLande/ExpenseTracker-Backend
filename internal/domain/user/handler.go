package domain_user

import (
	"net/http"

	"ExpenseTracker-Backend/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req UpdateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.service.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) PresignProfilePicture(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req PresignProfilePictureRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	resp, err := h.service.PresignProfilePicture(r.Context(), userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) CompleteProfilePicture(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req CompleteProfilePictureRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.service.CompleteProfilePicture(r.Context(), userID, req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}