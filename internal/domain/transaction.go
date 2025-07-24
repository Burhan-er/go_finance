package domain

import "time"

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)

type StatusType string

const (
	Pending   StatusType = "pending"
	Completed StatusType = "completed"
	Failed    StatusType = "failed"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Type        TransactionType `json:"type"`
	Status      StatusType      `json:"status"`
	Amount      int             `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}
