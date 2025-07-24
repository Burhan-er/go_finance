package domain

import "time"

type Balance struct {
	UserID        string    `json:"user_id"`
	Amount        int       `json:"amount"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}
