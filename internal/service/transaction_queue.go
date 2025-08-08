package service

import (
	"context"
	"database/sql"
	"fmt"
	"go_finance/internal/domain"
	"go_finance/internal/repository"
	"go_finance/pkg/utils"
	"sync"

	"github.com/shopspring/decimal"
)

type Job struct {
	TransactionID string
	Type          domain.TransactionType
	FromUserID    string
	ToUserID      string
	Amount        decimal.Decimal
}

type TransactionProcessor struct {
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
	db              *sql.DB
	auditLogService AuditLogService

	jobs chan Job
	wg   sync.WaitGroup
	quit chan struct{}
}

func NewTransactionProcessor(
	numWorkers int,
	jobQueueSize int,
	tr repository.TransactionRepository,
	br repository.BalanceRepository,
	db *sql.DB,
	auditLogService AuditLogService,
) *TransactionProcessor {
	return &TransactionProcessor{
		transactionRepo: tr,
		balanceRepo:     br,
		db:              db,
		auditLogService: auditLogService,
		jobs:            make(chan Job, jobQueueSize),
		quit:            make(chan struct{}),
	}
}

func (p *TransactionProcessor) Start(numWorkers int) {
	utils.Logger.Info("Starting transaction processor", "num_workers", numWorkers)
	for i := 1; i <= numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *TransactionProcessor) Stop() {
	utils.Logger.Info("Stopping transaction processor...")
	close(p.quit)
	p.wg.Wait()
	close(p.jobs)
	utils.Logger.Info("Transaction processor stopped.")
}

func (p *TransactionProcessor) SubmitJob(job Job) {
	p.jobs <- job
}

func (p *TransactionProcessor) worker(id int) {
	defer p.wg.Done()
	utils.Logger.Info("Worker started", "worker_id", id)
	for {
		select {
		case job := <-p.jobs:
			utils.Logger.Info("Worker picked up job", "worker_id", id, "transaction_id", job.TransactionID, "type", job.Type)
			p.processTransaction(context.Background(), job)
		case <-p.quit:
			utils.Logger.Info("Worker shutting down", "worker_id", id)
			return
		}
	}
}

func (p *TransactionProcessor) processTransaction(ctx context.Context, job Job) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		utils.Logger.Error("Worker failed to begin transaction", "transaction_id", job.TransactionID, "error", err)
		p.updateTransactionStatus(job.TransactionID, domain.Failed, "failed to begin db transaction")
		return
	}
	defer tx.Rollback()

	if job.Type == domain.TransactionType(domain.Credit) || job.Type == domain.TransactionType(domain.Transfer) {
		balance, err := p.balanceRepo.GetBalanceByUserID(ctx, job.FromUserID)
		if err != nil {
			utils.Logger.Error("Worker failed to get balance", "transaction_id", job.TransactionID, "user_id", job.FromUserID, "error", err)
			p.updateTransactionStatus(job.TransactionID, domain.Failed, "failed to retrieve sender balance")
			return
		}
		if balance.Amount.Cmp(job.Amount) < 0 {
			utils.Logger.Warn("Worker: insufficient funds", "transaction_id", job.TransactionID, "user_id", job.FromUserID)
			p.updateTransactionStatus(job.TransactionID, domain.Failed, "insufficient funds")
			return
		}
	}

	var senderAmount, receiverAmount decimal.Decimal

	switch job.Type {
	case domain.TransactionType(domain.Credit):
		senderAmount = decimal.Decimal{} //there are empty i expect be able to maybe bank account
		receiverAmount = job.Amount
	case domain.TransactionType(domain.Debit):
		senderAmount = job.Amount
		receiverAmount = decimal.Decimal{} //""""
	case domain.TransactionType(domain.Transfer):
		senderAmount = job.Amount.Neg()
		receiverAmount = job.Amount
		
	default:
		utils.Logger.Error("Worker: unknown transaction type", "transaction_id", job.TransactionID, "type", job.Type)
		p.updateTransactionStatus(job.TransactionID, domain.Failed, fmt.Sprintf("unknown transaction type: %s", job.Type))
		return
	}

	if job.Type != domain.TransactionType(domain.Credit) {
		if err := p.balanceRepo.UpdateBalance(ctx, tx, job.FromUserID, senderAmount); err != nil {
			utils.Logger.Error("Worker failed to update sender balance", "transaction_id", job.TransactionID, "error", err)
			p.updateTransactionStatus(job.TransactionID, domain.Failed, "failed to update sender balance")
			return
		}
	}
	if err := p.balanceRepo.UpdateBalance(ctx, tx, job.ToUserID, receiverAmount); err != nil {
		utils.Logger.Error("Worker failed to update receiver balance", "transaction_id", job.TransactionID, "error", err)
		p.updateTransactionStatus(job.TransactionID, domain.Failed, "failed to update receiver balance")
		return
	}

	if err := p.transactionRepo.UpdateTransactionStatus(ctx, tx, nil, job.TransactionID, domain.Completed); err != nil {
		utils.Logger.Error("Worker failed to update transaction status to completed", "transaction_id", job.TransactionID, "error", err)
		return
	}

	if err := tx.Commit(); err != nil {
		utils.Logger.Error("Worker failed to commit transaction", "transaction_id", job.TransactionID, "error", err)
		p.updateTransactionStatus(job.TransactionID, domain.Failed, "failed to commit transaction")
		return
	}

	p.auditLogService.Create(ctx, "transaction", job.TransactionID, fmt.Sprintf("%s_success", job.Type), fmt.Sprintf("Transaction %s completed successfully", job.Type))
	utils.Logger.Info("Worker completed job successfully", "transaction_id", job.TransactionID)
}

func (p *TransactionProcessor) updateTransactionStatus(transactionID string, status domain.StatusType, reason string) {
	err := p.transactionRepo.UpdateTransactionStatus(context.Background(), nil, p.db, transactionID, status)
	if err != nil {
		utils.Logger.Error("CRITICAL: Worker failed to update transaction status to FAILED", "transaction_id", transactionID, "error", err)
	}
	p.auditLogService.Create(context.Background(), "transaction", transactionID, "transaction_failed", reason)
}
