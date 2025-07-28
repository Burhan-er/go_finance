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

func NewRouter(authMiddleware *mWare.AuthMiddleware,h Handlers) http.Handler {
	r := chi.NewRouter()

	// Temel middleware'ler
	r.Use(middleware.Logger)    // Gelen istekleri loglar
	r.Use(middleware.Recoverer) // Panic durumlarında sunucunun çökmesini engeller
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register) // PUT /api/v1/auth/register
			r.Post("/login", h.Auth.Login)       // PUT /api/v1/auth/login

			// r.Post("/refresh", authHandler.Refresh) // TODO: Implement refresh token logic
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/", h.User.ListUsers)         // GET /api/v1/users
			r.Get("/{id}", h.User.GetUserByID)   // GET /api/v1/users/{id}
			r.Put("/{id}", h.User.UpdateUser)    // PUT /api/v1/users/{id}
			r.Delete("/{id}", h.User.DeleteUser) // DELETE /api/v1/users/{id}
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Post("/credit", h.Transaction.CreditTransaction)     // POST /api/v1/transactions/credit
			r.Post("/debit", h.Transaction.DebitTransaction)       // POST /api/v1/transactions/debit
			r.Post("/transfer", h.Transaction.TransferTransaction) // POST /api/v1/transactions/transfer
			r.Get("/history", h.Transaction.TransactionHistory)    // GET /api/v1/transactions/history
			r.Get("/{id}", h.Transaction.GetTransactionByID)       // GET /api/v1/transactions/{id}
		})

		r.Route("/balances",func(r chi.Router) {
			r.Get("/current", h.Balance.GetCurrentBalance)        // GET /api/v1/balances/current
			r.Get("/historical", h.Balance.GetHistoricalBalances) // GET /api/v1/balances/historical
			r.Get("/at-time", h.Balance.GetBalanceAtTime)         // GET /api/v1/balances/at-time
		})
	})

	return r
}
