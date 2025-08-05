package api

import (
	"go_finance/internal/api/handler"
	"net/http"

	mWare "go_finance/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	Auth        *handler.AuthHandler
	User        *handler.UserHandler
	Transaction *handler.TransactionHandler
	Balance     *handler.BalanceHandler
}

func NewRouter(authMiddleware *mWare.AuthMiddleware, h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/", h.User.ListUsers)
			r.Get("/{id}", h.User.GetUserByID)
			r.Put("/{id}", h.User.UpdateUser)
			r.Delete("/{id}", h.User.DeleteUser)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/credit", h.Transaction.CreditTransaction)
			r.Post("/debit", h.Transaction.DebitTransaction)
			r.Post("/transfer", h.Transaction.TransferTransaction)
			r.Get("/history", h.Transaction.TransactionHistory)
			r.Get("/{id}", h.Transaction.GetTransactionByID)
		})

		r.Route("/balances", func(r chi.Router) {
			r.Get("/current", h.Balance.GetCurrentBalance)
			r.Get("/historical", h.Balance.GetHistoricalBalances)
			r.Get("/at-time", h.Balance.GetBalanceAtTime)
		})
	})

	return r
}
