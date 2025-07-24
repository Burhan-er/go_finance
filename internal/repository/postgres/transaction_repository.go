package postgres

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"
	"time"
)

type transactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

// CreateTransaction creates a new transaction in the database
func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	query := `INSERT INTO transactions (user_id, type, status, amount, description, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	transaction.CreatedAt = time.Now()
	if transaction.Status == "" {
		transaction.Status = domain.Pending
	}

	err := r.db.QueryRowContext(ctx, query, transaction.UserID, transaction.Type, transaction.Status, transaction.Amount, transaction.Description, transaction.CreatedAt).Scan(&transaction.ID)
	return err
}

// GetTransactionByID retrieves a single transaction by its ID
func (r *transactionRepository) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `SELECT id, user_id, type, status, amount, description, created_at FROM transactions WHERE id = $1`

	var transaction domain.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Status,
		&transaction.Amount, &transaction.Description, &transaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUserID retrieves all transactions for a given user
func (r *transactionRepository) GetTransactionsByUserID(ctx context.Context, userID string) ([]*domain.Transaction, error) {
	query := `SELECT id, user_id, type, status, amount, description, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Type, &transaction.Status, &transaction.Amount, &transaction.Description, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransactionStatus updates the status of a specific transaction
func (r *transactionRepository) UpdateTransactionStatus(ctx context.Context, id string, status domain.StatusType) error {
	query := `UPDATE transactions SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}