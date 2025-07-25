package api

import (
	"go_finance/internal/api/handler"
	"net/http"

	mWare "go_finance/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, authMiddleware *mWare.AuthMiddleware, userHandler *handler.UserHandler) http.Handler {
	r := chi.NewRouter()

	// Temel middleware'ler
	r.Use(middleware.Logger)    // Gelen istekleri loglar
	r.Use(middleware.Recoverer) // Panic durumlarında sunucunun çökmesini engeller
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register) // PUT /api/v1/auth/register
			r.Post("/login", authHandler.Login)       // PUT /api/v1/auth/login

			// r.Post("/refresh", authHandler.Refresh) // TODO: Implement refresh token logic
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/", userHandler.ListUsers)       // GET /api/v1/users
			r.Get("/{id}", userHandler.GetUserByID) // GET /api/v1/users/{id}
			//r.Put("/{id}", userHandler.UpdateUser)    // PUT /api/v1/users/{id}
			//r.Delete("/{id}", userHandler.DeleteUser) // DELETE /api/v1/users/{id}
		})

		// r.Route("/transactions", func(r chi.Router) {
		// 	r.Post("/credit", handler.CreditTransaction)     // POST /api/v1/transactions/credit
		// 	r.Post("/debit", handler.DebitTransaction)       // POST /api/v1/transactions/debit
		// 	r.Post("/transfer", handler.TransferTransaction) // POST /api/v1/transactions/transfer
		// 	r.Get("/history", handler.TransactionHistory)    // GET /api/v1/transactions/history
		// 	r.Get("/{id}", handler.GetTransactionByID)       // GET /api/v1/transactions/{id}
		// })

		// r.Route("/balances", func(r chi.Router) {
		// 	r.Get("/current", handler.GetCurrentBalance)        // GET /api/v1/balances/current
		// 	r.Get("/historical", handler.GetHistoricalBalances) // GET /api/v1/balances/historical
		// 	r.Get("/at-time", handler.GetBalanceAtTime)         // GET /api/v1/balances/at-time
		// })
	})

	return r
}
