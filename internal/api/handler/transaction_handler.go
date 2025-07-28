package handler

import (
	"encoding/json"
	"fmt"
	"go_finance/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(ts service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: ts}
}

// POST /api/v1/transactions/credit
func (h *TransactionHandler) CreditTransaction(w http.ResponseWriter, r *http.Request) {

	var req service.PostTransactionCreditRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Servis katmanını çağır
	transaction, err := h.transactionService.Credit(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// POST /api/v1/transactions/debit
func (h *TransactionHandler) DebitTransaction(w http.ResponseWriter, r *http.Request) {
	var req service.PostTransactionDebitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.Debit(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// POST /api/v1/transactions/transfer
func (h *TransactionHandler) TransferTransaction(w http.ResponseWriter, r *http.Request) {
	var req service.PostTransactionTransferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.transactionService.Transfer(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// GET /api/v1/transactions/history
func (h *TransactionHandler) TransactionHistory(w http.ResponseWriter, r *http.Request) {

	var req service.GetTransactionHistoryRequest

	pageStr := r.URL.Query().Get("Page")
	limitStr := r.URL.Query().Get("Limit")

	page, pageConvertErr := strconv.Atoi(pageStr)
	if pageConvertErr != nil {
		http.Error(w, fmt.Sprintf("%s is not a Number", pageStr), http.StatusBadRequest)
		return
	}

	limit, limitConvertError := strconv.Atoi(limitStr)
	if limitConvertError != nil {
		http.Error(w, fmt.Sprintf("%s is not a Number", limitStr), http.StatusBadRequest)
		return
	}

	req.Page = page
	req.Limit = limit

	history, err := h.transactionService.GetHistory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

// GET /api/v1/transactions/{id}
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	var req service.GetTransactionByIdRequest

	transactionID := chi.URLParam(r, "id")
	if transactionID == "" {
		http.Error(w, "Transaction ID is required", http.StatusBadRequest)
		return
	}
	req.ID = transactionID

	transaction, err := h.transactionService.GetByID(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}
