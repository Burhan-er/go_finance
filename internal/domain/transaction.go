package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionQueryOption interface{}

type Limit int
type Offset int
type StatusType string
type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)

const (
	Pending   StatusType = "pending"
	Completed StatusType = "completed"
	Failed    StatusType = "failed"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	ToUserID    string          `json:"to_user_id"`
	Type        TransactionType `json:"type"`
	Status      StatusType      `json:"status"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
