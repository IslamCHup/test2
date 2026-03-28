package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/islamchupanov/tz1/internal/config"
	database "github.com/islamchupanov/tz1/internal/db"
	"github.com/islamchupanov/tz1/internal/handler"
	appLogger "github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/repository"
	"github.com/islamchupanov/tz1/internal/router"
	"github.com/islamchupanov/tz1/internal/service"

	_ "github.com/islamchupanov/tz1/docs"
)

// @title Device API
// @version 1.0
// @description API for managing devices
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	// Валидация конфигурации
	if err := cfg.Validate(); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

	logger := appLogger.InitLog(cfg.LogLevel)

	// Инициализация БД с retry
	dbConn, sqlDB, err := database.InitDBWithRetry(cfg.DB, logger, 5, 2*time.Second)
	if err != nil {
		logger.Error("failed to initialize database after retries", "error", err)
		os.Exit(1)
	}

	logger.Info("database connection verified")

	// Используем SQL миграции вместо AutoMigrate для контроля схемы
	// AutoMigrate удален - используем migrations/*.sql через goose или другой инструмент
	// Для простоты оставляем только базовую проверку подключения
	logger.Info("database connection established, migrations should be applied separately via goose")

	deviceRepo := repository.NewDeviceRepository(dbConn, logger)
	deviceService := service.NewDeviceService(deviceRepo, logger)
	deviceHandler := handler.NewDeviceHandler(deviceService, logger)

	r, err := router.SetupRouter(deviceHandler)
	if err != nil {
		logger.Error("failed to setup router", "error", err)
		os.Exit(1)
	}

	// Валидация порта
	port := cfg.Port
	if port == "" {
		port = "8080"
		logger.Warn("APP_PORT is empty, using default port 8080")
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	// Закрытие соединения с БД (используем уже полученный sqlDB)
	if db, ok := sqlDB.(*sql.DB); ok {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database connection", "error", err)
		}
	}

	logger.Info("server exited gracefully")
}
