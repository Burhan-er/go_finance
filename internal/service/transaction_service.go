package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_finance/internal/api/middleware"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"go_finance/pkg/utils"
	"time"
)

type transactionService struct {
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
	db              *sql.DB
	auditLogService AuditLogService
	processor       *TransactionProcessor
}

func NewTransactionService(
	tr repository.TransactionRepository,
	br repository.BalanceRepository,
	db *sql.DB,
	auditLogService AuditLogService,
	processor *TransactionProcessor,
) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		balanceRepo:     br,
		db:              db,
		auditLogService: auditLogService,
		processor:       processor,
	}
}

func (s *transactionService) Credit(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error) {
	s.auditLogService.Create(ctx, "transaction", "", "credit_queued", fmt.Sprintf("Credit transaction queued from %s, amount: %v", req.ToUserID, req.Amount))

	newTransaction := &domain.Transaction{
		UserID:    req.FromUserID,
		ToUserID:  req.ToUserID,
		Type:      domain.TransactionType(domain.Credit),
		Amount:    req.Amount,
		Status:    domain.Pending, 
		CreatedAt: time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, s.db, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create initial credit transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	job := Job{
		TransactionID: newTransaction.ID,
		Type:          newTransaction.Type,
		FromUserID:    newTransaction.UserID,
		ToUserID:      newTransaction.ToUserID,
		Amount:        newTransaction.Amount,
	}
	s.processor.SubmitJob(job)

	utils.Logger.Info("Credit transaction queued for processing", "transaction_id", newTransaction.ID)

	resp := &PostTransactionCreditResponse{
		Transaction: newTransaction,
		Message:     "Credit transaction has been accepted for processing.",
	}

	return resp, nil
}

func (s *transactionService) Debit(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error) {
	s.auditLogService.Create(ctx, "transaction", "", "debit_queued", fmt.Sprintf("Debit transaction queued from %s, amount: %v", req.FromUserID, req.Amount))

	currentBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to pre-check user balance for debit", "user_id", req.FromUserID, "error", err)
		return nil, fmt.Errorf("failed to check user balance: %w", err)
	}
	if currentBalance.Amount.Cmp(req.Amount) < 0 {
		utils.Logger.Warn("Insufficient funds on pre-check for debit", "user_id", req.FromUserID)
		return nil, errors.New("insufficient funds")
	}

	newTransaction := &domain.Transaction{
		UserID:    req.FromUserID,
		ToUserID:  req.ToUserID,
		Type:      domain.TransactionType(domain.Debit),
		Amount:    req.Amount,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, s.db, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create initial debit transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	job := Job{
		TransactionID: newTransaction.ID,
		Type:          newTransaction.Type,
		FromUserID:    newTransaction.UserID,
		ToUserID:      newTransaction.ToUserID,
		Amount:        newTransaction.Amount,
	}
	s.processor.SubmitJob(job)

	utils.Logger.Info("Debit transaction queued for processing", "transaction_id", newTransaction.ID)

	resp := &PostTransactionDebitResponse{
		Transaction: newTransaction,
		Message:     "Debit transaction has been accepted for processing.",
	}

	return resp, nil
}

func (s *transactionService) Transfer(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error) {
	s.auditLogService.Create(ctx, "transaction", "", "transfer_queued", fmt.Sprintf("Transfer transaction queued from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))

	senderBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to pre-check sender's balance for transfer", "user_id", req.FromUserID, "error", err)
		return nil, fmt.Errorf("failed to check sender's balance: %w", err)
	}
	if senderBalance.Amount.Cmp(req.Amount) < 0 {
		utils.Logger.Warn("Insufficient funds on pre-check for transfer", "user_id", req.FromUserID)
		return nil, errors.New("insufficient funds for transfer")
	}

	newTransaction := &domain.Transaction{
		UserID:    req.FromUserID,
		ToUserID:  req.ToUserID,
		Type:      domain.TransactionType(domain.Transfer),
		Amount:    req.Amount,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, s.db, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create initial transfer transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	job := Job{
		TransactionID: newTransaction.ID,
		Type:          newTransaction.Type,
		FromUserID:    newTransaction.UserID,
		ToUserID:      newTransaction.ToUserID,
		Amount:        newTransaction.Amount,
	}
	s.processor.SubmitJob(job)

	utils.Logger.Info("Transfer transaction queued for processing", "transaction_id", newTransaction.ID)

	resp := &PostTransactionTransferResponse{
		Transaction: newTransaction,
		Message:     "Transfer has been accepted for processing.",
	}

	return resp, nil
}

func (s *transactionService) GetHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error) {
	if ctx.Value(middleware.UserIDKey) != req.UserID && ctx.Value(middleware.UserRoleKey) != string(domain.AdminRole) {
		utils.Logger.Warn("Unauthorized access attempt for transaction history",
			"requesting_user_id", ctx.Value(middleware.UserIDKey),
			"target_user_id", *req.UserID,
			"role", ctx.Value(middleware.UserRoleKey))
		return nil, errors.New("you do not have permission to view this transaction history")
	}

	var opts []domain.TransactionQueryOption
	if req.Limit != nil {
		opts = append(opts, domain.Limit(*req.Limit))
	}
	if req.Offset != nil {
		opts = append(opts, domain.Offset(*req.Offset))
	}
	if req.Type != nil && *req.Type != "" {
		opts = append(opts, domain.TransactionType(*req.Type))
	}

	transactions, err := s.transactionRepo.GetTransactionsByUserID(ctx, *req.UserID, opts...)
	if err != nil {
		utils.Logger.Error("Failed to get transaction history", "user_id", *req.UserID, "error", err)
		return nil, fmt.Errorf("failed to get transactions by user id: %w", err)
	}

	resp := &GetTransactionHistoryResponse{
		Transactions: transactions,
	}

	return resp, nil
}

func (s *transactionService) GetByID(ctx context.Context, req GetTransactionByIdRequest) (*GetTransactionByIdResponse, error) {
	transaction, err := s.transactionRepo.GetTranscaptionByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Logger.Warn("Transaction not found", "transaction_id", req.ID)
			return nil, fmt.Errorf("transaction with ID %s not found", req.ID)
		}
		utils.Logger.Error("Failed to get transaction by ID", "transaction_id", req.ID, "error", err)
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	if ctx.Value(middleware.UserIDKey) != transaction.UserID && ctx.Value(middleware.UserRoleKey) != domain.AdminRole {
		utils.Logger.Warn("User dont have access for get transaction by id", "userID", ctx.Value(middleware.UserIDKey), "requestID", req.ID)
		return nil, fmt.Errorf("user has not access to get transaction by ID: %s", req.ID)
	}

	resp := &GetTransactionByIdResponse{
		Transaction: transaction,
	}

	return resp, nil
}
