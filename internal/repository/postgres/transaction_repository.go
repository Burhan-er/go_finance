package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go_finance/internal/domain"
	"time"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, db *sql.DB, transaction *domain.Transaction) (string, error) {
	query := `INSERT INTO transactions (from_user_id, to_user_id, type, status, amount, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	transaction.CreatedAt = time.Now()
	if transaction.Status == "" {
		transaction.Status = domain.Pending
	}
	var insertedID string

	err := db.QueryRowContext(ctx, query,
		transaction.UserID,
		transaction.ToUserID,
		transaction.Type,
		transaction.Status,
		transaction.Amount,
		transaction.CreatedAt,
	).Scan(&insertedID)

	return insertedID, err
}

func (r *transactionRepository) GetTransactionsByUserID(ctx context.Context, id string, opts ...domain.TransactionQueryOption) ([]*domain.Transaction, error) {
	limit := -1
	offset := 0
	var ttype domain.StatusType
	var hasType bool

	for _, opt := range opts {
		switch v := opt.(type) {
		case domain.Limit:
			limit = int(v)
		case domain.Offset:
			offset = int(v)
		case domain.StatusType:
			ttype = domain.StatusType(v)
			hasType = true
		}
	}
	query := `SELECT id, from_user_id, to_user_id, type, status, amount, created_at 
	          FROM transactions WHERE from_user_id = $1`
	args := []interface{}{id}
	argIndex := 2

	if hasType {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, ttype)
		argIndex++
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
		argIndex++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction

	for rows.Next() {
		var t domain.Transaction
		err := rows.Scan(
			&t.ID, &t.UserID, &t.ToUserID, &t.Type, &t.Status,
			&t.Amount, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) GetTranscaptionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	query := `SELECT id, from_user_id, to_user_id, type, status, amount, created_at
	          FROM transactions WHERE id = $1`

	var t domain.Transaction
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.UserID,
		&t.ToUserID,
		&t.Type,
		&t.Status,
		&t.Amount,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *transactionRepository) UpdateTransactionStatus(ctx context.Context, tx *sql.Tx, db *sql.DB, id string, status domain.StatusType) error {
    query := `UPDATE transactions SET status = $1 WHERE id = $2`

    if tx != nil {
        _, err := tx.ExecContext(ctx, query, status, id)
        return err
    }

    if db != nil {
        _, err := db.ExecContext(ctx, query, status, id)
        return err
    }

    return fmt.Errorf("neither tx nor db provided")
}

