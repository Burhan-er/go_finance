package main

import (
	"context"
	"database/sql"
	"go_finance/internal/api"
	"go_finance/internal/api/handler"
	"go_finance/internal/api/middleware"
	"go_finance/internal/config"
	"go_finance/internal/repository/postgres"
	"go_finance/internal/service"
	"go_finance/pkg/database"
	"go_finance/pkg/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		utils.Logger.Error("Could not load config", "error", err)
		os.Exit(1)
	}

	utils.InitLogger()
	utils.Logger.Info("Logger initialized")
	
	//DB
	maxRetries := 3
	var (
		db  *sql.DB
	)
	for i := 0; i < maxRetries; i++ {

		db, _, err = database.ConnectAndMigrateDB(cfg.DatabaseURL, "migrations")
		if err == nil {
			utils.Logger.Error("Migration failed", "error", err)
			break	
		}
		utils.Logger.Error("Migration attempt failed", "attempt", i+1, "error", err)
		if i<maxRetries-1{
			time.Sleep(2*time.Second)
		} 	
	}
	if err != nil {
		utils.Logger.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	utils.Logger.Info("Database connection established")

	//Repo And Services
	numWorkers,_ := strconv.Atoi(cfg.NumWorkers)
	jobQueueSize,_ := strconv.Atoi(cfg.JobQueueSize)
	auditLogRepo := postgres.NewAuditLogRepository(db)
	auditLogService := service.NewAuditLogService(auditLogRepo)

	userRepo := postgres.NewUserRepository(db)
	balanceRepo := postgres.NewBalanceRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	transactionProcessor := service.NewTransactionProcessor(
		numWorkers,
		jobQueueSize,
		transactionRepo,
		balanceRepo,
		db,
		auditLogService,
	)

	userService := service.NewUserService(userRepo, balanceRepo, cfg.JWTSecret, cfg.JWTExpiresIn, auditLogService)
	balanceService := service.NewBalanceService(balanceRepo, auditLogService)
	transactionService := service.NewTransactionService(transactionRepo, balanceRepo, db, auditLogService,transactionProcessor)

	transactionProcessor.Start(numWorkers)

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
		utils.Logger.Info("Server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Error("Could not listen", "addr", server.Addr, "error", err)
			os.Exit(1)
		}
	}()
		//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	transactionProcessor.Stop()
	utils.Logger.Info("Shutting down server...")


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	utils.Logger.Info("Server exiting")
}
