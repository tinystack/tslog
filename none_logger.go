// Package tslog provides a no-operation logger implementation.
// This file contains the NoneLogger that implements the Logger interface
// but performs no actual logging operations, useful for disabling logging
// entirely with zero performance overhead.
package tslog

// NoneLogger is a no-operation logger that implements the Logger interface
// but discards all log messages. This is useful when you want to disable
// logging entirely while maintaining the same API.
//
// All methods are implemented as no-ops with zero allocations and minimal
// performance overhead, making it safe to use in performance-critical code
// paths where logging needs to be disabled.
//
// Example usage:
//
//	var logger tslog.Logger = &tslog.NoneLogger{}
//	logger.Info("This message will be discarded")
type NoneLogger struct{}

// NewNoneLogger creates a new NoneLogger instance.
// This function is provided for consistency with other logger constructors,
// though creating a NoneLogger directly with &NoneLogger{} is also valid.
//
// Example:
//
//	logger := tslog.NewNoneLogger()
func NewNoneLogger() Logger {
	return &NoneLogger{}
}

// Debug discards the debug message. This is a no-op method.
// Arguments are ignored and no processing is performed.
func (*NoneLogger) Debug(args ...interface{}) {}

// Info discards the info message. This is a no-op method.
// Arguments are ignored and no processing is performed.
func (*NoneLogger) Info(args ...interface{}) {}

// Warn discards the warning message. This is a no-op method.
// Arguments are ignored and no processing is performed.
func (*NoneLogger) Warn(args ...interface{}) {}

// Error discards the error message. This is a no-op method.
// Arguments are ignored and no processing is performed.
func (*NoneLogger) Error(args ...interface{}) {}

// Debugf discards the formatted debug message. This is a no-op method.
// Format string and arguments are ignored and no processing is performed.
func (*NoneLogger) Debugf(format string, args ...interface{}) {}

// Infof discards the formatted info message. This is a no-op method.
// Format string and arguments are ignored and no processing is performed.
func (*NoneLogger) Infof(format string, args ...interface{}) {}

// Warnf discards the formatted warning message. This is a no-op method.
// Format string and arguments are ignored and no processing is performed.
func (*NoneLogger) Warnf(format string, args ...interface{}) {}

// Errorf discards the formatted error message. This is a no-op method.
// Format string and arguments are ignored and no processing is performed.
func (*NoneLogger) Errorf(format string, args ...interface{}) {}

// Debugt discards the structured debug message. This is a no-op method.
// Message and structured fields are ignored and no processing is performed.
func (*NoneLogger) Debugt(msg string, args T) {}

// Infot discards the structured info message. This is a no-op method.
// Message and structured fields are ignored and no processing is performed.
func (*NoneLogger) Infot(msg string, args T) {}

// Warnt discards the structured warning message. This is a no-op method.
// Message and structured fields are ignored and no processing is performed.
func (*NoneLogger) Warnt(msg string, args T) {}

// Errort discards the structured error message. This is a no-op method.
// Message and structured fields are ignored and no processing is performed.
func (*NoneLogger) Errort(msg string, args T) {}

// NewNoneDriver creates a Driver function that returns a NoneLogger.
// This can be used with NewLogger to create a no-op logger instance.
//
// Example:
//
//	logger := tslog.NewLogger(tslog.WithDriver(tslog.NewNoneDriver))
func NewNoneDriver(opts *Options) Logger {
	return NewNoneLogger()
}
