package service

import (
	"context"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"go_finance/pkg/utils"
	"time"
)

type AuditLogService interface {
	Create(ctx context.Context, entityType, entityID, action, details string) (*domain.AuditLog, error)
	List(ctx context.Context, entityType string) ([]*domain.AuditLog, error)
}

type auditLogService struct {
	repo repository.AuditLogRepository
}

func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Create(ctx context.Context, entityType, entityID, action, details string) (*domain.AuditLog, error) {
	logEntry := &domain.AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}
	err := s.repo.CreateAuditLog(ctx, logEntry)
	if err != nil {
		utils.Logger.Error("Failed to create audit log", "error", err)
		return nil, err
	}
	utils.Logger.Info("Audit log created", "entity_type", entityType, "entity_id", entityID, "action", action)
	return logEntry, nil
}

func (s *auditLogService) List(ctx context.Context, entityType string) ([]*domain.AuditLog, error) {
	logs, err := s.repo.ListAuditLogs(ctx, entityType)
	if err != nil {
		utils.Logger.Error("Failed to list audit logs", "error", err)
		return nil, err
	}
	utils.Logger.Info("Audit logs listed", "entity_type", entityType, "count", len(logs))
	return logs, nil
}
