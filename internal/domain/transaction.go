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
type TransactionJobType string


const (
	Pending   StatusType = "pending"
	Completed StatusType = "completed"
	Failed    StatusType = "failed"
)
const (
	Credit TransactionJobType = "credit"
	Debit TransactionJobType = "debit"
	Transfer TransactionJobType = "transfer"
)


type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	ToUserID    string          `json:"to_user_id"`
	Type        TransactionType `json:"type"`
	Status      StatusType      `json:"status"`
	Amount      decimal.Decimal `json:"amount"`
	CreatedAt   time.Time       `json:"created_at"`
}
