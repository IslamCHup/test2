package logger

import (
"log"
"os"
"strings"
)

type Logger struct {
info  *log.Logger
error *log.Logger
warn  *log.Logger
debug *log.Logger
level string
}

func NewLogger(level string) *Logger {
return &Logger{
info:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
warn:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
debug: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
level: strings.ToLower(level),
}
}

func (l *Logger) Info(msg string, args ...interface{}) {
if l.level != "debug" && l.level != "info" && l.level != "warn" && l.level != "error" {
l.level = "info"
}
if l.level <= "info" {
l.info.Printf(msg, args...)
}
}

func (l *Logger) Error(msg string, args ...interface{}) {
if l.level <= "error" {
l.error.Printf(msg, args...)
}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
if l.level <= "warn" {
l.warn.Printf(msg, args...)
}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
if l.level == "debug" {
l.debug.Printf(msg, args...)
}
}
