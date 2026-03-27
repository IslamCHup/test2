package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/islamchupanov/tz1/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.DBConfig, logger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if logger != nil {
			logger.Error("failed to connect to postgres", "error", err)
		}
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		if logger != nil {
			logger.Error("failed to get sql.DB from gorm.DB", "error", err)
		}
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		if logger != nil {
			logger.Error("failed to ping db", "error", err)
		}
		return nil, err
	}

	logger.Info("connected to postgres",
		"host", cfg.Host,
		"port", cfg.Port,
		"db", cfg.Name,
	)

	return db, nil
}
