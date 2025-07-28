package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_finance/internal/domain"
	"go_finance/internal/repository" // Repository'lerinizin bu pakette olduğunu varsayıyorum
	"time"

	"github.com/google/uuid"
)

// transactionService, TransactionService arayüzünü uygular.
type transactionService struct {
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
	db              *sql.DB // Veritabanı transactionlarını yönetmek için
}

// NewTransactionService, transactionService'in yeni bir örneğini oluşturur.
// Bu, bağımlılık enjeksiyonu için kullanılır.
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

// Credit, kullanıcının hesabına para yatırma işlemini yönetir.
func (s *transactionService) Credit(ctx context.Context, req PostTransactionCreditRequest) (*PostTransactionCreditResponse, error) {
	// Atomikliği sağlamak için bir veritabanı transaction'ı başlat
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Fonksiyon sonunda bir hata olursa transaction'ı geri al (rollback)
	defer tx.Rollback()

	// Yeni bir transaction domain nesnesi oluştur
	newTransaction := &domain.Transaction{
		ID:          uuid.New().String(),
		UserID:      req.FromUserID,
		ToUserID:    req.ToUserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Status:      domain.Pending,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Transaction'ı 'Pending' durumuyla veritabanına kaydet
	if err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Kullanıcının bakiyesini güncelle
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// Transaction durumunu 'Completed' olarak güncelle
	// NOT: Verdiğiniz repo arayüzünde UpdateTransactionStatus *sql.Tx almıyor.
	// Ancak atomiklik için alması gerekir. Kodun bu şekilde olması gerektiğini varsayarak yazıyorum.
	// Gerekirse repository'nizi `UpdateTransactionStatus(ctx context.Context, tx *sql.Tx, id string, status domain.StatusType) error` şeklinde güncelleyin.
	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Her şey yolunda gittiyse, transaction'ı onayla (commit)
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newTransaction.Status = domain.Completed // Dönen nesnenin durumunu da güncelle

	resp := &PostTransactionCreditResponse{
		Transaction: newTransaction,
		Message:     "Credit transaction completed successfully.",
	}

	return resp, nil
}

// Debit, kullanıcının hesabından para çekme işlemini yönetir.
func (s *transactionService) Debit(ctx context.Context, req PostTransactionDebitRequest) (*PostTransactionDebitResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Para çekmeden önce bakiye kontrolü yap
	currentBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	if currentBalance.Amount < req.Amount { // last I'm here
		return nil, errors.New("insufficient funds")
	}

	newTransaction := &domain.Transaction{
		ID:          uuid.New().String(),
		UserID:      req.FromUserID,
		Type:        req.Type,
		Amount:      req.Amount,
		Status:      domain.Pending,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.CreateTransaction(ctx, tx, newTransaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Bakiyeyi azaltmak için negatif miktar gönder
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, -req.Amount); err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newTransaction.Status = domain.Completed

	resp := &PostTransactionDebitResponse{
		Transaction: newTransaction,
		Message:     "Debit transaction completed successfully.",
	}

	return resp, nil
}

// Transfer, iki kullanıcı arasında para transferi işlemini yönetir.
func (s *transactionService) Transfer(ctx context.Context, req PostTransactionTransferRequest) (*PostTransactionTransferResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Gönderen kullanıcının bakiye kontrolü
	senderBalance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender's balance: %w", err)
	}

	if senderBalance.Amount < req.Amount {
		return nil, errors.New("insufficient funds for transfer")
	}

	// Not: domain.Transaction'da transferin karşı tarafını tutmak için
	// RelatedUserID gibi bir alan olması faydalı olur.
	newTransaction := &domain.Transaction{
		ID:          uuid.New().String(),
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
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Gönderenden parayı düş
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.FromUserID, -req.Amount); err != nil {
		return nil, fmt.Errorf("failed to debit from sender: %w", err)
	}

	// Alıcıya parayı ekle
	if err := s.balanceRepo.UpdateBalance(ctx, tx, req.ToUserID, req.Amount); err != nil {
		return nil, fmt.Errorf("failed to credit to receiver: %w", err)
	}

	if err := s.transactionRepo.UpdateTransactionStatus(ctx, tx, newTransaction.ID, domain.Completed); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newTransaction.Status = domain.Completed

	resp := &PostTransactionTransferResponse{
		Transaction: newTransaction,
		Message:     "Transfer completed successfully.",
	}

	return resp, nil
}

// GetHistory, bir kullanıcının işlem geçmişini getirir.
// NOT: Verdiğiniz GetTransactionHistoryRequest'te UserID yok. Genellikle işlem geçmişi
// belirli bir kullanıcı için istenir. Bu yüzden `GetTransactionsByUserID` kullanıldı.
// Gerekirse isteğe UserID eklenmeli veya tüm işlemleri getiren bir repo metodu yazılmalıdır.
// Ayrıca repoda sayfalama (pagination) desteği eklenmesi (offset, limit) daha verimli olacaktır.
func (s *transactionService) GetHistory(ctx context.Context, req GetTransactionHistoryRequest) (*GetTransactionHistoryResponse, error) {
	// Bu örnekte, UserID'nin bir şekilde (örneğin JWT'den) alınıp
	// GetTransactionsByUserID'ye verildiği varsayılmıştır.
	// Şimdilik bu metot, verdiğiniz arayüzlere göre tam olarak implemente edilemez.
	// Örnek olarak boş bir liste döndürüyorum ve bir hata mesajı veriyorum.
	// Gerçek bir senaryoda bu kısmın tasarıma göre doldurulması gerekir.

	// Örnek: userID := "jwt-den-gelen-id"
	// transactions, err := s.transactionRepo.GetTransactionsByUserID(ctx, userID)
	return nil, errors.New("GetHistory requires a UserID, which is missing from the request struct; also repository needs pagination support")

	// Eğer tüm işlemleri getiren bir repo metodunuz olsaydı (ör: GetAllTransactions(ctx, limit, offset)) şöyle olurdu:
	// transactions, err := s.transactionRepo.GetAllTransactions(ctx, req.Limit, (req.Page-1)*req.Limit)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get transaction history: %w", err)
	// }
	//
	// resp := &GetTransactionHistoryResponse{
	// 	Transactions: transactions,
	// }
	// return resp, nil
}

// GetByID, ID'ye göre tek bir işlemi getirir.
func (s *transactionService) GetByID(ctx context.Context, req GetTransactionByIdRequest) (*GetTransactionByIdResponse, error) {
	transaction, err := s.transactionRepo.GetTransactionByID(ctx, req.ID)
	if err != nil {
		// Veritabanında bulunamama durumunu (sql.ErrNoRows) daha spesifik bir hata ile dönebiliriz.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction with ID %s not found", req.ID)
		}
		return nil, fmt.Errorf("failed to get transaction by ID: %w", err)
	}

	resp := &GetTransactionByIdResponse{
		Transaction: transaction,
	}

	return resp, nil
}
