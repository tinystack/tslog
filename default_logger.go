// Package tslog provides default logger functionality for quick and easy logging.
// This file contains the default logger instance and convenience functions
// that wrap the default logger methods.
package tslog

import (
	"sync"

	"github.com/tinystack/tslog/writer"
)

// defaultLogger holds the global default logger instance.
// It's initialized once during package initialization and can be
// safely accessed from multiple goroutines.
var defaultLogger Logger

// defaultLoggerMutex protects defaultLogger updates to ensure thread safety.
var defaultLoggerMutex sync.RWMutex

// defaultLogLevel defines the initial log level for the default logger.
// This can be overridden by creating a new logger and calling UpdateDefaultLogger.
var (
	defaultLogLevel = DebugLevel
)

// init initializes the default logger with sensible defaults.
// This ensures that the package can be used immediately without any configuration.
// The default logger uses:
// - Debug level logging
// - Console encoder for human-readable output
// - Standard output as the destination
// - Caller information disabled for cleaner output
func init() {
	funcOpts := []FuncOption{
		WithLevel(defaultLogLevel),
		WithWriter(writer.NewStdoutWriter()),
		WithEncoder(EncoderConsole),
		WithCaller(false),
	}

	opts := defaultOptions()
	for _, f := range funcOpts {
		if f != nil {
			f(opts)
		}
	}

	defaultLogger = opts.driver(opts)
}

// DefaultLogger returns the current default logger instance.
// This function is thread-safe and can be called from multiple goroutines.
//
// Example:
//
//	logger := tslog.DefaultLogger()
//	logger.Info("Using default logger")
func DefaultLogger() Logger {
	defaultLoggerMutex.RLock()
	defer defaultLoggerMutex.RUnlock()
	return defaultLogger
}

// UpdateDefaultLogger replaces the global default logger with a new instance.
// This operation is thread-safe and will affect all subsequent calls to
// the package-level logging functions (Debug, Info, Warn, Error, etc.).
//
// Example:
//
//	customLogger := tslog.NewLogger(tslog.WithLevel(tslog.InfoLevel))
//	tslog.UpdateDefaultLogger(customLogger)
//
// Note: This function should be called early in your application's lifecycle,
// preferably during initialization, to avoid confusion about which logger
// is being used.
func UpdateDefaultLogger(l Logger) {
	if l == nil {
		return // Ignore nil loggers to prevent panics
	}

	defaultLoggerMutex.Lock()
	defer defaultLoggerMutex.Unlock()
	defaultLogger = l
}

// Package-level convenience functions that delegate to the default logger.
// These functions provide a simple API for applications that don't need
// multiple logger instances or complex configuration.

// Debug logs a message at Debug level using the default logger.
// Arguments are handled in the manner of fmt.Print.
//
// Example:
//
//	tslog.Debug("Debug message")
//	tslog.Debug("User ID:", userID, "Action:", action)
func Debug(args ...interface{}) {
	DefaultLogger().Debug(args...)
}

// Info logs a message at Info level using the default logger.
// Arguments are handled in the manner of fmt.Print.
//
// Example:
//
//	tslog.Info("Application started")
//	tslog.Info("Processing request for user:", userID)
func Info(args ...interface{}) {
	DefaultLogger().Info(args...)
}

// Warn logs a message at Warn level using the default logger.
// Arguments are handled in the manner of fmt.Print.
//
// Example:
//
//	tslog.Warn("Deprecated API used")
//	tslog.Warn("High memory usage detected:", memUsage)
func Warn(args ...interface{}) {
	DefaultLogger().Warn(args...)
}

// Error logs a message at Error level using the default logger.
// Arguments are handled in the manner of fmt.Print.
//
// Example:
//
//	tslog.Error("Failed to connect to database")
//	tslog.Error("Error processing request:", err)
func Error(args ...interface{}) {
	DefaultLogger().Error(args...)
}

// Debugf logs a formatted message at Debug level using the default logger.
// Arguments are handled in the manner of fmt.Printf.
//
// Example:
//
//	tslog.Debugf("User %d performed action %s", userID, action)
func Debugf(format string, args ...interface{}) {
	DefaultLogger().Debugf(format, args...)
}

// Infof logs a formatted message at Info level using the default logger.
// Arguments are handled in the manner of fmt.Printf.
//
// Example:
//
//	tslog.Infof("Request processed in %v", duration)
func Infof(format string, args ...interface{}) {
	DefaultLogger().Infof(format, args...)
}

// Warnf logs a formatted message at Warn level using the default logger.
// Arguments are handled in the manner of fmt.Printf.
//
// Example:
//
//	tslog.Warnf("Memory usage is %d%%, consider optimization", memPercent)
func Warnf(format string, args ...interface{}) {
	DefaultLogger().Warnf(format, args...)
}

// Errorf logs a formatted message at Error level using the default logger.
// Arguments are handled in the manner of fmt.Printf.
//
// Example:
//
//	tslog.Errorf("Failed to process request: %v", err)
func Errorf(format string, args ...interface{}) {
	DefaultLogger().Errorf(format, args...)
}

// Debugt logs a message with structured fields at Debug level using the default logger.
// This method is useful for structured logging with key-value pairs.
//
// Example:
//
//	tslog.Debugt("User action", tslog.T{
//	    "user_id": 123,
//	    "action": "login",
//	    "ip": "192.168.1.1",
//	})
func Debugt(msg string, args T) {
	DefaultLogger().Debugt(msg, args)
}

// Infot logs a message with structured fields at Info level using the default logger.
// This method is useful for structured logging with key-value pairs.
//
// Example:
//
//	tslog.Infot("Request completed", tslog.T{
//	    "method": "GET",
//	    "path": "/api/users",
//	    "duration": "150ms",
//	})
func Infot(msg string, args T) {
	DefaultLogger().Infot(msg, args)
}

// Warnt logs a message with structured fields at Warn level using the default logger.
// This method is useful for structured logging with key-value pairs.
//
// Example:
//
//	tslog.Warnt("High latency detected", tslog.T{
//	    "service": "database",
//	    "latency": "2.5s",
//	    "threshold": "1s",
//	})
func Warnt(msg string, args T) {
	DefaultLogger().Warnt(msg, args)
}

// Errort logs a message with structured fields at Error level using the default logger.
// This method is useful for structured logging with key-value pairs.
//
// Example:
//
//	tslog.Errort("Database connection failed", tslog.T{
//	    "host": "db.example.com",
//	    "port": 5432,
//	    "error": err.Error(),
//	})
func Errort(msg string, args T) {
	DefaultLogger().Errort(msg, args)
}
