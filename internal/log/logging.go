package log

import (
	"fmt"
	"log"
	"os"
)

// Logger represents the logging interface
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// StandardLogger implements the Logger interface using the standard log package
type StandardLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.LstdFlags),
	}
}

// Info logs info level messages
func (l *StandardLogger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.infoLogger.Println(msg)
}

// Error logs error level messages
func (l *StandardLogger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.errorLogger.Println(msg)
}

// Debug logs debug level messages
func (l *StandardLogger) Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.debugLogger.Println(msg)
}

// DefaultLogger is the default logger instance
var DefaultLogger = NewStandardLogger()

// Info logs an info message using the default logger
func Info(format string, args ...interface{}) {
	DefaultLogger.Info(format, args...)
}

// Error logs an error message using the default logger
func Error(format string, args ...interface{}) {
	DefaultLogger.Error(format, args...)
}

// Debug logs a debug message using the default logger
func Debug(format string, args ...interface{}) {
	DefaultLogger.Debug(format, args...)
}
