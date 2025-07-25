package handler

import (
	"encoding/json"
	"go_finance/internal/service"
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

	req.ID = paramId

	user, err := h.userService.GetUserByID(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var req service.GetAllUsersRequest

	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// defer r.Body.Close()

	users, err := h.userService.GetAllUsers(r.Context(), req)
	if err != nil {
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

	req.ID = paramId

	user, err := h.userService.UpdateUser(r.Context(), req)
	if err != nil {
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

	req.ID = paramId

	user, err := h.userService.DeleteUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(user)

}
