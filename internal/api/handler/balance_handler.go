package handler

import (
	"encoding/json"
	"go_finance/internal/service"
	"net/http"
	"time"
)

type BalanceHandler struct {
	balanceService service.BalanceService
}

func NewBalanceHandler(bs service.BalanceService) *BalanceHandler {
	return &BalanceHandler{balanceService: bs}
}

func (h *BalanceHandler) GetCurrentBalance(w http.ResponseWriter, r *http.Request) {
	var req service.GetBalanceCurrentRequest
	req.UserID = r.URL.Query().Get("user_id")
	balance, err := h.balanceService.GetCurrent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func (h *BalanceHandler) GetHistoricalBalances(w http.ResponseWriter, r *http.Request) {
	var req service.GetBalanceHistoricalRequest

	// pageStr := r.URL.Query().Get("Page")
	// limitStr := r.URL.Query().Get("Limit")

	// page, pageConvertErr := strconv.Atoi(pageStr)
	// if pageConvertErr != nil {
	// 	http.Error(w, fmt.Sprintf("%s is not a Number", pageStr), http.StatusBadRequest)
	// 	return
	// }

	// limit, limitConvertError := strconv.Atoi(limitStr)
	// if limitConvertError != nil {
	// 	http.Error(w, fmt.Sprintf("%s is not a Number", limitStr), http.StatusBadRequest)
	// 	return
	// }

	// req.Page = page
	// req.Limit = limit

	balances, getBalanceErr := h.balanceService.GetHistorical(r.Context(), req)
	if getBalanceErr != nil {
		http.Error(w, getBalanceErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balances)
}

func (h *BalanceHandler) GetBalanceAtTime(w http.ResponseWriter, r *http.Request) {
	var req service.GetBalanceAtTimeRequest

	timestampStr := r.URL.Query().Get("timestamp")
	if timestampStr == "" {
		http.Error(w, "Query parameter 'timestamp' is required", http.StatusBadRequest)
		return
	}

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
