package handler

import (
	"encoding/json"
	"go_finance/internal/service"
	"go_finance/pkg/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	var req service.GetUserByIdRequest
	paramId := chi.URLParam(r, "id")

	if paramId == "" {
		utils.Logger.Warn("user ID is missing from URL parameters")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	req.ID = paramId

	user, err := h.userService.GetUserByID(r.Context(), &req)
	if err != nil {
		utils.Logger.Error("failed to get user by id", "user_id", paramId, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var req service.GetAllUsersRequest

	users, err := h.userService.GetAllUsers(r.Context(), req)
	if err != nil {
		utils.Logger.Error("failed to list users", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req service.PutUserByIdRequest
	paramId := chi.URLParam(r, "id")

	if paramId == "" {
		utils.Logger.Warn("user ID is missing from URL parameters for update")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Warn("failed to read request.body for user update", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.ID = paramId 

	user, err := h.userService.UpdateUser(r.Context(), req)
	if err != nil {
		utils.Logger.Error("failed to update user", "user_id", paramId, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req service.DeleteUserByIdRequest
	paramId := chi.URLParam(r, "id")

	if paramId == "" {
		utils.Logger.Warn("user ID is missing from URL parameters for delete")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	req.ID = paramId

	_, err := h.userService.DeleteUser(r.Context(), req)
	if err != nil {
		utils.Logger.Error("failed to delete user", "user_id", paramId, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
