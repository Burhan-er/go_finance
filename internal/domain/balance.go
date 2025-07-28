package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	UserID        string          `json:"user_id"`
	Amount        decimal.Decimal `json:"amount"`
	LastUpdatedAt time.Time       `json:"last_updated_at"`
}
