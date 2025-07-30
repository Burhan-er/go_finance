package domain

import "time"

type AuditLog struct {
	ID         string     `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   string     `json:"entity_id"`
	Action     string    `json:"action"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}
