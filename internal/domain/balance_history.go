package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type BalanceHistory struct {
	ID         string
	UserID     string
	Amount     decimal.Decimal
	RecordedAt time.Time
}
