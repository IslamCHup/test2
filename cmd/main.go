package main

import (
	"context"
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
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/repository"
	"github.com/islamchupanov/tz1/internal/router"
	"github.com/islamchupanov/tz1/internal/service"
)

// @title Device API
// @version 1.0
// @description API for managing devices
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	appLogger := appLogger.InitLog(cfg.LogLevel)

	dbConn, err := database.InitDB(cfg.DB, appLogger)
	if err != nil {
		appLogger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Проверка подключения к БД через ping
	if err := dbConn.DB().Ping(); err != nil {
		appLogger.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	appLogger.Info("database connection verified")

	if err := dbConn.AutoMigrate(&model.Device{}); err != nil {
		appLogger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	appLogger.Info("database migration completed")

	deviceRepo := repository.NewDeviceRepository(dbConn, appLogger)
	deviceService := service.NewDeviceService(deviceRepo, appLogger)
	deviceHandler := handler.NewDeviceHandler(deviceService, appLogger)

	r := router.SetupRouter(deviceHandler)

	// Валидация порта
	port := cfg.Port
	if port == "" {
		port = "8080"
		appLogger.Warn("APP_PORT is empty, using default port 8080")
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
		appLogger.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("server forced to shutdown", "error", err)
	}

	// Закрытие соединения с БД
	sqlDB, err := dbConn.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("failed to close database connection", "error", err)
		}
	}

	appLogger.Info("server exited gracefully")
}
