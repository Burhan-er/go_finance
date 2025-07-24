package service

import (
	"context"
	"go_finance/internal/domain"

)

// RegisterRequest struct for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User *domain.User `json:"user"`
	Message string `json:"message"`
}

// LoginRequest struct for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse struct contains tokens
type LoginResponse struct {
	User        *domain.User `json:"user"`
	AccessToken string       `json:"access_token"`
	Message string `json:"message"`
}

// UserService defines the interface for user-related business logic
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}