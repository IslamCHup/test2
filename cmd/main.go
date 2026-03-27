package main

import (
	"log"
	"os"

	"github.com/islamchupanov/tz1/internal/config"
	database "github.com/islamchupanov/tz1/internal/db"
	"github.com/islamchupanov/tz1/internal/handler"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/repository"
	"github.com/islamchupanov/tz1/internal/router"
	"github.com/islamchupanov/tz1/internal/service"
)

func main() {
	cfg := config.Load()

	logger := logger.InitLog(cfg.LogLevel)

	dbConn, err := database.InitDB(cfg.DB, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	if err := dbConn.AutoMigrate(&model.Device{}); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	logger.Info("database migration completed")

	deviceRepo := repository.NewDeviceRepository(dbConn, logger)
	deviceService := service.NewDeviceService(deviceRepo, logger)
	deviceHandler := handler.NewDeviceHandler(deviceService, logger)

	r := router.SetupRouter(deviceHandler)

	logger.Info("starting server", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
