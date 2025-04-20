package logger

import (
	"fmt"
	"time"
)

func (l *Logger) Info(msg string) error {
	return l.log("INFO", msg)
}

func (l *Logger) Error(msg string) error {
	return l.log("ERROR", msg)
}

func (l *Logger) Debug(msg string) error {
	return l.log("DEBUG", msg)
}

func (l *Logger) Warning(msg string) error {
	return l.log("WARNING", msg)
}

func (l *Logger) log(level, msg string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] [%s] %s\n", timestamp, l.name, level, msg)

	l.mu.Lock()
	defer l.mu.Unlock()

	_, err := l.file.WriteString(logLine)
	return err
}
