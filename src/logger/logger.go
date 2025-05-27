package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

const LevelNone = 1000 // Highest number (3) but should prevent all logging

type Logger struct {
	name     string
	logger   *slog.Logger
	file     *os.File
	filePath string
	level    slog.Level
}

var (
	registry   = make(map[string]*Logger)
	registryMu sync.Mutex
)

func NewApplicationLogger(appName, filePath string) (*Logger, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Slog is thread-safe, each log Operation like Info() or Error() are atomic
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(file, opts)

	slogger := slog.New(handler).With("application", appName)

	logger := &Logger{
		name:     appName,
		logger:   slogger,
		file:     file,
		filePath: filePath,
		level:    slog.LevelInfo,
	}

	registry[appName] = logger
	return logger, nil
}

func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		name:     l.name,
		logger:   l.logger.With("component", component),
		file:     l.file,
		filePath: l.filePath,
		level:    l.level,
	}
}

func GetInstance(name string) (*Logger, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	logger, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("logger instance '%s' not found", name)
	}
	return logger, nil
}

func (l *Logger) Close() error {
	registryMu.Lock()
	delete(registry, l.name)
	registryMu.Unlock()

	return l.file.Close()
}

func (l *Logger) SetDebugLevel(level slog.Level) {
	registryMu.Lock()
	defer registryMu.Unlock()

	l.level = level
	opts := &slog.HandlerOptions{Level: level}
	newHandler := slog.NewJSONHandler(l.file, opts)
	l.logger = slog.New(newHandler).With("logger", l.name)
}

func ConvertLogLevelFromConfig(level string) slog.Level {
	switch level {
	case "info":
		return slog.LevelInfo
	case "none":
		return LevelNone
	case "debug":
		return slog.LevelDebug
	case "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *Logger) GetDebugLevel() slog.Level {
	registryMu.Lock()
	defer registryMu.Unlock()
	return l.level
}

func (l *Logger) SetInfoMode() {
	l.SetDebugLevel(slog.LevelInfo)
}

func (l *Logger) SetWarningMode() {
	l.SetDebugLevel(slog.LevelWarn)
}

func (l *Logger) SetErrorMode() {
	l.SetDebugLevel(slog.LevelError)
}

func (l *Logger) SetDebugMode() {
	l.SetDebugLevel(slog.LevelDebug)
}

func (l *Logger) SetNoneMode() {
	l.SetDebugLevel(LevelNone)
}

func (l *Logger) Debug(msg string) error {
	if l.level <= slog.LevelDebug && l.level != LevelNone {
		l.logger.Debug(msg)
	}
	return nil
}

func (l *Logger) Info(msg string) error {
	if l.level <= slog.LevelInfo && l.level != LevelNone {
		l.logger.Info(msg)
	}
	return nil
}

func (l *Logger) Warning(msg string) error {
	if l.level <= slog.LevelWarn && l.level != LevelNone {
		l.logger.Warn(msg)
	}
	return nil
}

func (l *Logger) Error(msg string) error {
	if l.level <= slog.LevelError && l.level != LevelNone {
		l.logger.Error(msg)
	}
	return nil
}

func (l *Logger) WithAttrs(attrs map[string]interface{}) *Logger {

	keyValues := make([]any, 0, len(attrs)*2)

	for k, v := range attrs {
		keyValues = append(keyValues, k, v)
	}

	newLogger := &Logger{
		name:     l.name,
		logger:   l.logger.With(keyValues...),
		file:     l.file,
		filePath: l.filePath,
		level:    l.level,
	}

	return newLogger
}
