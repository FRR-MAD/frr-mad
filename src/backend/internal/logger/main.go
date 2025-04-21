package logger

import (
	"fmt"
	"os"
	"sync"
)

const (
	DebugLevelNormal = iota
	DebugLevelError
	DebugLevelDebug
)

type Logger struct {
	name       string
	file       *os.File
	mu         sync.Mutex
	filePath   string
	debugLevel int
}

var (
	registry   = make(map[string]*Logger)
	registryMu sync.Mutex
)

func NewLogger(name, filePath string) (*Logger, error) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if logger, exists := registry[name]; exists {
		return logger, nil
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := &Logger{
		name:       name,
		file:       file,
		filePath:   filePath,
		debugLevel: DebugLevelNormal,
	}

	registry[name] = logger
	return logger, nil
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
	l.mu.Lock()
	defer l.mu.Unlock()

	registryMu.Lock()
	delete(registry, l.name)
	registryMu.Unlock()

	return l.file.Close()
}

func (l *Logger) SetDebugLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debugLevel = level
}

func (l *Logger) GetDebugLevel() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.debugLevel
}
