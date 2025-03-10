package logger

import (
	"fmt"
	"sync"
	"time"
)

// LoggerService represents a custom logger with different log levels
type LoggerService struct {
	mu sync.Mutex
}

var instance *LoggerService
var once sync.Once

// GetInstance returns the singleton instance of LoggerService
func GetInstance() *LoggerService {
	once.Do(func() {
		instance = &LoggerService{}
	})
	return instance
}

// log formats and prints the log message with a timestamp
func (l *LoggerService) log(level string, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("%s [%s] %s", timestamp, level, message)
	fmt.Println(logMessage)
}

// Info logs an info message with formatting support
func (l *LoggerService) Info(format string, args ...any) {
	l.log("INFO", format, args...)
}

// Warn logs a warning message with formatting support
func (l *LoggerService) Warn(format string, args ...any) {
	l.log("WARN", format, args...)
}

// Error logs an error message with formatting support
func (l *LoggerService) Error(format string, args ...any) {
	l.log("ERROR", format, args...)
}
