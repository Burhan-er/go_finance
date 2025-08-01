package postgres

import (
	"context"
	"database/sql"
	"go_finance/internal/domain"
	"time"
)

type auditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *auditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (entity_type, entity_id, action, details, created_at) VALUES ($1, $2, $3, $4, $5)`
	if log.EntityID == "" {
		query = `INSERT INTO audit_logs (entity_type, action, details, created_at) VALUES ($1, $2, $3, $4)`
		_, err := r.db.ExecContext(ctx, query, log.EntityType, log.Action, log.Details, time.Now())
		return err
	}
	_, err := r.db.ExecContext(ctx, query, log.EntityType, log.EntityID, log.Action, log.Details, time.Now())
	return err
}

func (r *auditLogRepository) ListAuditLogs(ctx context.Context, entityType string) ([]*domain.AuditLog, error) {
	query := `SELECT id, entity_type, entity_id, action, details, created_at FROM audit_logs WHERE entity_type = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, entityType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		var log domain.AuditLog
		if err := rows.Scan(&log.ID, &log.EntityType, &log.EntityID, &log.Action, &log.Details, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}
