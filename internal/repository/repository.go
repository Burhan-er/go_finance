


	//Struct ile tanımlamalar yapılabilir



package repository

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"
)

// UserRepository defines the methods for interacting with user data.
type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUserByID(ctx context.Context, id string) (*domain.User, error)
	DeleteUserByID(ctx context.Context, id string) (*domain.User, error)
}

// TransactionRepository defines the methods for interacting with transaction data.
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *domain.Transaction) error
	GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID string) ([]domain.Transaction, error)


	//update Transaction status eklencek


}

// BalanceRepository defines the methods for interacting with balance data.
type BalanceRepository interface {


	//Create Balance eklenecek


	GetBalanceByUserID(ctx context.Context, userID string) (*domain.Balance, error)
	UpdateUserBalance(ctx context.Context, tx *sql.Tx, userID string, newAmount int) error
}
