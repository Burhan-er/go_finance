package service

import (
	"context"
	"go_finance/internal/domain"
	"time"
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
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type PutUserByIdResponse struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
}

type DeleteUserByIdRequest struct {
	ID string `json:"id"`
}

type DeleteUserByIdResponse struct {
	Message string `json:"message"`
}

type PostTransactionCreditRequest struct {
	UserID      string `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type PostTransactionCreditResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type PostTransactionDebitRequest struct {
	UserID      string `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type PostTransactionDebitResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type PostTransactionTransferRequest struct {
	FromUserID  string `json:"from_user_id"`
	ToUserID    string `json:"to_user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

type PostTransactionTransferResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

type GetTransactionHistoryRequest struct {
	UserID string `json:"user_id"`
}

type GetTransactionHistoryResponse struct {
	Transactions []*domain.Transaction `json:"transactions"`
}

type GetTransactionByIdRequest struct {
	ID string `json:"id"`
}

type GetTransactionByIdResponse struct {
	Transaction *domain.Transaction `json:"transaction"`
}

type GetBalanceHistoricalRequest struct {
	UserID string `json:"user_id"`
}

type GetBalanceHistoricalResponse struct {
	UserID  string            `json:"user_id"`
	History []*domain.Balance `json:"history"`
}

type GetBalanceCurrentRequest struct {
	UserID string `json:"user_id"`
}

type GetBalanceCurrentResponse struct {
	Balance *domain.Balance `json:"balance"`
}

type GetBalanceAtTimeRequest struct {
	UserID string    `json:"user_id"`
	Time   time.Time `json:"time"`
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
	CreditTransaction(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error)
	DebitTransaction(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error)
	TransferTransaction(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error)
	GetTransactionHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error)
	GetTransactionByID(ctx context.Context, req GetTransactionByIdRequest) (*GetTransactionByIdResponse, error)
}

// BalanceService defines the interface for balance-related business logic
type BalanceService interface {
	GetCurrentBalance(ctx context.Context, req GetBalanceCurrentRequest) (*GetBalanceCurrentResponse, error)
	GetHistoricalBalance(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error)
	GetBalanceAtTime(ctx context.Context, req GetBalanceAtTimeRequest) (*GetBalanceAtTimeResponse, error)
}
