package logging

import (
	"log"
	"os"
)

// Logger is a simple logger that writes to the console.
type Logger struct {
	*log.Logger
}

// NewLogger creates a new Logger.
func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info logs an informational message.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Printf("INFO: "+msg, args...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Printf("ERROR: "+msg, args...)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Printf("DEBUG: "+msg, args...)
}
