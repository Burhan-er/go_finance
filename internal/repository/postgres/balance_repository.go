package postgres

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"
	"time"

	"github.com/shopspring/decimal"
)

type balanceRepository struct {
	db *sql.DB
}

// NewBalanceRepository creates a new balance repository
func NewBalanceRepository(db *sql.DB) *balanceRepository {
	return &balanceRepository{db: db}
}

// CreateBalance creates an initial balance for a user
func (r *balanceRepository) CreateBalance(ctx context.Context, balance *domain.Balance) error {
	query := `INSERT INTO balances (user_id, amount, last_updated_at)
			  VALUES ($1, $2, $3)`

	balance.LastUpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, balance.UserID, balance.Amount, balance.LastUpdatedAt)
	return err
}

// GetBalanceByUserID retrieves the current balance for a given user
func (r *balanceRepository) GetBalanceByUserID(ctx context.Context, userID string) (*domain.Balance, error) {
	query := `SELECT user_id, amount, last_updated_at FROM balances WHERE user_id = $1`

	var balance domain.Balance
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.UserID, &balance.Amount, &balance.LastUpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// UpdateBalance updates a user's balance
// This often happens within a database transaction to ensure atomicity
func (r *balanceRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, userID string, amount decimal.Decimal) error {
	query := `UPDATE balances SET amount = amount + $1, last_updated_at = $2 WHERE user_id = $3`

	_, err := tx.ExecContext(ctx, query, amount, time.Now(), userID)
	return err
}
