package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

type Logger struct {
	out   io.Writer
	err   io.Writer
	level Level
}

func InitLog(level string) *Logger {
	lvl := parseLevel(level)

	logger := &Logger{
		out:   os.Stdout,
		err:   os.Stderr,
		level: lvl,
	}

	logger.Info("logger initialized", "level", level)
	return logger
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log("INFO", msg, args...)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log("DEBUG", msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log("WARN", msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logErr("ERROR", msg, args...)
}

func (l *Logger) log(prefix, msg string, args ...interface{}) {
	formatted := formatLog(msg, args...)
	fmt.Fprintf(l.out, "{\"level\":\"%s\",\"msg\":\"%s\"}\n", prefix, formatted)
}

func (l *Logger) logErr(prefix, msg string, args ...interface{}) {
	formatted := formatLog(msg, args...)
	fmt.Fprintf(l.err, "{\"level\":\"%s\",\"msg\":\"%s\"}\n", prefix, formatted)
}

func formatLog(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}

	var sb strings.Builder
	sb.WriteString(msg)

	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				sb.WriteString(fmt.Sprintf(",\"%s\":\"%v\"", key, args[i+1]))
			}
		}
	}

	return sb.String()
}
