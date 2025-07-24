package main

import (
	"context"
	"go_finance/internal/api"
	"go_finance/internal/api/handler"
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
	// 1. Konfigürasyonu yükle
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Logger'ı ayarla (Şimdilik standart log, pkg/logger ile geliştirilebilir)
	log.Println("Logger initialized")

	// 3. Veritabanı bağlantısı kur
	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established")

	// ---- Bağımlılıkları Oluşturma (Dependency Injection) ----

	// 4. Repository katmanını oluştur
	userRepo := postgres.NewUserRepository(db)

	// 5. Servis katmanını oluştur
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiresIn)

	// 6. Handler (API) katmanını oluştur
	authHandler := handler.NewAuthHandler(userService)

	// 7. Router'ı oluştur ve handler'ları kaydet
	router := api.NewRouter(authHandler)

	// 8. HTTP Sunucusunu ayarla
	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	// 9. Graceful Shutdown (Düzgün Kapatma) mekanizmasını kur
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
