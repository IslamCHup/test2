package database

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/islamchupanov/tz1/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(cfg config.DBConfig, logger *slog.Logger) (*gorm.DB, error) {
	// Use SQLite for simplicity - stores data in devices.db file
	dbPath := filepath.Join(".", "devices.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		if logger != nil {
			logger.Error("failed to connect to sqlite", "error", err)
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

	logger.Info("connected to sqlite", "path", dbPath)

	return db, nil
}
