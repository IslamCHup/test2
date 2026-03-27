package config

import "os"

type Config struct {
	Env      string
	LogLevel string
	Port     string
	DB       DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func Load() *Config {
	return &Config{
		Env:      getEnv("APP_ENV", "local"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Port:     getEnv("APP_PORT", "8080"),

		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "12345"),
			Name:     getEnv("DB_NAME", "subscriptions"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}
