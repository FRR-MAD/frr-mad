package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// Level constants that match your original levels
const (
	LevelNormal = iota
	LevelError
	LevelDebug
)
const LevelNone = 99 // Highest number (3) but should prevent all logging

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	name     string
	logger   *slog.Logger
	file     *os.File
	filePath string
	level    int
}

var (
	registry   = make(map[string]*Logger)
	registryMu sync.Mutex
)

// NewApplicationLogger creates the main application logger
func NewApplicationLogger(appName, filePath string) (*Logger, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

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
		level:    LevelNormal,
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

// GetInstance retrieves an existing logger by name
func GetInstance(name string) (*Logger, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	logger, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("logger instance '%s' not found", name)
	}
	return logger, nil
}

// Close closes the logger and removes it from the registry
func (l *Logger) Close() error {
	registryMu.Lock()
	delete(registry, l.name)
	registryMu.Unlock()

	return l.file.Close()
}

// SetDebugLevel sets the debug level for the logger
func (l *Logger) SetDebugLevel(level int) {
	registryMu.Lock()
	defer registryMu.Unlock()
	l.level = level

	// Update the slog level based on our custom level
	var slogLevel slog.Level
	switch level {
	case LevelNormal:
		slogLevel = slog.LevelInfo
	case LevelError:
		slogLevel = slog.LevelWarn
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelNone:
		// Use a level that's higher than any possible log level to suppress all logging
		slogLevel = slog.LevelError + 1000
	default:
		slogLevel = slog.LevelInfo
	}

	// Replace the handler with a new one at the updated level
	opts := &slog.HandlerOptions{Level: slogLevel}
	newHandler := slog.NewJSONHandler(l.file, opts)
	l.logger = slog.New(newHandler).With("logger", l.name)
}

// GetDebugLevel returns the current debug level
func (l *Logger) GetDebugLevel() int {
	registryMu.Lock()
	defer registryMu.Unlock()
	return l.level
}

// SetNormalMode sets the logger to normal mode
func (l *Logger) SetNormalMode() {
	l.SetDebugLevel(LevelNormal)
}

// SetErrorMode sets the logger to error mode
func (l *Logger) SetErrorMode() {
	l.SetDebugLevel(LevelError)
}

// SetDebugMode sets the logger to debug mode
func (l *Logger) SetDebugMode() {
	l.SetDebugLevel(LevelDebug)
}

// SetNoneMode sets the logger to none mode (no logging)
func (l *Logger) SetNoneMode() {
	l.SetDebugLevel(LevelNone)
}

// Info logs an info message
func (l *Logger) Info(msg string) error {
	if l.level != LevelNone {
		l.logger.Info(msg)
	}
	return nil
}

// Error logs an error message
func (l *Logger) Error(msg string) error {
	if l.level != LevelNone {
		l.logger.Error(msg)
	}
	return nil
}

// Debug logs a debug message if the level is high enough
func (l *Logger) Debug(msg string) error {
	if l.level >= LevelDebug && l.level != LevelNone {
		l.logger.Debug(msg)
	}
	return nil
}

// Warning logs a warning message if the level is high enough
func (l *Logger) Warning(msg string) error {
	if l.level >= LevelError && l.level != LevelNone {
		l.logger.Warn(msg)
	}
	return nil
}

// WithAttrs adds structured fields to the logger
func (l *Logger) WithAttrs(attrs map[string]interface{}) *Logger {
	// Create a slice to hold the key-value pairs
	keyValues := make([]any, 0, len(attrs)*2)

	// Populate the slice with alternating keys and values
	for k, v := range attrs {
		keyValues = append(keyValues, k, v)
	}

	// Create a new slog.Logger with the attributes
	newLogger := &Logger{
		name:     l.name,
		logger:   l.logger.With(keyValues...),
		file:     l.file,
		filePath: l.filePath,
		level:    l.level,
	}

	return newLogger
}
