package logger

import (
	"fmt"
	"sync"
	"time"
)

// Colors for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

// Log levels
const (
	INFO    = "INFO"
	SUCCESS = "SUCCESS"
	ERROR   = "ERROR"
	WARN    = "WARN"
	DEBUG   = "DEBUG"
)

// Logger provides attractive logging functionality
type Logger struct {
	prefix string
}

var (
	instance *Logger
	once     sync.Once
)

// Get returns the singleton logger instance
func Get() *Logger {
	once.Do(func() {
		instance = &Logger{
			prefix: "gozephyr",
		}
	})
	return instance
}

// SetPrefix sets the prefix for the logger
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// formatMessage formats the log message with timestamp and color
func (l *Logger) formatMessage(level, message string, color string) string {
	timestamp := time.Now().Format("15:04:05")
	return fmt.Sprintf("%s[%s]%s %s%s%s %s%s",
		color, timestamp, Reset,
		color, level, Reset,
		l.prefix, message)
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(l.formatMessage(INFO, message, Cyan))
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(l.formatMessage(SUCCESS, message, Green))
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(l.formatMessage(ERROR, message, Red))
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(l.formatMessage(WARN, message, Yellow))
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(l.formatMessage(DEBUG, message, Gray))
}

// Section prints a section header
func (l *Logger) Section(title string) {
	fmt.Printf("\n%s%s%s\n", Purple, title, Reset)
	fmt.Printf("%s%s%s\n", Purple, "="+repeat("=", len(title))+"=", Reset)
}

// SubSection prints a subsection header
func (l *Logger) SubSection(title string) {
	fmt.Printf("\n%s%s%s\n", Blue, title, Reset)
	fmt.Printf("%s%s%s\n", Blue, "-"+repeat("-", len(title))+"-", Reset)
}

// repeat returns a string repeated n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
