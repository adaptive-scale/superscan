package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of the log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
)

var (
	// logLevelNames maps LogLevel to string representation
	logLevelNames = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		ERROR: "ERROR",
	}

	// logLevelColors maps LogLevel to ANSI color codes
	logLevelColors = map[LogLevel]string{
		DEBUG: "\033[36m", // Cyan
		INFO:  "\033[32m", // Green
		ERROR: "\033[31m", // Red
	}

	// resetColor resets the ANSI color
	resetColor = "\033[0m"
)

// Logger provides structured logging functionality
type Logger struct {
	*log.Logger
	level LogLevel
}

// New creates a new Logger instance
func New(level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

// formatMessage formats the log message with timestamp, level, and caller information
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)

	// Format timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Format message
	message := fmt.Sprintf(format, args...)

	// Combine all parts
	return fmt.Sprintf("%s%s [%s] %s:%d - %s%s",
		logLevelColors[level],
		timestamp,
		logLevelNames[level],
		file,
		line,
		message,
		resetColor,
	)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.Print(l.formatMessage(DEBUG, format, args...))
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= INFO {
		l.Print(l.formatMessage(INFO, format, args...))
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.Print(l.formatMessage(ERROR, format, args...))
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current logging level
func (l *Logger) GetLevel() LogLevel {
	return l.level
} 