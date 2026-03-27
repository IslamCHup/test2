package logger

import (
	"log/slog"
	"os"
	"strings"
)

func parseLog(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func InitLog(level string) *slog.Logger {
	levelLog := parseLog(level)

	handlersLogger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: levelLog,
	})

	logger := slog.New(handlersLogger)
	slog.SetDefault(logger)

	slog.Info("logger инициализирован: ", "level", levelLog.String())

	return logger
}

