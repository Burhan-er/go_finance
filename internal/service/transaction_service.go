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
}

func NewTransactionService(
	tr repository.TransactionRepository,
	br repository.BalanceRepository,
	db *sql.DB,
	auditLogService AuditLogService,
) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		balanceRepo:     br,
		db:              db,
		auditLogService: auditLogService,
	}
}

func (s *transactionService) Credit(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error) {
	s.auditLogService.Create(ctx, "transaction","" , "credit_attempt", fmt.Sprintf("Credit attempt from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to begin credit transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	newTransaction := &domain.Transaction{
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Status:      domain.Pending,
		CreatedAt:   time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create credit transaction record", "error", err)
		s.auditLogService.Create(ctx, "transaction", "", "credit_failed", "Failed to create transaction record")
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to update sender balance for credit", "user_id", req.FromUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.FromUserID, "credit_balance_failed", "Failed to update sender balance for credit")
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to update receiver balance for credit", "user_id", req.ToUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.ToUserID, "credit_balance_failed", "Failed to update receiver balance for credit")
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update credit transaction status", "transaction_id", newTransaction.ID, "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "credit_status_failed", "Failed to update transaction status")
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit credit transaction", "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "credit_commit_failed", "Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "credit_success", fmt.Sprintf("Credit transaction completed from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	utils.Logger.Info("Credit transaction completed", "transaction_id", newTransaction.ID, "user_id", req.FromUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionCreditResponse{
		Transaction: newTransaction,
		Message:     "Credit transaction completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) Debit(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error) {
	s.auditLogService.Create(ctx, "transaction", "", "debit_attempt", fmt.Sprintf("Debit attempt from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to begin debit transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	currentBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to get user balance for debit", "user_id", req.FromUserID, "error", err)
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	if currentBalance.Amount.Compare(req.Amount) != +1 {
		utils.Logger.Warn("Insufficient funds for debit", "user_id", req.FromUserID, "requested_amount", req.Amount, "current_balance", currentBalance.Amount)
		return nil, errors.New("insufficient funds")
	}

	newTransaction := &domain.Transaction{
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Status:      domain.Pending,
		CreatedAt:   time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create debit transaction record", "error", err)
		s.auditLogService.Create(ctx, "transaction", "", "debit_failed", "Failed to create transaction record")
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to update sender balance for debit", "user_id", req.FromUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.FromUserID, "debit_balance_failed", "Failed to update sender balance for debit")
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to update receiver balance for debit", "user_id", req.ToUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.ToUserID, "debit_balance_failed", "Failed to update receiver balance for debit")
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}
	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update debit transaction status", "transaction_id", newTransaction.ID, "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "debit_status_failed", "Failed to update transaction status")
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}
	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit debit transaction", "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "debit_commit_failed", "Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "debit_success", fmt.Sprintf("Debit transaction completed from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	utils.Logger.Info("Debit transaction completed", "transaction_id", newTransaction.ID, "user_id", req.FromUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionDebitResponse{
		Transaction: newTransaction,
		Message:     "Debit transaction completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) Transfer(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error) {
	s.auditLogService.Create(ctx, "transaction", "", "transfer_attempt", fmt.Sprintf("Transfer attempt from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to begin transfer transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	senderBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to get sender's balance for transfer", "user_id", req.FromUserID, "error", err)
		return nil, fmt.Errorf("failed to get sender's balance: %w", err)
	}

	if senderBalance.Amount.Cmp(req.Amount) == -1 {
		utils.Logger.Warn("Insufficient funds for transfer", "user_id", req.FromUserID, "requested_amount", req.Amount, "current_balance", senderBalance.Amount)
		return nil, errors.New("insufficient funds for transfer")
	}

	newTransaction := &domain.Transaction{
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        domain.Transfer,
		Amount:      req.Amount,
		Status:      domain.Pending,
		CreatedAt:   time.Now(),
	}

	insertedID, err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction)
	if err != nil {
		utils.Logger.Error("Failed to create transfer transaction record", "error", err)
		s.auditLogService.Create(ctx, "transaction", "", "transfer_failed", "Failed to create transaction record")
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}
	newTransaction.ID = insertedID

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to debit from sender for transfer", "user_id", req.FromUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.FromUserID, "transfer_balance_failed", "Failed to debit from sender for transfer")
		return nil, fmt.Errorf("failed to debit from sender: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to credit to receiver for transfer", "user_id", req.ToUserID, "error", err)
		s.auditLogService.Create(ctx, "balance", req.ToUserID, "transfer_balance_failed", "Failed to credit to receiver for transfer")
		return nil, fmt.Errorf("failed to credit to receiver: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update transfer transaction status", "transaction_id", newTransaction.ID, "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "transfer_status_failed", "Failed to update transaction status")
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit transfer transaction", "error", err)
		s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "transfer_commit_failed", "Failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	s.auditLogService.Create(ctx, "transaction", newTransaction.ID, "transfer_success", fmt.Sprintf("Transfer transaction completed from %s to %s, amount: %v", req.FromUserID, req.ToUserID, req.Amount))
	utils.Logger.Info("Transfer transaction completed", "transaction_id", newTransaction.ID, "from_user_id", req.FromUserID, "to_user_id", req.ToUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionTransferResponse{
		Transaction: newTransaction,
		Message:     "Transfer completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) GetHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error) {
	if ctx.Value(middleware.UserIDKey) != *req.UserID && ctx.Value(middleware.UserRoleKey) != domain.AdminRole {
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
		utils.Logger.Warn("User dont have access", "userID", ctx.Value(middleware.UserIDKey), "requestID", req.ID)
		return nil, fmt.Errorf("user has not access to get transaction by ID: %s", req.ID)
	}

	resp := &GetTransactionByIdResponse{
		Transaction: transaction,
	}

	return resp, nil
}
