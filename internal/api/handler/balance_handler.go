package handler

import (
	"encoding/json"
	"fmt"
	"go_finance/internal/service" // Proje yolunuzu buna göre güncelleyin
	"net/http"
	"strconv"
	"time"
)

// BalanceHandler, bakiye ile ilgili HTTP isteklerini yönetir.
type BalanceHandler struct {
	balanceService service.BalanceService
}

// NewBalanceHandler, yeni bir BalanceHandler örneği oluşturur.
func NewBalanceHandler(bs service.BalanceService) *BalanceHandler {
	return &BalanceHandler{balanceService: bs}
}

// GetCurrentBalance, kullanıcının mevcut bakiyesini getirir.
// GET /api/v1/balances/current
func (h *BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	// Genellikle JWT token gibi bir yerden kullanıcı ID'si alınır ve
	// context ile servis katmanına iletilir.
	// Şimdilik servis katmanının bunu hallettiğini varsayıyoruz.
	
	var req service.GetBalanceCurrentRequest


	balance, err := h.balanceService.GetCurrent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

// GetHistoricalBalances, kullanıcının geçmiş bakiye kayıtlarını (anlık görüntülerini) getirir.
// GET /api/v1/balances/historical
func (h *BalanceHandler) GetHistoricalBalances(w http.ResponseWriter, r *http.Request) {
	// Pagination parametreleri (page, limit) burada da kullanılabilir.
	// r.URL.Query() ile alabilirsiniz.
	var req service.GetBalanceHistoricalRequest

	pageStr := r.URL.Query().Get("Page")
	limitStr := r.URL.Query().Get("Limit")

	page,pageConvertErr := strconv.Atoi(pageStr);
	if  pageConvertErr !=nil{
		http.Error(w,fmt.Sprintf("%s is not a Number",pageStr),http.StatusBadRequest)
		return
	}
	
	limit,limitConvertError:= strconv.Atoi(limitStr);
	if  limitConvertError !=nil{
		http.Error(w,fmt.Sprintf("%s is not a Number",limitStr),http.StatusBadRequest)
		return
	}

	req.Page = page
	req.Limit = limit

	balances, getBalanceErr := h.balanceService.GetHistorical(r.Context(),req)
	if getBalanceErr != nil {
		http.Error(w, getBalanceErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balances)
}

// GetBalanceAtTime, belirli bir zamandaki bakiyeyi getirir.
// GET /api/v1/balances/at-time?timestamp=2023-10-27T10:00:00Z
func (h *BalanceHandler) GetBalanceAtTime(w http.ResponseWriter, r *http.Request) {
	var req service.GetBalanceAtTimeRequest
	
	timestampStr := r.URL.Query().Get("timestamp")
	if timestampStr == "" {
		http.Error(w, "Query parameter 'timestamp' is required", http.StatusBadRequest)
		return
	}

	// Gelen string'i time.Time objesine çevir. RFC3339 formatı (ISO 8601) standarttır.
	t, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		http.Error(w, "Invalid timestamp format. Use RFC3339 (e.g., 2023-10-27T10:00:00Z)", http.StatusBadRequest)
		return
	}
	req.Timestamp = t

	balance, err := h.balanceService.GetAtTime(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}