package client

import (
	"fmt"
	"log"
	"os"
)

// DefaultLogger provides a simple logger implementation
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "[ControlApp] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *DefaultLogger) Debug(args ...interface{}) {
	l.logger.Println("DEBUG:", fmt.Sprint(args...))
}

func (l *DefaultLogger) Info(args ...interface{}) {
	l.logger.Println("INFO:", fmt.Sprint(args...))
}

func (l *DefaultLogger) Warn(args ...interface{}) {
	l.logger.Println("WARN:", fmt.Sprint(args...))
}

func (l *DefaultLogger) Error(args ...interface{}) {
	l.logger.Println("ERROR:", fmt.Sprint(args...))
}

func (l *DefaultLogger) WithField(key string, value interface{}) Logger {
	// For simplicity, just return self - real implementation would create a new logger with fields
	return l
}

func (l *DefaultLogger) WithFields(fields map[string]interface{}) Logger {
	// For simplicity, just return self - real implementation would create a new logger with fields
	return l
}
