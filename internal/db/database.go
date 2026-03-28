package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/islamchupanov/tz1/internal/config"
	"github.com/islamchupanov/tz1/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildDSN(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)
}

// InitDB — базовая инициализация БД
func InitDB(cfg config.DBConfig, log *logger.Logger) (*gorm.DB, *sql.DB, error) {
	dsn := buildDSN(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if log != nil {
			log.Error("failed to connect to postgres", "error", err)
		}
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		if log != nil {
			log.Error("failed to get sql.DB from gorm", "error", err)
		}
		return nil, nil, err
	}

	// pool настройки
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// проверка соединения
	if err := sqlDB.Ping(); err != nil {
		if log != nil {
			log.Error("failed to ping db", "error", err)
		}
		return nil, nil, err
	}

	if log != nil {
		log.Info("connected to postgres",
			"host", cfg.Host,
			"port", cfg.Port,
			"dbname", cfg.Name,
		)
	}

	return db, sqlDB, nil
}

// InitDBWithRetry — инициализация с retry
func InitDBWithRetry(
	cfg config.DBConfig,
	log *logger.Logger,
	maxRetries int,
	delay time.Duration,
) (*gorm.DB, *sql.DB, error) {

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		db, sqlDB, err := InitDB(cfg, log)
		if err == nil {
			return db, sqlDB, nil
		}

		lastErr = err

		if i < maxRetries-1 {
			if log != nil {
				log.Warn("database connection failed, retrying",
					"attempt", i+1,
					"max_retries", maxRetries,
					"error", err,
				)
			}
			time.Sleep(delay)
		}
	}

	return nil, nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, lastErr)
}