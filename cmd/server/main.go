package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpdelivery "github.com/company/payment-service/internal/delivery/http"
	"github.com/company/payment-service/internal/domain/usecase"
	"github.com/company/payment-service/internal/infrastructure/config"
	"github.com/company/payment-service/internal/infrastructure/database"
	"github.com/company/payment-service/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 載入 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 載入配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日誌
	logger, err := logger.NewLogger(logger.Config{
		Level:      cfg.Logger.Level,
		Format:     cfg.Logger.Format,
		OutputPath: cfg.Logger.OutputPath,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 初始化數據庫連接
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// 初始化 repositories
	paymentRepo := database.NewPaymentRepository(db)
	merchantRepo := database.NewMerchantRepository(db)
	// customerRepo := database.NewCustomerRepository(db) // 需要實現

	// 初始化 use cases
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, merchantRepo, nil) // 暫時傳入 nil

	// 設置路由
	router := httpdelivery.SetupRouter(paymentUseCase, merchantRepo)

	// 創建 HTTP 服務器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 啟動服務器
	go func() {
		logger.Info(fmt.Sprintf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}