package repository

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"

	"github.com/shopspring/decimal"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUserByID(ctx context.Context, id string) (*domain.User, error)
	DeleteUserByID(ctx context.Context, id string) (*domain.User, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) (string,error)
	GetTransactionsByUserID(ctx context.Context, userID string, opts ...domain.TransactionQueryOption) ([]*domain.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, tx *sql.Tx, id string, status domain.StatusType) error
	GetTranscaptionByID(ctx context.Context, id string)(*domain.Transaction, error)

}

type BalanceRepository interface {
	GetBalanceByUserID(ctx context.Context, userID string) (*domain.Balance, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, userID string, amount decimal.Decimal) error
	CreateBalance(ctx context.Context, balance *domain.Balance) error
}

type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, log *domain.AuditLog) error
	ListAuditLogs(ctx context.Context, entityType string) ([]*domain.AuditLog, error)
}