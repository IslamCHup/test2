package database

import (
	"fmt"
	"time"

	"github.com/islamchupanov/tz1/internal/config"
	"github.com/islamchupanov/tz1/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.DBConfig, logger *logger.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if logger != nil {
			logger.Error("failed to connect to postgres: %v", err)
		}
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		if logger != nil {
			logger.Error("failed to get sql.DB from gorm.DB: %v", err)
		}
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		if logger != nil {
			logger.Error("failed to ping db: %v", err)
		}
		return nil, err
	}

	logger.Info("connected to postgres: host=%s port=%s dbname=%s", cfg.Host, cfg.Port, cfg.Name)

	return db, nil
}
