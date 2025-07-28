package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// Modern ve yapısallaştırılmış loglama için
	// Projenizin domain ve repository paketlerinin yollarını buraya yazın
	"go_finance/internal/domain"
	"go_finance/internal/repository"
)

type balanceService struct {
	balanceRepo repository.BalanceRepository
	//logger      *slog.Logger
}

// NewBalanceService, BalanceService'in yeni bir örneğini oluşturur.
// Bu, "Dependency Injection" için kullanılan bir yapıcı fonksiyondur.
func NewBalanceService(repo repository.BalanceRepository) BalanceService {
	return &balanceService{
		balanceRepo: repo,
		// logger:      logger,
	}
}

// GetCurrent, bir kullanıcının mevcut bakiyesini alır.
func (s *balanceService) GetCurrent(ctx context.Context, req GetBalanceCurrentRequest) (*GetBalanceCurrentResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	balance, err := s.balanceRepo.GetBalanceByUserID(ctx, req.UserID)
	if err != nil {
		// s.logger.Error("failed to get balance from repository", "user_id", req.UserID, "error", err)
		// Burada hatayı sarmalayarak daha anlamlı bir mesaj verebiliriz.
		// Örneğin sql.ErrNoRows hatasını "kullanıcı bulunamadı" gibi bir hataya çevirebiliriz.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with id: %s", req.UserID)
		}
		return nil, fmt.Errorf("internal server error") // Kullanıcıya iç hatayı göstermeyiz.
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

// GetHistorical, bir kullanıcının belirli bir zaman aralığındaki bakiye geçmişini alır.
// NOT: Sağlanan BalanceRepository bu işlevselliği desteklememektedir.
// Bu, gerçek bir senaryoda nasıl görünebileceğini gösteren bir yer tutucudur.
func (s *balanceService) GetHistorical(ctx context.Context, req GetBalanceHistoricalRequest) (*GetBalanceHistoricalResponse, error) {
	// s.logger.Warn("GetHistorical is not fully implemented due to repository limitations", "user_id", req.UserID)

	// GERÇEK UYGULAMA NASIL OLURDU?
	// 1. Repository'de `GetTransactionHistory(ctx, userID, startDate, endDate)` gibi bir metodunuz olurdu.
	// 2. Bu metot, `transactions` tablosundan ilgili kayıtları çekerdi.
	//    transactions, err := s.transactionRepo.GetHistory(ctx, req.UserID, req.StartDate, req.EndDate)
	// 3. Servis, bu veriyi işleyip `GetBalanceHistoricalResponse`'a dönüştürürdü.

	// Placeholder (Yer Tutucu) Cevap:
	// Bu kodun derlenmesini ve arayüzü tatmin etmesini sağlar.
	return &GetBalanceHistoricalResponse{
		//UserID: req.Limit,
		/*	Transactions: []BalanceTransaction{
				// Örnek veri
				{Timestamp: time.Now().Add(-24 * time.Hour), Amount: 100, Description: "Maaş Yattı"},
				{Timestamp: time.Now().Add(-12 * time.Hour), Amount: -20, Description: "Kahve"},
			},
		}*/}, nil // Şimdilik hata döndürmüyoruz.

}

// GetAtTime, bir kullanıcının belirli bir zamandaki bakiyesini hesaplar.
// NOT: Sağlanan BalanceRepository bu işlevselliği desteklememektedir.
// Bu, gerçek bir senaryoda nasıl görünebileceğini gösteren bir yer tutucudur.
func (s *balanceService) GetAtTime(ctx context.Context, req GetBalanceAtTimeRequest) (*GetBalanceAtTimeResponse, error) {
	//s.logger.Warn("GetAtTime is not fully implemented due to repository limitations", "user_id", req.UserID)

	// GERÇEK UYGULAMA NASIL OLURDU?
	// Bu, genellikle "Event Sourcing" benzeri bir yaklaşımla çözülür.
	// 1. Belirtilen zamana (`req.Timestamp`) kadar olan TÜM işlemleri (`transactions`) veritabanından çekerdiniz.
	//    `transactions, err := s.transactionRepo.GetTransactionsUntil(ctx, req.UserID, req.Timestamp)`
	// 2. Başlangıç bakiyesinden (genellikle 0) başlayarak bu işlemlerin toplamını hesaplardınız.
	//    var calculatedAmount int = 0
	//    for _, tx := range transactions {
	//        calculatedAmount += tx.Amount
	//    }
	// 3. Hesaplanan bu tutar, o zamanki bakiyeyi verirdi.

	// Placeholder (Yer Tutucu) Cevap:
	return &GetBalanceAtTimeResponse{
		// UserID:    req.UserID,
		// Amount:    500, // Örnek bir değer
		// Timestamp: req.Timestamp,
	}, nil // Şimdilik hata döndürmüyoruz.
}
