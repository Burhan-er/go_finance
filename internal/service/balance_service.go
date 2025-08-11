package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go_finance/internal/api/middleware"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"go_finance/pkg/cache"
	"go_finance/pkg/utils"
	"time"
)

type balanceService struct {
	balanceRepo     repository.BalanceRepository
	auditLogService AuditLogService
	cache           *cache.Cache
}

func NewBalanceService(repo repository.BalanceRepository, auditLogService AuditLogService, cache *cache.Cache) BalanceService {
	return &balanceService{
		balanceRepo:     repo,
		auditLogService: auditLogService,
		cache:           cache,
	}
}

func (s *balanceService) GetCurrent(ctx context.Context, req GetBalanceCurrentRequest) (*GetBalanceCurrentResponse, error) {

	s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_attempt", "Get current balance attempt")

	if req.UserID == "" {
		utils.Logger.Warn("User ID required for current balance retrieval")
		return nil, errors.New("user id required")
	}
	// CACHİNG
	cacheKey := fmt.Sprintf("balance:%s", req.UserID)
	val, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		var bal domain.Balance
		if jsonErr := json.Unmarshal([]byte(val), &bal); jsonErr == nil {
			return &GetBalanceCurrentResponse{
				Balance: &bal,
			}, nil
		}

	}

	balance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Warn("User not found while retrieving balance", "userID", req.UserID)
			s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_failed", "User not found for balance")
			return nil, errors.New("user not found")
		}
		utils.Logger.Error("Error retrieving balance", "userID", req.UserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_failed", "Internal server error for balance")
		return nil, errors.New("internal server error")
	}

	s.auditLogService.Create(ctx, "balance", req.UserID, "get_current_success", "Current balance retrieved successfully")
	utils.Logger.Info("Current balance retrieved successfully", "userID", req.UserID)

	if err := s.cache.Set(ctx, cacheKey, balance, time.Duration(10*time.Second)); err != nil {
		utils.Logger.Error("GetCurrentBalanceOperation Cache Error", "userID", req.UserID)
	}

	response := &GetBalanceCurrentResponse{
		Balance: &domain.Balance{
			UserID:        balance.UserID,
			Amount:        balance.Amount,
			LastUpdatedAt: balance.LastUpdatedAt,
		},
	}

	return response, nil
}

func (s *balanceService) GetHistorical(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error) {
	s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_attempt", "Get historical balance attempt")
	if req.UserID == "" {
		utils.Logger.Warn("User ID required for historical balance retrieval")
		return nil, errors.New("user id required")
	}
	if req.StartDate == "" || req.EndDate == "" {
		utils.Logger.Warn("Date range missing for historical balance", "start", req.StartDate, "end", req.EndDate)
		return nil, errors.New("start date and end date are required")
	}

	history, err := s.balanceRepo.GetBalanceHistoryByUserID(ctx, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		utils.Logger.Error("Error retrieving balance history", "userID", req.UserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_failed", "Internal server error for balance history")
		return nil, errors.New("internal server error")
	}

	s.auditLogService.Create(ctx, "balance", req.UserID, "get_historical_success", "Historical balance retrieved successfully")
	utils.Logger.Info("Historical balance retrieved successfully", "userID", req.UserID)

	response := &GetBalanceHistoricalResponse{
		History: history,
	}
	return response, nil
}

func (s *balanceService) GetAtTime(ctx context.Context, req GetBalanceAtTimeRequest) (*GetBalanceAtTimeResponse, error) {
	s.auditLogService.Create(ctx, "balance", "", "get_at_time_attempt", "Get balance at specific time attempt")

	if req.UserID != ctx.Value(middleware.UserIDKey) && ctx.Value(middleware.UserRoleKey) != domain.AdminRole {
		utils.Logger.Warn("Unauthorized access to balance at time", "requestUser", req.UserID)
		return nil, fmt.Errorf("you dont have access")
	}

	if req.Timestamp.IsZero() {
		utils.Logger.Warn("Timestamp is required for balance at time query")
		return nil, errors.New("timestamp is required")
	}

	balance, err := s.balanceRepo.GetBalanceAtTime(ctx, req.UserID, req.Timestamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Warn("Balance not found at given time", "userID", req.UserID, "timestamp", req.Timestamp)
			s.auditLogService.Create(ctx, "balance", "", "get_at_time_failed", "Balance not found at given time")
			return nil, errors.New("balance not found at given time")
		}
		utils.Logger.Error("Error retrieving balance at time", "userID", req.UserID, "timestamp", req.Timestamp, "error", err)
		s.auditLogService.Create(ctx, "balance", "", "get_at_time_failed", "Internal server error for balance at time")
		return nil, errors.New("internal server error")
	}

	s.auditLogService.Create(ctx, "balance", "", "get_at_time_success", "Balance at given time retrieved successfully")
	utils.Logger.Info("Balance at time retrieved successfully", "userID", req.UserID, "timestamp", req.Timestamp)

	response := &GetBalanceAtTimeResponse{
		Balance: balance,
	}
	return response, nil
}
