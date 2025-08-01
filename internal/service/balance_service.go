package service

import (
	"context"
	"database/sql"
	"errors"
	"go_finance/internal/api/middleware"
	"go_finance/internal/domain"
	"go_finance/internal/repository"

	"github.com/shopspring/decimal"
)

type balanceService struct {
	balanceRepo     repository.BalanceRepository
	auditLogService AuditLogService
}

func NewBalanceService(repo repository.BalanceRepository, auditLogService AuditLogService) BalanceService {
	return &balanceService{
		balanceRepo:     repo,
		auditLogService: auditLogService,
	}
}

func (s *balanceService) GetCurrent(ctx context.Context, req GetBalanceCurrentRequest) (*GetBalanceCurrentResponse, error) {
	s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_attempt", "Get current balance attempt")
	if req.UserID == "" {
		return nil, errors.New("user id required")
	}
	balance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_failed", "User not found for balance")
			return nil, errors.New("user not found")
		}
		s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_failed", "Internal server error for balance")
		return nil, errors.New("internal server error")
	}
	s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_success", "Current balance retrieved successfully")

	response := &GetBalanceCurrentResponse{
		Balance: &domain.Balance{
			UserID:        balance.UserID,
			Amount:        balance.Amount,
			LastUpdatedAt: balance.LastUpdatedAt,
		},
	}

	return response, nil
}
func (s *balanceService) GetHistorical(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error){
	return nil,nil
}
// func (s *balanceService) GetHistorical(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error) {
// 	s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_attempt", "Get historical balance attempt")
// 	logs, err := s.auditLogService.List(ctx, "balance")
// 	if err != nil {
// 		s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_failed", "Failed to retrieve historical balance logs")
// 		return nil, errors.New("failed to retrieve historical balance logs")
// 	}
// 	s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_success", "Historical balance logs retrieved successfully")
// 	return &GetBalanceHistoricalResponse{
// 		History: convertAuditLogsToBalances(logs),
// 	}, nil
// }

func (s *balanceService) GetAtTime(ctx context.Context, req GetBalanceAtTimeRequest) (*GetBalanceAtTimeResponse, error) {
	var userID string
	userIDValue := ctx.Value(middleware.UserIDKey)
	if v, ok := userIDValue.(string); ok {
		userID = v
	} else {
		return nil, errors.New("user id not found in context")
	}
	s.auditLogService.Create(ctx, "balance", userID, "get_at_time_attempt", "Get balance at specific time attempt")
	logs, err := s.auditLogService.List(ctx, "balance")
	if err != nil {
		s.auditLogService.Create(ctx, "balance", userID, "get_at_time_failed", "Failed to retrieve balance logs for specific time")
		return nil, errors.New("failed to retrieve balance logs for specific time")
	}
	var balanceAtTime *domain.Balance
	amount, _ := decimal.NewFromString("0.00")
	for _, log := range logs {
		if log.CreatedAt.Equal(req.Timestamp) || log.CreatedAt.Before(req.Timestamp) {
			balanceAtTime = &domain.Balance{
				UserID:        userID,
				Amount:        amount,
				LastUpdatedAt: log.CreatedAt,
			}
		}
	}
	s.auditLogService.Create(ctx, "balance", userID, "get_at_time_success", "Balance at specific time retrieved successfully")
	return &GetBalanceAtTimeResponse{
		Balance: balanceAtTime,
	}, nil
}
