package service

import (
	"context"
	"go_finance/internal/domain"
	"time"

	"github.com/shopspring/decimal"
)

// RegisterRequest struct for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
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
	Message     string       `json:"message"`
}
type GetAllUsersRequest struct{}

type GetAllUsersResponse struct {
	Users []*domain.User `json:"users"`
}

type GetUserByIdRequest struct {
	ID string
}

type GetUserByIdResponse struct {
	User *domain.User `json:"user"`
}

type PutUserByIdRequest struct {
	ID       string
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type PutUserByIdResponse struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
}

type DeleteUserByIdRequest struct {
	ID string
}

type DeleteUserByIdResponse struct {
	Message string `json:"message"`
}

type PostTransactionCreditRequest struct {
	ToUserID    string                 `json:"to_user_id"`
	FromUserID  string                 `json:"from_user_id"`
	Type        domain.TransactionType `json:"type"`
	Amount      decimal.Decimal        `json:"amount"`
	Description string                 `json:"description"`
}

type PostTransactionCreditResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type PostTransactionDebitRequest struct {
	ToUserID    string                 `json:"to_user_id"`
	FromUserID  string                 `json:"from_user_id"`
	Type        domain.TransactionType `json:"type"`
	Amount      decimal.Decimal        `json:"amount"`
	Description string                 `json:"description"`
}

type PostTransactionDebitResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type PostTransactionTransferRequest struct {
	FromUserID  string          `json:"from_user_id"`
	ToUserID    string          `json:"to_user_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
}

type PostTransactionTransferResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type GetTransactionHistoryRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetTransactionHistoryResponse struct {
	Transactions []*domain.Transaction `json:"transactions"`
}

type GetTransactionByIdRequest struct {
	ID string
}

type GetTransactionByIdResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}

type GetBalanceHistoricalRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetBalanceHistoricalResponse struct {
	UserID  string            `json:"user_id"`
	History []*domain.Balance `json:"history"`
}

type GetBalanceCurrentRequest struct {
	UserID string
}

type GetBalanceCurrentResponse struct {
	Balance *domain.Balance `json:"balance"`
}

type GetBalanceAtTimeRequest struct {
	Timestamp time.Time
}

type GetBalanceAtTimeResponse struct {
	Balance *domain.Balance `json:"balance"`
}

// UserService defines the interface for user-related business logic
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	GetAllUsers(ctx context.Context, req GetAllUsersRequest) (*GetAllUsersResponse, error)
	GetUserByID(ctx context.Context, req *GetUserByIdRequest) (*GetUserByIdResponse, error)
	UpdateUser(ctx context.Context, req PutUserByIdRequest) (*PutUserByIdResponse, error)
	DeleteUser(ctx context.Context, req DeleteUserByIdRequest) (*DeleteUserByIdResponse, error)
}

// TransactionService defines the interface for transaction-related business logic
type TransactionService interface {
	Credit(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error)
	Debit(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error)
	Transfer(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error)
	GetHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error)
	GetByID(ctx context.Context, req GetTransactionByIdRequest) (*GetTransactionByIdResponse, error)
}

// BalanceService defines the interface for balance-related business logic
type BalanceService interface {
	GetCurrent(ctx context.Context, req GetBalanceCurrentRequest) (*GetBalanceCurrentResponse, error)
	GetHistorical(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error)
	GetAtTime(ctx context.Context, req GetBalanceAtTimeRequest) (*GetBalanceAtTimeResponse, error)
}
