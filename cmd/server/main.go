package main

import (
	"context"
	"go_finance/internal/api"
	"go_finance/internal/api/handler"
	"go_finance/internal/api/middleware"
	"go_finance/internal/config"
	"go_finance/internal/repository/postgres"
	"go_finance/internal/service"
	"go_finance/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	log.Println("Logger initialized")

	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	//Repo And Services
	userRepo := postgres.NewUserRepository(db)
	balanceRepo := postgres.NewBalanceRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	balanceService := service.NewBalanceService(balanceRepo)
	userService := service.NewUserService(userRepo, balanceRepo,cfg.JWTSecret, cfg.JWTExpiresIn)
	transactionService := service.NewTransactionService(transactionRepo, balanceRepo, db)

	//Handler
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(userService)
	balanceHandler := handler.NewBalanceHandler(balanceService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	//MiddleWare
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	HandlerOfAll := api.Handlers{
		Auth:        authHandler,
		User:        userHandler,
		Transaction: transactionHandler,
		Balance:     balanceHandler,
	}

	router := api.NewRouter(authMiddleware, HandlerOfAll)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
