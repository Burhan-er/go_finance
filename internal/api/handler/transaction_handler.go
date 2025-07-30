package handler

import (
	"encoding/json"
	"fmt"
	"go_finance/internal/domain"
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

func (h *TransactionHandler) CreditTransaction(w http.ResponseWriter, r *http.Request) {

	var req service.PostTransactionCreditRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.Credit(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

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

type StrConv struct {
	IntValue    int
	StringValue string
}

func (h *TransactionHandler) TransactionHistory(w http.ResponseWriter, r *http.Request) {

	var req service.GetTransactionHistoryRequest

	userId := r.URL.Query().Get("user_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	ttype := r.URL.Query().Get("type")

	var offset int
	var offsetConvertErr error
	if offsetStr != "" {
		offset, offsetConvertErr = strconv.Atoi(offsetStr)
		if offsetConvertErr != nil {
			http.Error(w, fmt.Sprintf("%s is not a Number", offsetStr), http.StatusBadRequest)
			return
		}
	}
	var limit int
	var limitConvertErr error
	if limitStr != "" {
		limit, limitConvertErr = strconv.Atoi(limitStr)
		if limitConvertErr != nil {
			http.Error(w, fmt.Sprintf("%s is not a Number", limitStr), http.StatusBadRequest)
			return
		}
	}

	req.UserID = &userId
	req.Offset = &offset
	ttypeVal := domain.TransactionType(ttype)
	req.Type = &ttypeVal
	req.Limit = &limit

	history, err := h.transactionService.GetHistory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}

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
