package handler

import (
	"encoding/json"
	"go_finance/internal/service"
	"go_finance/pkg/utils" // Logger için import eklendi
	"net/http"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(us service.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Warn("failed to decode register request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := h.userService.Register(r.Context(), req)
	if err != nil {
		utils.Logger.Error("user registration failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Warn("failed to decode login request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	res, err := h.userService.Login(r.Context(), req)
	if err != nil {
		utils.Logger.Warn("user login attempt failed", "email", req.Email, "error", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}