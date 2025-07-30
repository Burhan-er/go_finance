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
}

func NewTransactionService(
	tr repository.TransactionRepository,
	br repository.BalanceRepository,
	db *sql.DB,
) TransactionService {
	return &transactionService{
		transactionRepo: tr,
		balanceRepo:     br,
		db:              db,
	}
}

func (s *transactionService) Credit(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error) {
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
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction); err != nil {
		utils.Logger.Error("Failed to create credit transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to update sender balance for credit", "error", err)
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to update receiver balance for credit", "error", err)
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update credit transaction status", "error", err)
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit credit transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	utils.Logger.Info("Credit transaction completed", "transaction_id", newTransaction.ID, "user_id", req.FromUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionCreditResponse{
		Transaction: newTransaction,
		Message:     "Credit transaction completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) Debit(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to begin debit transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	currentBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to get user balance for debit", "error", err)
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	if currentBalance.Amount.Compare(req.Amount) != +1 {
		utils.Logger.Warn("Insufficient funds for debit", "user_id", req.FromUserID, "amount", req.Amount)
		return nil, errors.New("insufficient funds")
	}

	newTransaction := &domain.Transaction{
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Status:      domain.Pending,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction); err != nil {
		utils.Logger.Error("Failed to create debit transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to update sender balance for debit", "error", err)
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to update receiver balance for debit", "error", err)
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update debit transaction status", "error", err)
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit debit transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	utils.Logger.Info("Debit transaction completed", "transaction_id", newTransaction.ID, "user_id", req.FromUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionDebitResponse{
		Transaction: newTransaction,
		Message:     "Debit transaction completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) Transfer(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to begin transfer transaction", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	senderBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		utils.Logger.Error("Failed to get sender's balance for transfer", "error", err)
		return nil, fmt.Errorf("failed to get sender's balance: %w", err)
	}

	if senderBalance.Amount.Cmp(req.Amount) == -1 {
		utils.Logger.Warn("Insufficient funds for transfer", "user_id", req.FromUserID, "amount", req.Amount)
		return nil, errors.New("insufficient funds for transfer")
	}

	newTransaction := &domain.Transaction{
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        domain.Transfer,
		Amount:      req.Amount,
		Status:      domain.Pending,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction); err != nil {
		utils.Logger.Error("Failed to create transfer transaction record", "error", err)
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount.Neg()); err != nil {
		utils.Logger.Error("Failed to debit from sender for transfer", "error", err)
		return nil, fmt.Errorf("failed to debit from sender: %w", err)
	}

	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount); err != nil {
		utils.Logger.Error("Failed to credit to receiver for transfer", "error", err)
		return nil, fmt.Errorf("failed to credit to receiver: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		utils.Logger.Error("Failed to update transfer transaction status", "error", err)
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Failed to commit transfer transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	utils.Logger.Info("Transfer transaction completed", "transaction_id", newTransaction.ID, "from_user_id", req.FromUserID, "to_user_id", req.ToUserID, "amount", req.Amount)

	newTransaction.Status = domain.Completed

	resp := &PostTransactionTransferResponse{
		Transaction: newTransaction,
		Message:     "Transfer completed successfully.",
	}

	return resp, nil
}

func (s *transactionService) GetHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	if ctx.Value(middleware.UserIDKey) != *req.UserID && ctx.Value(middleware.UserRoleKey) != domain.AdminRole {
		return nil, errors.New("you dont have any access")
	}
	var opts []domain.TransactionQueryOption

	if req.Limit != nil {
		opts = append(opts, domain.Limit(*req.Limit))
	}
	if req.Offset != nil {
		opts = append(opts, domain.Offset(*req.Offset))
	}
	if req.Type != nil {
		opts = append(opts, domain.StatusType(*req.Type))
	}

	transactions, err := s.transactionRepo.GetTransactionsByUserID(ctx, *req.UserID, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by user id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	resp := &GetTransactionHistoryResponse{
		Transactions: transactions,
	}

	return resp, nil
}

func (s *transactionService) GetByID(ctx context.Context, req GetTransactionByIdRequest) (*GetTransactionByIdResponse, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin get by id transaction: %w", err)
	}
	defer tx.Rollback()

	transaction, err := s.transactionRepo.GetTranscaptionByID(ctx, tx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction with ID %s not found", req.ID)
		}
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	_ = tx.Commit()

	resp := &GetTransactionByIdResponse{
		Transaction: transaction,
	}

	return resp, nil
}
